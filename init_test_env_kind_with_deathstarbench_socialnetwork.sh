#!/bin/sh

echo ====== Start ======
echo 1. Delete Previous Environment and Create new Environment
kind delete cluster --name my-cluster
kind create cluster --name my-cluster --config ./yaml/kind/kind.yaml

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

echo 2-3. Install helm packages
kubectl delete horizontalpodautoscalers.autoscaling --all=true --now=true --wait=true
helm uninstall social-network --wait
helm install social-network --wait ./DeathStarBench/socialNetwork/helm-chart/socialnetwork/

echo 3. Start Metrics-server
kubectl delete -n kube-system deployments.apps metrics-server
kubectl apply -f yaml/metrics-server/metrics-server.yaml
sleep 10

echo 4. Upstream controller and agents
cd ./controller/scripts
./run.sh k &
cd -
sleep 20
cd ./agent/scripts

./run.sh k &
cd -
sleep 20

echo 5. Auth for using Kubelet API agent
kubectl apply -f yaml/kubelet_auth/service_account.yaml
kubectl apply -f yaml/kubelet_auth/cluster_role_binding_auth.yaml

echo 6. Run Log
cd ./DeathStarBench/socialNetwork/benchmark_scripts
baseLogPath=./log
currentTime=`date +"%m-%d_%H%M%S"`
mkdir $baseLogPath/$currentTime
logPath=$baseLogPathPath/$currentTime
./log.sh $logPath & log=$!
cd -

kill -9 $log
echo ====== Done ======