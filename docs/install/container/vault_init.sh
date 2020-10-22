# source vault token
. ~/vault/vault-env

sh ~/vault/vault-unseal
sh ~/vault/vault-plugins

# creates a namespace for password generation and binds secrets generation plugin
vault secrets enable -path="gen" -plugin-name="secrets-gen" plugin
