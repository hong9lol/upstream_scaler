#!/bin/sh

if [ "$1" = "d" ]; then
    cd ../docker
    mkdir src
    
    # copy srouces
    cp ../main.py src/
    cp -r ../db src/
    cp -r ../kube_client src/
    cp -r ../manager src/
    cp -r ../utils src/
    cp ../requirements.txt src/
    
    # build docker & push 
    docker build -t upsteam_controller -f Dockerfile .
    docker image tag upsteam_controller:latest hong9lol/upstream-controller:0.1
    docker push hong9lol/upstream-controller:0.1
    
    rm -rf src

    # run
    cd -
    #./run.sh k
else 
    pip3 install -r ../requirements.txt
    ./run.sh
fi

