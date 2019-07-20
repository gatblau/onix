# Deploy Onix containers without orchestration <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

***NOTE***: These scripts as they stand, are configured to only use ephemeral storage and are  for test and demo purposes only.

## Using Docker Compose

To spin up Onix using [Docker Compose](https://docs.docker.com/compose/) just run the following command:

```bash
sh up_with_compose.sh
```

***NOTE***: The above script will call compose but also the readyness probe to deploy the database.

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
```