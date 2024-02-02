#!/bin/sh

echo ====== Start ======
echo 1. Delete Previous Environment and Create new Environment
kind delete cluster --name my-cluster
kind create cluster --name my-cluster --config ./yaml/kind/kind.yaml


echo 2. Start Simple Servers
kubectl apply -f yaml/simple_server/simple_server.yaml
kubectl apply -f yaml/simple_server/simple_server2.yaml
sleep 10
kubectl expose deployment simple-server --type=LoadBalancer --port=8080 & 
kubectl expose deployment simple-server2 --type=LoadBalancer --port=8080 & 
sleep 10

echo 3. Start Metrics-server
kubectl delete -n kube-system deployments.apps metrics-server
kubectl apply -f yaml/metrics-server/metrics-server.yaml
sleep 60

echo 4. Start HPA
kubectl apply -f yaml/simple_server/simple_server_hpa.yaml
kubectl apply -f yaml/simple_server/simple_server_hpa2.yaml

echo ====== Done ======