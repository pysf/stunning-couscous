#!/usr/bin/env bash

export POSTGRESQL_PASSWORD=gHteuivwdvkew4wt
export POSTGRESQL_USERNAME=postgres
export POSTGRESQL_PORT=5432
export POSTGRESQL_DATABASE=stunning-couscous
export POSTGRESQL_HOST=127.0.0.1

go test ./... -v