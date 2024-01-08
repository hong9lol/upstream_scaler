#!/bin/sh

if [ "$1" = "k" ]; then
    kubectl create namespace upstream-system
    kubectl delete -n upstream-system daemonsets.apps upstream-agent
    kubectl apply -f ../k8s/agent.yaml
else
    ../build/output/agent
fi