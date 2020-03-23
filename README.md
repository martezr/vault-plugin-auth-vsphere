HashiCorp Vault VMware vSphere Authentication Plugin
=======

[![Build Status](https://img.shields.io/travis/martezr/vault-plugin-auth-vsphere/master.svg)][travis]
[![GoReportCard][report-badge]][report]
[![GitHub release](https://img.shields.io/github/release/martezr/vault-plugin-auth-vsphere.svg)](https://github.com/martezr/vault-plugin-auth-vsphere/releases/)
[![license](https://img.shields.io/github/license/martezr/vault-plugin-auth-vsphere.svg)](https://github.com/martezr/vault-plugin-auth-vsphere/blob/master/LICENSE)

[travis]: https://travis-ci.org/martezr/vault-plugin-auth-vsphere

[report-badge]: https://goreportcard.com/badge/github.com/martezr/vault-plugin-auth-vsphere
[report]: https://goreportcard.com/report/github.com/martezr/vault-plugin-auth-vsphere

This is a standalone backend plugin for use with HashiCorp Vault. This plugin allows for VMware vSphere virtual machines (VMs) to authenticate with Vault. This plugin requires the vAuth platform to be deployed and configured.

## Getting Started

This is a [Vault plugin](https://www.vaultproject.io/docs/internals/plugins.html) and is meant to work with Vault. This guide assumes you have already installed Vault and have a basic understanding of how Vault works.

To learn specifically about how plugins work, see documentation on [Vault plugins](https://www.vaultproject.io/docs/internals/plugins.html).

## Setup

Download the latest plugin binary from the [Releases](https://github.com/martezr/vault-plugin-auth-vsphere/releases) page on GitHub and move the plugin binary into Vault's configured *plugin_directory*.

```
$ mv vault-plugin-auth-vsphere /etc/vault/plugins/vault-plugin-auth-vsphere
```

Calculate the checksum of the plugin and register it in Vault's plugin catalog. It is highly recommended that you use the published checksums on the Release page to verify integrity.

```
$ export SHA256_SUM=$(shasum -a 256 "/opt/vault/plugins/vault-plugin-auth-vsphere" | cut -d' ' -f1)
$ vault write sys/plugins/catalog/auth/vsphere command="vault-plugin-auth-vsphere" sha_256="${SHA256_SUM}"
```

Enable authentication with the plugin.

```
$ vault auth enable vsphere
```

Configure the vAuth URL
```
$ vault write auth/vsphere/config vauth_url="https://vauth.grt.local"
```

## Developing

If you wish to work on this plugin, you'll first need
[Go](https://www.golang.org) installed on your machine
(version 1.13+ is *required*).

Compile the plugin binary for use with Vault

```shell
go build -o plugins/vsphere cmd/vault-plugin-auth-vsphere/main.go
```

Run HashiCorp Vault in dev mode with the plugin automatically loaded

```shell
vault server -dev -dev-root-token-id=root -dev-listen-address=0.0.0.0:8200 -dev-plugin-dir=./plugins
```

```shell
export VAULT_ADDR="http://127.0.0.1:8200"
```

Enable the vSphere authentication plugin

```
vault auth enable -path="vsphere" -plugin-name="vsphere" plugin
```

## License

|                |                                                  |
| -------------- | ------------------------------------------------ |
| **Author:**    | Martez Reed (<martez.reed@greenreedtech.com>)    |
| **Copyright:** | Copyright (c) 2019 Green Reed Technology    |
| **License:**   | Apache License, Version 2.0                      |

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.