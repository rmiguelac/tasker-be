#!/usr/bin/env bash

function log {
    echo "[$1] $2"
}

if [[ $1 == "create" ]]; then
    docker run -d \
    	--name postgres \
        --restart always \
    	-e POSTGRES_PASSWORD=testpwd \
    	-e PGDATA=/var/lib/postgresql/data/pgdata \
    	-v ${HOME}/repositories/tasker/db:/var/lib/postgresql/data \
    	-v ${HOME}/repositories/tasker/scripts:/docker-entrypoint-initdb.d \
    	postgres:14.3

    cat <<'EOF' > ${HOME}/repositories/tasker/scripts/init.sql
CREATE USER tasker;
CREATE DATABASE tasks;
GRANT ALL PRIVILEGES ON DATABASE tasks TO tasker;
CREATE TABLE tasks (id serial, name char(124));
EOF

elif [[ $1 == "delete" ]]; then
    sudo rm -Rf ${HOME}/repositories/tasker/db 
    docker stop postgres
    docker rm postgres
else
    log "ERROR" "No argument provided. Either 'create' or 'delete' should be given"

fi
