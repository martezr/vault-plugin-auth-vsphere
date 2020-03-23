#!/bin/bash

export VAULT_ADDR="http://127.0.0.1:8200"

vault login root

vault auth enable vsphere

vault write auth/vsphere/config vauth_url="https://localhost"

vault policy write webapp vault-policy.hcl

vault write auth/vsphere/role/webapp policies=webapp

vault read auth/vsphere/role/webapp

vault list auth/vsphere/roles