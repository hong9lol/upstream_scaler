#!/bin/sh

if [ "$1" = "k" ]; then
    kubectl create namespace upstream-system
    kubectl delete -n upstream-system deployments.apps upstream-controller
    kubectl delete -n upstream-system clusterrole upstream-cluster-role
    kubectl delete -n upstream-system clusterrolebindings.rbac.authorization.k8s.io upstream-cluster-rolebinding
    kubectl apply -f ../k8s/cluster-role-binding.yaml
    kubectl apply -f ../k8s/controller.yaml

    # if minikube
    minikube service -n upstream-system upstream-controller
else
    python3 ../main.py
fi