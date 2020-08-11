# downloads prometheus
curl -L https://github.com/prometheus/prometheus/releases/download/v2.20.0/prometheus-2.20.0.darwin-amd64.tar.gz | tar -zx --strip-components=1 prometheus-2.20.0.darwin-amd64/prometheus
# downloads alert-manager
curl -L https://github.com/prometheus/alertmanager/releases/download/v0.21.0/alertmanager-0.21.0.darwin-amd64.tar.gz | tar -zx --strip-components=1 alertmanager-0.21.0.darwin-amd64/alertmanager
# downloads etcd
curl -L https://github.com/etcd-io/etcd/releases/download/v3.4.10/etcd-v3.4.10-darwin-amd64.zip -o etcd-v3.4.10-darwin-amd64.zip && unzip -qq etcd-v3.4.10-darwin-amd64.zip && mv etcd-v3.4.10-darwin-amd64/etcd ./ && rm -rf etcd-v3.4.10-darwin-amd64 && rm etcd-v3.4.10-darwin-amd64.zip


