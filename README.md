# WinRM

## Check Format

```yaml
- name:
  release:
    org: compscore
    repo: smb
    tag: latest
  credentials:
    username:
    password:
  target:
  command:
  expectedOutput:
  weight:
  options:
    status_code:
    match:
    substring_match:
    regex_match:
```

## Parameters

|    parameter     |          path           |   type   | default  | required | description                                     |
| :--------------: | :---------------------: | :------: | :------: | :------: | :---------------------------------------------- |
|      `name`      |         `.name`         | `string` |   `""`   |  `true`  | `name of check (must be unique)`                |
|      `org`       |     `.release.org`      | `string` |   `""`   |  `true`  | `organization that check repository belongs to` |
|      `repo`      |     `.release.repo`     | `string` |   `""`   |  `true`  | `repository of the check`                       |
|      `tag`       |     `.release.tag`      | `string` | `latest` | `false`  | `tagged version of check`                       |
|    `username`    | `.credentials.username` | `string` |   `""`   | `false`  | `username for winrm user`                       |
|    `password`    | `.credentials.password` | `string` |   `""`   | `false`  | `default password for winrm user`               |
|     `target`     |        `.target`        | `string` |   `""`   |  `true`  | `network target for winrm server`               |
|    `command`     |       `.command`        | `string` |   `""`   | `false`  | `command to execute as remote user`             |
| `expectedOutput` |    `.expectedOutput`    | `string` |   `""`   | `false`  | `expected output of provided command`           |
|     `weight`     |        `.weight`        |  `int`   |   `0`    |  `true`  | `amount of points a successful check is worth`  |
|     `https`      |    `.options.https`     |  `bool`  | `false`  | `false`  | `use https to establish winrm connection`       |
|    `insecure`    |   `.options.insecure`   |  `bool`  | `false`  | `false`  | `establish winrm connection as insecure`        |
|     `cacert`     |    `.options.cacert`    | `string` |   `""`   | `false`  | `establish winrm connection with pinned cacert` |
|      `cert`      |     `.options.cert`     | `string` |   `""`   | `false`  | `establish winrm connection with client cert`   |
|      `key`       |     `.options.key`      | `string` |   `""`   | `false`  | `establish winrm connection with client key`    |

## Examples

```yaml
- name: host_a-winrm
  release:
    org: compscore
    repo: winrm
    tag: latest
  credentials:
    username: Administrator
    password: changeme
  target: 10.{{ .Team }}.1.1:5985
  command: whoami
  expectedOutput: Administrator
  weight: 2
  options:
    https: false
    insecure: true
```
