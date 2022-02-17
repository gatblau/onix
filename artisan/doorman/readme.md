# Artisan's Doorman

A gateway for securely transferring and resigning application release artefacts between networks.

mc admin config get local/ notify_webhook
mc admin config set local/ notify_webhook:1 queue_limit="0"  endpoint="http://localhost:8080/event/minio" queue_dir=""