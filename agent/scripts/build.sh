#!/bin/sh


if [ "$1" = "d" ]; then
    GOOS=linux GOARCH=amd64 go build -o ../build/output/agent ./../cmd/agent.go 
    cp ../build/output/agent ../build/docker/agent
    cd ../build/docker
    docker build -t upsteam_agent -f Dockerfile .
    docker image tag upsteam_agent:latest hong9lol/upstream-agent:0.1
    docker push hong9lol/upstream-agent:0.1
    cd -
    #./run.sh k
else 
    go build -o ../build/output/agent ./../cmd/agent.go 
fi