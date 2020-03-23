#!/bin/bash

export SHA256_SUM=$(sha256sum "/opt/vault/plugins/vault-plugin-auth-vsphere" | cut -d' ' -f1)

vault write sys/plugins/catalog/auth/vsphere command="vault-plugin-auth-vsphere" sha_256="${SHA256_SUM}"

vault auth enable vsphere

vault write auth/vsphere/config vauth_url="https://localhost"

vault write auth/vsphere/role/webapp policies=webapp
