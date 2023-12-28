#!/bin/sh

echo ====== Start ======

echo 1. Init Minikube
minikube delete
minikube start --nodes=3 --cpus=max --cni=calico

echo 2. Start Simple Application
kubectl apply -f yaml/simple_app.yaml
sleep 10
kubectl expose deployment simple-app-deployment --type=LoadBalancer --port=8080 & 

sleep 10
minikube service simple-app-deployment
minikube service list -n default -o json | jq '.[1].URLs[0]' > target_url.txt

echo 3. Start Metrics-server
kubectl delete -n kube-system deployments.apps metrics-server
minikube addons enable metrics-server
# image=`kubectl get deployments.apps -n kube-system metrics-server -o wide --no-headers | awk '/[[:space:]]/ {print $7}'`
# echo $image
# kubectl patch -n kube-system deployments.apps metrics-server -p '{"spec":{"template":{"spec":{"containers":[{"name":"metrics-server", "image": "'${image}'", "args":["--cert-dir=/tmp","--secure-port=4443","--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname","--kubelet-use-node-status-port","--metric-resolution=15s","--kubelet-insecure-tls"]}]}}}}'
# sleep 10
# kubectl rollout restart -n kube-system deployments.apps metrics-server
# kubectl get deployments.apps -n kube-system metrics-server --template='{{range $k := .spec.template.spec.containers}}{{$k.args}}{{"\n"}}{{end}}' | grep -o 'metric-resolution=[^ ]*'
sleep 120

echo ====== Done ======