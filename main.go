package winrm

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/masterzen/winrm"
)

type optionsStruct struct {
	// Target has HTTPS WinRM endpoint
	HTTPS bool `compspore:"https"`

	// Target is running in insecure mode
	Insecure bool `compspore:"insecure"`

	// CA cert for the target
	CACert []byte `compspore:"cacert"`

	// Client cert for the target
	Cert []byte `compspore:"cert"`

	// Client key for the target
	Key []byte `compspore:"key"`
}

func (o *optionsStruct) Unmarshal(options map[string]interface{}) error {
	httpsInterface, ok := options["https"]
	if ok {
		https, ok := httpsInterface.(bool)
		if ok {
			o.HTTPS = https
		} else {
			o.HTTPS = false
		}
	}

	insecureInterface, ok := options["insecure"]
	if ok {
		insecure, ok := insecureInterface.(bool)
		if ok {
			o.Insecure = insecure
		} else {
			o.Insecure = false
		}
	}

	caCertInterface, ok := options["cacert"]
	if ok {
		caCert, ok := caCertInterface.(string)
		if ok {
			o.CACert = []byte(caCert)
		} else {
			return fmt.Errorf("cacert must be a string, provided: %T", caCertInterface)
		}
	}

	certInterface, ok := options["cert"]
	if ok {
		cert, ok := certInterface.(string)
		if ok {
			o.Cert = []byte(cert)
		} else {
			return fmt.Errorf("cert must be a string, provided: %T", certInterface)
		}
	}

	keyInterface, ok := options["key"]
	if ok {
		key, ok := keyInterface.(string)
		if ok {
			o.Key = []byte(key)
		} else {
			return fmt.Errorf("key must be a string, provided: %T", keyInterface)
		}
	}

	return nil
}
func Run(ctx context.Context, target string, command string, expectedOutput string, username string, password string, options map[string]interface{}) (bool, string) {
	optionsStruct := optionsStruct{}
	err := optionsStruct.Unmarshal(options)
	if err != nil {
		return false, fmt.Sprintf("failed to unmarshal options: %s", err)
	}

	var port int

	if strings.Contains(target, ":") {
		split := strings.Split(target, ":")
		if len(split) != 2 {
			return false, "target must be in the format of <host>:<port> or <host>"
		}

		target = split[0]
		p, err := strconv.Atoi(split[1])
		if err != nil {
			return false, "port must be an integer"
		}

		port = p
	} else {
		if optionsStruct.HTTPS {
			port = 5986
		} else {
			port = 5985
		}
	}

	endpoint := winrm.NewEndpoint(
		target,
		port,
		optionsStruct.HTTPS,
		optionsStruct.Insecure,
		optionsStruct.CACert,
		optionsStruct.Cert,
		optionsStruct.Key,
		0,
	)

	client, err := winrm.NewClient(endpoint, username, password)
	if err != nil {
		return false, fmt.Sprintf("failed to create winrm client: %s", err)
	}

	outputChan := make(chan string)
	errChan := make(chan error)

	go func() {
		defer close(outputChan)
		defer close(errChan)

		stdout, stderr, _, err := client.RunWithContextWithString(ctx, command, "")
		if err != nil {
			errChan <- err
			return
		}

		if stderr != "" {
			errChan <- fmt.Errorf("command resulted in error: %s", stderr)
			return
		}

		outputChan <- stdout
	}()

	select {
	case <-ctx.Done():
		return false, fmt.Sprintf("command timed out: %s", ctx.Err())
	case out := <-outputChan:
		if expectedOutput != "" {
			if strings.TrimSpace(out) != strings.TrimSpace(expectedOutput) {
				return false, fmt.Sprintf("expected output not found; expected: \"%s\"; recieved: \"%s\"", expectedOutput, out)
			}
		}

		return true, ""
	case err := <-errChan:
		return false, fmt.Sprintf("command resulted in error: %s", err)
	}
}
