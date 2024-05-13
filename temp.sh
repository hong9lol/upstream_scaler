#!/bin/bash
minikube addons enable metrics-server
sleep 10
image=`kubectl get deployments.apps -n kube-system metrics-server -o wide --no-headers | awk '/[[:space:]]/ {print $7}'`
echo $image
kubectl patch -n kube-system deployments.apps metrics-server -p '{"spec":{"template":{"spec":{"containers":[{"name":"metrics-server", "image": "'${image}'", "args":["--cert-dir=/tmp","--secure-port=4443","--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname","--kubelet-use-node-status-port","--metric-resolution=1s","--kubelet-insecure-tls"]}]}}}}'
sleep 10
kubectl rollout restart -n kube-system deployments.apps metrics-server
kubectl get deployments.apps -n kube-system metrics-server --template='{{range $k := .spec.template.spec.containers}}{{$k.args}}{{"\n"}}{{end}}' | grep -o 'metric-resolution=[^ ]*'
sleep 60
k create namespace cadvisor
kubectl apply -f yaml/simple_server/simple_server.yaml
kubectl apply -f cadvisor.yaml