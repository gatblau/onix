# Deploy Onix containers without orchestration <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

***NOTE***: These scripts as they stand, are configured to only use ephemeral storage and are  for test and demo purposes only.

## Using Docker Compose

To spin up Onix using [Docker Compose](https://docs.docker.com/compose/):

- copy the [.env](.env), the [docker-compose.yaml](docker-compose.yaml) and [up_with_compose.sh](up_with_compose.sh) files to a local folder
- open the .env file and update the following password variables:
    - PG_ADMIN_PWD: the postgres admin database user password
    - ONIX_DB_PWD: the onix database user password
    - DBMAN_HTTP_PWD: the DbMan Web API HTTP admin user password
    - ONIX_HTTP_ADMIN_PWD: the Onix Web API HTTP admin user password
    - ONIX_HTTP_READER_PWD: the Onix Web API HTTP reader user password
    - ONIX_HTTP_WRITER_PWD: the Onix Web API HTTP writer user password
- finally run the up_with_compose.sh script, which will ensure schemas are installed and passwords updated accordingly
  
```bash
sh up_with_compose.sh
```

To destroy the containers:

```bash
docker-conpose down
```

## Using bare Docker

If you do not want to install [Docker Compose](https://docs.docker.com/compose/), then:

```bash
sh up_without_compose.sh
```

To destroy the containers:

```bash
docker rm -f oxdb
docker rm -f ox
docker rm -f oxku
docker rm -f dbman
```