#!/usr/bin/env bash

function log {
    echo "[$1] $2"
}

if [[ $1 == "create" ]]; then
    docker run -d \
    	--name postgres \
        --restart always \
        -p 5432:5432 \
    	-e POSTGRES_PASSWORD=testpwd \
    	-e PGDATA=/var/lib/postgresql/data/pgdata \
    	-v ${HOME}/repositories/tasker/db:/var/lib/postgresql/data \
    	-v ${HOME}/repositories/tasker/schemas:/docker-entrypoint-initdb.d \
    	postgres:14.3

elif [[ $1 == "delete" ]]; then
    sudo rm -Rf ${HOME}/repositories/tasker/db 
    docker stop postgres
    docker rm postgres
else
    log "ERROR" "No argument provided. Either 'create' or 'delete' should be given"

fi
