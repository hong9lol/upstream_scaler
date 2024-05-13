#!/bin/bash

echo ====== Start ======

sudo sysctl -w fs.inotify.max_user_watches=2099999999
sudo sysctl -w fs.inotify.max_user_instances=2099999999
sudo sysctl -w fs.inotify.max_queued_events=2099999999

echo 1. Delete Previous Environment and Create new Environment
kind delete cluster 
kind create cluster --config ./yaml/kind/kind_1s_cadvisor.yaml


echo 2. Install application
echo 2-1. Install MetalLB
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.7/config/manifests/metallb-native.yaml
kubectl wait --namespace metallb-system \
                --for=condition=ready pod \
                --selector=app=metallb \
                --timeout=90s

echo 2-2. IP address pool for service loadbalancing
baseIP=$(docker network inspect -f '{{.IPAM.Config}}' kind | awk '/[[:space:]]/ {print $1}' | grep -oE '[0-9]+\.[0-9]+\.[0-9]')
export KIND_IP_RANGE="${baseIP}.200-${baseIP}.250"
echo $KIND_IP_RANGE
envsubst < yaml/metal-lb/metal-lb.yaml | kubectl apply -f-
sleep 10

kind load docker-image hong9lol/custom-metrics-server:v0.0.1
kind load docker-image hong9lol/simple_app:latest

echo 3. Start Metrics-server
kubectl delete -n kube-system deployments.apps metrics-server
# 0부터 30 사이의 랜덤 값을 생성
RANDOM_NUMBER=$((RANDOM % 10))

# 생성된 랜덤 값을 이용해 sleep 호출
echo "Sleeping for $RANDOM_NUMBER seconds..."
sleep $RANDOM_NUMBER   
kubectl apply -f yaml/metrics-server/metrics-server_1s.yaml
sleep 60

echo 4. Start Simple Servers
kubectl apply -f yaml/simple_server/simple_server.yaml
kubectl expose deployment simple-server --type=LoadBalancer --port=8080 & 
# 0부터 30 사이의 랜덤 값을 생성
RANDOM_NUMBER=$((RANDOM % 60))

# 생성된 랜덤 값을 이용해 sleep 호출
echo "Sleeping for $RANDOM_NUMBER seconds..."
sleep $RANDOM_NUMBER 
sleep 60

kubectl apply -f yaml/simple_server/simple_server_hpa.yaml
# 0부터 30 사이의 랜덤 값을 생성
RANDOM_NUMBER=$((RANDOM % 15))

# 생성된 랜덤 값을 이용해 sleep 호출
echo "Sleeping for $RANDOM_NUMBER seconds..."
sleep $RANDOM_NUMBER   

# set target IP
kubectl get svc | grep simple-server |  awk '/[[:space:]]/ {print $4}' > target_url.txt

./run_benchmark_simple1.sh
# sleep 15
# set target IP
# ./run_benchmark_simple2.sh

echo ====== Done ======