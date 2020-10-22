#!/bin/bash

# The following deploys out a basic Vault instance based on filesystem storage (into Docker volumes)
#
# It will deliberately destory any existing Vault instance - you have been warned !
#
# Notes:
# - password generator plugin is installed
# - GUI is disabled
# - memory lock is enabled

###################################################################################
# Options
###################################################################################
# Only option is to choose where the script outputs it's helper files
# As per comments below, these files contain *all* keys to Vault - in production this information would
# be managed in a much more secure way :)
HELPER_FILES_DIR=~/vault

###################################################################################
# Main script
###################################################################################
if ! command -v vault &> /dev/null
then
    echo "Vault could not be found, please install it before running this script"
    exit
fi

sh vault_down.sh

echo ==============================================================================
echo Starting fresh Vault ...
docker run -d \
  --name vault \
  --restart unless-stopped \
  --cap-add IPC_LOCK \
  -e 'VAULT_LOCAL_CONFIG={"backend": {"file": {"path": "/vault/file"}}, "listener": {"tcp": {"address": "0.0.0.0:8888", "tls_disable": 1}}, "ui": 0, "default_lease_ttl": "168h", "max_lease_ttl": "720h", "plugin_directory": "/vault/plugins"}' \
  -e 'VAULT_ADDR=http://0.0.0.0:8888' \
  -e 'VAULT_API_ADDR=http://0.0.0.0:8888' \
  -v vault-config:/vault/config \
  -v vault-policies:/vault/policies \
  -v vault-data:/vault/data \
  -v vault-logs:/vault/logs \
  -p 8888:8888 \
  dcgsteve/vault:002 server

export VAULT_ADDR=http://127.0.0.1:8888

echo ==============================================================================
echo Pause for Vault to come up ...
sleep 1

echo ==============================================================================
echo "Initialising Vault ..."
mkdir -p $HELPER_FILES_DIR
vault operator init > $HELPER_FILES_DIR/vault-info

echo ==============================================================================
echo "Writing out developer helper files ..."
echo "  vault-env     = environment variables"
echo "  vault-unseal  = unseal script"
echo "  vault-info    = full key and token info for Vault"
echo "  vault-plugins = enable use of secrets plugin"
echo ""
echo "Obviously these helper files are only for development not for production !"

echo export VAULT_ADDR=http://127.0.0.1:8888 > $HELPER_FILES_DIR/vault-env
echo export VAULT_TOKEN=$( cat $HELPER_FILES_DIR/vault-info | grep 'Initial Root Token' | awk -F ' ' '{print $4}' ) >> $HELPER_FILES_DIR/vault-env

echo vault operator unseal $( cat $HELPER_FILES_DIR/vault-info | grep 'Unseal Key 1' | awk -F ' ' '{print $4}' ) > $HELPER_FILES_DIR/vault-unseal
echo vault operator unseal $( cat $HELPER_FILES_DIR/vault-info | grep 'Unseal Key 2' | awk -F ' ' '{print $4}' ) >> $HELPER_FILES_DIR/vault-unseal
echo vault operator unseal $( cat $HELPER_FILES_DIR/vault-info | grep 'Unseal Key 3' | awk -F ' ' '{print $4}' ) >> $HELPER_FILES_DIR/vault-unseal

echo vault plugin register -sha256="ddaef75e7b7653e34e8b5efebe6253381a423428b68544cd79149deaff8b5f4e" -command="vault-secrets-gen" secret secrets-gen > $HELPER_FILES_DIR/vault-plugins
echo docker exec vault setcap cap_ipc_lock=+ep /vault/plugins/vault-secrets-gen >> $HELPER_FILES_DIR/vault-plugins

chmod u+x $HELPER_FILES_DIR/vault-unseal
chmod u+x $HELPER_FILES_DIR/vault-plugins

ls -l $HELPER_FILES_DIR

echo ==============================================================================
echo Done