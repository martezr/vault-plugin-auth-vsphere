#!/bin/bash

VAULT_ADDR="http://10.0.0.206:8200"
NAME=$(/usr/bin/vmtoolsd --cmd 'info-get guestinfo.vault.vmname')
KEY=$(/usr/bin/vmtoolsd --cmd 'info-get guestinfo.vault.secretkey')
DC=$(/usr/bin/vmtoolsd --cmd 'info-get guestinfo.vault.datacenter')
ROLE=$(/usr/bin/vmtoolsd --cmd 'info-get guestinfo.vault.role')
echo $ROLE

curl --request POST --data '{"secretkey":"'"$KEY"'","vmname": "'"$NAME"'","datacenter":"'"$DC"'","role": "'"$ROLE"'"}' http://10.0.0.206:8200/v1/auth/vsphere/login