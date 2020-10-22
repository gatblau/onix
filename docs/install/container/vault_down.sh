echo ==============================================================================
echo Stopping and destroying any existing Vault ...
docker stop vault
docker rm vault
docker volume rm vault-config
docker volume rm vault-data
docker volume rm vault-logs
docker volume rm vault-policies