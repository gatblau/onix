# Rem Backend

This section contains scripts to start up all required REM backend services:

- PostgreSQL database server
- Onix database
- REM database
- Onix Web API Service
- REM HTTP Service
- DbMan instances for database provisioning

![setup](rem_setup.png)

### Start services

In a host with docker and docker-compose installed:

```bash
$ sh up.sh
```

### Stopping services

In a host with docker and docker-compose installed:

```bash
$ sh down.sh
```

### Endpoints

| service | endpoint |
|---|---|
| onix database | postgresql://onix:PWD@localhost:5432/onix |
| rem database | postgresql://rem:PWD@localhost:5432/rem |
| onix WAPI | http://localhost:8080 |
| rem WAPI | http://localhost:8088 |

### Configuration Variables

Are held in the file [here](.env)


