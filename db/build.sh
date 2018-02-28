#!/usr/bin/env bash
# builds an onix postgresql database image using the S2I tool
./s2i build ./image_conf/ centos/postgresql-96-centos7 onix-db:0.0.1-0