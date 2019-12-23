#!/usr/bin/env bash

docker kill pg-db || :
mkdir -p $HOME/docker/volumes/postgres
docker run --rm --name pg-db \
  -v $HOME/docker/volumes/postgres:/var/lib/postgresql/data \
  -e POSTGRES_PASSWORD=develop -d -p 5432:5432 \
  postgres

# psql -h localhost -U postgres -d postgres
