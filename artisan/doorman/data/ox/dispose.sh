#!/bin/bash
# application 'Onix Config Manager' dispose script using docker-compose
# auto-generated by Artisan on 2022-06-01 15:52:30.100386 +0000 UTC


# bring down services
docker-compose down

# remove docker volumes
docker volume rm db