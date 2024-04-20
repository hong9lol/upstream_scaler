#!/bin/sh

echo ====== Start ======

# solution of "failed to create fsnotify watcher: too many open files"
sudo sysctl -w fs.inotify.max_user_watches=2099999999
sudo sysctl -w fs.inotify.max_user_instances=2099999999
sudo sysctl -w fs.inotify.max_queued_events=2099999999

echo 1. Delete Previous Environment and Create new Environment
#if [ "$1" = "fast" ]; then
#    kind create cluster --name my-cluster --config ./yaml/kind/kind_1s.yaml
#else
#    kind create cluster --name my-cluster --config ./yaml/kind/kind.yaml
#fi


echo 2. Install application
#echo 2-1. Install MetalLB
#kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.7/config/manifests/metallb-native.yaml
#kubectl wait --namespace metallb-system \
#                --for=condition=ready pod \
#                --selector=app=metallb \
#                --timeout=90s

echo 2-2. IP address pool for service loadbalancing
kubectl apply -f yaml/metal-lb/routing.yaml 
#baseIP=$(docker network inspect -f '{{.IPAM.Config}}' kind | awk '/[[:space:]]/ {print $1}' | grep -oE '[0-9]+\.[0-9]+\.[0-9]')
#export KIND_IP_RANGE="${baseIP}.200-${baseIP}.250"
#echo $KIND_IP_RANGE
#envsubst < yaml/metal-lb/metal-lb.yaml | kubectl apply -f-
#sleep 10

echo 2-3. Install helm packages
#kind load docker-image mongo:4.4.6 -n my-cluster # move local image to kind cluster nodes
#kind load docker-image redis:6.2.4 -n my-cluster
#kind load docker-image memcached:1.6.7 -n my-cluster
#kind load docker-image alpine/git:2.43.0 -n my-cluster
#kind load docker-image yg397/openresty-thrift:xenial -n my-cluster
#kind load docker-image yg397/media-frontend:xenial -n my-cluster
#kind load docker-image jaegertracing/all-in-one:1 -n my-cluster
#kind load docker-image deathstarbench/social-network-microservices:0.3.0 -n my-cluster

kubectl create secret docker-registry secret-jake --docker-username=hong9lol --docker-password=dlwoghd12@ 
kubectl create secret docker-registry secret-jake --docker-username=hong9lol --docker-password=dlwoghd12@ -n upstream-system

kubectl delete horizontalpodautoscalers.autoscaling --all=true --now=true --wait=true
helm uninstall social-network --wait
date +"%T"
sleep 330
echo install application
helm install social-network --wait ./DeathStarBench/socialNetwork/helm-chart/socialnetwork/

echo 3. Start Metrics-server
kubectl delete -n kube-system deployments.apps metrics-server
sleep 10
if [ "$1" = "fast" ]; then
    kubectl apply -f yaml/metrics-server/metrics-server_15s.yaml
else
    kubectl apply -f yaml/metrics-server/metrics-server_60s.yaml
fi
sleep 90

kubectl delete deployments.apps -n upstream-system upstream-controller
kubectl delete daemonsets.apps -n upstream-system upstream-agent

if [ "$1" = "default" ]; then
    echo Skip Step 4, 5 for upstream scaler
elif [ "$1" = "fast" ]; then
    echo Skip Step 4, 5 for upstream scaler
else
    echo 4. Upstream controller and agents
    # build
    cd ./controller/scripts
    ./build.sh d
    cd -
    cd ./agent/scripts
    ./build.sh d
    cd -

    cd ./controller/scripts
    ./run.sh k &
    cd -
    cd ./agent/scripts
    ./run.sh k &
    cd -

    echo 5. Auth for using Kubelet API agent
    kubectl apply -f yaml/kubelet_auth/service_account.yaml
    kubectl apply -f yaml/kubelet_auth/cluster_role_binding_auth.yaml
fi
sleep 120

echo 6. Run Log
cd ./DeathStarBench/socialNetwork/benchmark_scripts
baseLogPath=./log
currentTime=`date +"%m-%d_%H%M%S"`
mkdir $baseLogPath/$currentTime
logPath=$baseLogPath/$currentTime
./log.sh $logPath & log=$!

echo 7. Run Benchmark
./run_social_benchmark.sh $logPath
cd -

echo 8. Kill Log Proc
kill -9 $log

# 파드 리스트 가져오기
./time_checker.sh $1
mv podcnt.txt ./DeathStarBench/socialNetwork/benchmark_scripts/$logPath/podcnt.txt

kubectl logs -n upstream-system deployments/upstream-controller >> controller.log
mv controller.log ./DeathStarBench/socialNetwork/benchmark_scripts/$logPath/controller.log

echo ====== Done ======
