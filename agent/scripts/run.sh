#!/bin/sh

if [ "$1" = "k" ]; then
    kubectl create namespace upstream-system
    kubectl delete daemonsets.apps -n upstream-system upstream-agent
    kubectl apply -f ../k8s/agent.yaml
else
    ../build/output/agent
fi