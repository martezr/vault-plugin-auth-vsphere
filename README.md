HashiCorp Vault VMware vSphere Authentication Plugin
=======

[![Build Status](https://img.shields.io/travis/martezr/vault-plugin-auth-vsphere/master.svg)][travis]
[![GoReportCard][report-badge]][report]
[![GitHub release](https://img.shields.io/github/release/martezr/vault-plugin-auth-vsphere.svg)](https://github.com/martezr/vault-plugin-auth-vsphere/releases/)
[![license](https://img.shields.io/github/license/martezr/vault-plugin-auth-vsphere.svg)](https://github.com/martezr/vault-plugin-auth-vsphere/blob/master/LICENSE)

[travis]: https://travis-ci.org/martezr/vault-plugin-auth-vsphere

[report-badge]: https://goreportcard.com/badge/github.com/martezr/vault-plugin-auth-vsphere
[report]: https://goreportcard.com/report/github.com/martezr/vault-plugin-auth-vsphere

This is a standalone backend plugin for use with HashiCorp Vault. This plugin allows for VMware vSphere virtual machines (VMs) to authenticate with Vault.

## Getting Started

This is a [Vault plugin](https://www.vaultproject.io/docs/internals/plugins.html) and is meant to work with Vault. This guide assumes you have already installed Vault and have a basic understanding of how Vault works.

To learn specifically about how plugins work, see documentation on [Vault plugins](https://www.vaultproject.io/docs/internals/plugins.html).

## Setup

Download the latest plugin binary from the [Releases](https://github.com/martezr/vault-plugin-auth-vsphere/releases) page on GitHub and move the plugin binary into Vault's configured *plugin_directory*.

```
$ mv vault-plugin-auth-vsphere /etc/vault/plugins/vault-plugin-auth-vspherre
```

Calculate the checksum of the plugin and register it in Vault's plugin catalog. It is highly recommended that you use the published checksums on the Release page to verify integrity.

```
$ export SHA256_SUM=$(shasum -a 256 "/etc/vault/plugins/vault-plugin-auth-vsphere" | cut -d' ' -f1)
$ vault write sys/plugins/catalog/auth/vsphere \
    command="vault-plugin-auth-vsphere" \
    sha_256="${SHA256_SUM}"
```

Enable authentication with the plugin.

```
$ vault auth enable -path="vsphere" -plugin-name="vsphere" plugin
```