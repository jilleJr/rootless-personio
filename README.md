<!--
SPDX-FileCopyrightText: 2023 Kalle Fagerberg

SPDX-License-Identifier: CC-BY-4.0
-->

# "Rootless" Personio API client

[![REUSE status](https://api.reuse.software/badge/github.com/jilleJr/rootless-personio)](https://api.reuse.software/info/github.com/jilleJr/rootless-personio)
[![Go Reference](https://pkg.go.dev/badge/github.com/jilleJr/rootless-personio/pkg/personio.svg)](https://pkg.go.dev/github.com/jilleJr/rootless-personio/pkg/personio)

Accessing [Personio's API](https://developer.personio.de/docs)
requires API credentials [which does not scope to the employee level](https://developer.personio.de/discuss/634e4b08a3f8d80051c52cfe),
meaning you can only get official API access as an admin user,
where you get access to the sensitive information of all the employees in your
company.

Instead, this package uses a different API: the same API as your web browser.

This is done by pretending to be a browser and logging in normally using
email and password.

## CLI

### Installing CLI

There is a CLI available that exposes some of the features found in the
Go library.

Currently, installing **requires that you have Go 1.19 (or later) installed.**
Then you can run the following:

```sh
go install github.com/jilleJr/rootless-personio@latest
```

### CLI usage

```console
$ rootless-personio --help
Access Personio via your own employee credentials,
instead of obtaining admin/root API credentials.

Usage:
  rootless-personio [command]

Available Commands:
  attendance  Group of commands for interacting with attendance
  completion  Generate the autocompletion script for the specified shell
  config      Prints the parsed config
  help        Help about any command
  raw         Send a raw HTTP request to the API

Flags:
      --auth.email string       Email used when logging in
      --config string           Config file (default is $HOME/.rootless-personio.yaml)
  -h, --help                    Show this help text
      --log.format log-format   Sets the logging format (default pretty)
      --log.level log-level     Sets the logging level (default warn)
      --no-login                Skip logging in before the request
  -o, --output out-format       Sets the output format (default pretty)
  -q, --quiet                   Disables logging (same as "--log.level disabled")
      --url string              Base URL used to access Personio
  -v, --verbose count           Shows verbose logging (-v=info, -vv=debug, -vvv=trace)

Use "rootless-personio [command] --help" for more information about a command.
```

### Configuration

The CLI is configured via YAML files.
See [`personio.yaml`](./personio.yaml) for the default values.

#### Configuration files

Certmgmt looks for config files in multiple locations, where the latter
overrides config fields from the former.

On Linux:

1. Default values *(see [`personio.yaml`](./personio.yaml))*
2. `/etc/personio/personio.yaml`
3. `~/.config/personio.yaml`
4. `~/.personio.yaml`
5. `.personio.yaml` *(in current directory)*

On Windows:

1. Default values *(see [`personio.yaml`](./personio.yaml))*
2. `%APPDATA%/personio.yaml`
3. `%USERPROFILE%/.personio.yaml`
4. `.personio.yaml` *(in current directory)*

On Mac:

1. Default values *(see [`personio.yaml`](./personio.yaml))*
2. `/etc/personio/personio.yaml`
3. `~/Library/Application Support/personio.yaml`
4. `~/.personio.yaml`
5. `.personio.yaml` *(in current directory)*

#### JSON Schema

There's also a [JSON Schema](https://json-schema.org/) for the config file,
which gives you warnings and completion support inside your IDE.

Make use of it via e.g:

- [YAML extension](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)
  for [VS Code](https://code.visualstudio.com/).

- [coc-yaml plugin](https://github.com/neoclide/coc-yaml)
  for [coc.nvim](https://github.com/neoclide/coc.nvim),
  an extension framework for both Vim and NeoVim.

To make use of it, add the following comment to the beginning of your
config file:

```yaml
# yaml-language-server: $schema=https://github.com/jilleJr/rootless-personio/raw/main/personio.schema.json
```

## License

This repository was created by [@jorie1234](https://github.com/jorie1234)
under the MIT license, but has been forked and is now maintained by
Kalle Fagerberg ([@jilleJr](https://github.com/jilleJr)) under a new license.

The code in this project is licensed under GNU General Public License v3.0
or later ([LICENSES/GPL-3.0-or-later.txt](LICENSES/GPL-3.0-or-later.txt)),
and documentation is licensed under Creative Commons Attribution 4.0
International ([LICENSES/CC-BY-4.0.txt](LICENSES/CC-BY-4.0.txt)).

## Credits

Code in this repository is heavily inspired by:

- the upstream work from [@jorie1234](https://github.com/jorie1234):
  <https://github.com/jorie1234/goPersonio>

- Eduardo Sánchez's ([@Whipshout](https://github.com/Whipshout))
  Rust implementation: <https://github.com/Whipshout/personio_tool>
