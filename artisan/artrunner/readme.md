<img src="https://github.com/gatblau/artisan/raw/master/artisan.png" width="150" align="right"/>

# Artisan Kubernetes Runner

Execute Artisan flows as Tekton pipelines.

The runner can be configured with the following variables:

| var | description | default |
|---|---|---|
| `OX_METRICS_ENABLED` | enable prometheus /metrics endpoint | true |
| `OX_SWAGGER_ENABLED` | enable swagger user interface under /api/ | true |
| `OX_HTTP_PORT` | the port on which the server listen for connections | 8080 |
| `OX_HTTP_UNAME` | the basic authentication username | admin |
| `OX_HTTP_PWD` | the basic authentication password | adm1n |

## Image location

quay.io/gtablau/artisan-runner