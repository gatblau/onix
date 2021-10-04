# Pilot Host Control Service

This section contains scripts to start up all required PilotCtl services:

- PostgreSQL database server
- Onix database
- PilotCtl database
- Onix Web API Service
- PilotCtl HTTP Service
- DbMan instances for database provisioning

![setup](pilotctl_setup.png)


### Helper Scripts

Helper scripts are available - please use these rather than direct Docker commands as they may do other things outside of Docker to help configure your stack:

script | description
--- | ---
deploy.sh | deploys a completely new control plane from scratch
stop.sh | stops the stack but keeps the containers and persistent data
start.sh | re-starts an existing stack
destroy.sh | destroys the stack completely, including any persistent data

### Notes:
- Prior to starting the stack, it is expected that a Nexus server is available to the Docker network, and is configured with a "raw (hosted)" repository named "artisan". Credentials for this repo should be updated into the `.env` file

- An attachable Docker network should already be available (this is done for you by the startup script if it doesn't exist).

- By default, all databases are *not* exposed outside of the Docker network to your host. If you want to access them please see the optional `docker-compose-database.yaml` file

### Endpoints

| service | endpoint |
|---|---|
| onix database | postgresql://onix:PWD@localhost:5432/onix |
| rem database | postgresql://rem:PWD@localhost:5432/rem |
| onix WAPI | http://localhost:8080 |
| rem WAPI | http://localhost:8088 |

### Configuration Variables

Are held in the file [here](.env)
