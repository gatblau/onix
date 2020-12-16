REGISTRY=https://index.docker.io/v2
#REGISTRY="https://registry.hub.docker.com/v2"
#REGISTRY="https://registry.docker.io/v2"
#REGISTRY="https://registry-1.docker.io/v2"
#REGISTRY="https://hub.docker.com/v2"

REPO=library
IMAGE=debian
# Could also be a repo digest
TAG=latest

# Query tags
#curl "$REGISTRY/repositories/$REPO/$IMAGE/tags/"

# Query manifest
#curl -iL "$REGISTRY/$REPO/$IMAGE/manifests/$TAG"
# HTTP/1.1 401 Unauthorized
# Www-Authenticate: Bearer realm="https://auth.docker.io/token",service="registry.docker.io",scope="repository:library/debian:pull"

TOKEN=$(curl -sSL "https://auth.docker.io/token?service=registry.docker.io&scope=repository:$REPO/$IMAGE:pull" \
  | jq --raw-output .token)
curl -LH "Authorization: Bearer ${TOKEN}" "$REGISTRY/$REPO/$IMAGE/manifests/$TAG"

# Some repos seem to return V1 Schemas by default

REPO=nginxinc
IMAGE=nginx-unprivileged
TAG=1.17.2

curl -LH "Authorization: Bearer $(curl -sSL "https://auth.docker.io/token?service=registry.docker.io&scope=repository:$REPO/$IMAGE:pull" | jq --raw-output .token)" \
 "$REGISTRY/$REPO/$IMAGE/manifests/$TAG"

# Solution: Set the Accept Header for V2

curl -LH "Authorization: Bearer $(curl -sSL "https://auth.docker.io/token?service=registry.docker.io&scope=repository:$REPO/$IMAGE:pull" | jq --raw-output .token)" \
  -H "Accept:application/vnd.docker.distribution.manifest.v2+json" \
 "$REGISTRY/$REPO/$IMAGE/manifests/$TAG"