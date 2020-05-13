# DbMan - Onix Database Manager

DbMan is written in go and is meant to run in a container providing:
- Database and schemas creation
- Database schema upgrades
- Database backups
- Database restores
- any other required database administration tasks

**Note** DbMan is work in progress, it will replace and augment the logic in the Onix Web API
and allow a Kubernetes operator to orchestrate the required database admin tasks
