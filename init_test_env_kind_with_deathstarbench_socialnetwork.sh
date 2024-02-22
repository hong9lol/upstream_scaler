#!/bin/sh

echo ====== Start ======

# solution of "failed to create fsnotify watcher: too many open files"
sudo sysctl -w fs.inotify.max_user_watches=2099999999
sudo sysctl -w fs.inotify.max_user_instances=2099999999
sudo sysctl -w fs.inotify.max_queued_events=2099999999

# echo 1. Delete Previous Environment and Create new Environment
# kind delete cluster --name my-cluster
# kind create cluster --name my-cluster --config ./yaml/kind/kind.yaml

# echo 2. Install application
# echo 2-1. Install MetalLB
# kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.7/config/manifests/metallb-native.yaml
# kubectl wait --namespace metallb-system \
#                 --for=condition=ready pod \
#                 --selector=app=metallb \
#                 --timeout=90s

# echo 2-2. IP address pool for service loadbalancing
# baseIP=$(docker network inspect -f '{{.IPAM.Config}}' kind | awk '/[[:space:]]/ {print $1}' | grep -oE '[0-9]+\.[0-9]+\.[0-9]')
# export KIND_IP_RANGE="${baseIP}.200-${baseIP}.250"
# echo $KIND_IP_RANGE
# envsubst < yaml/metal-lb/metal-lb.yaml | kubectl apply -f-
# sleep 10

echo 2-3. Install helm packages
kubectl create secret docker-registry secret-jake --docker-username=hong9lol --docker-password=dlwoghd12@
kubectl create secret docker-registry secret-jake --docker-username=hong9lol --docker-password=dlwoghd12@ -n upstream-system
kubectl delete horizontalpodautoscalers.autoscaling --all=true --now=true --wait=true
helm uninstall social-network --wait
helm install social-network --wait ./DeathStarBench/socialNetwork/helm-chart/socialnetwork/

echo 3. Start Metrics-server
kubectl delete -n kube-system deployments.apps metrics-server
kubectl apply -f yaml/metrics-server/metrics-server.yaml
sleep 10

if [ "$1" = "default" ]; then
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
sleep 60

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

echo ====== Done ======