# get authentication token
UNAME=gatblau
UPASS=Andr3s4561
TOKEN=$(curl -s -H "Content-Type: application/json" -X POST -d '{"username": "'${UNAME}'", "password": "'${UPASS}'"}' https://hub.docker.com/v2/users/login/ | ./jq -r .token)

#token=$(echo "$jsonToken" | ./jq '.access_token')

curl "https://hub.docker.com/v2/gatblau/manifests/onix-snapshot" \
  -H "Accept: application/vnd.docker.distribution.manifest.v2+json" \
  -H "Authorization: JWT ${TOKEN}" -v