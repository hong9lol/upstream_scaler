#!/bin/sh
#test environment in minikube with simple api servers

echo ====== Start ======
echo 1. Delete Previous Environment and Init Minikube 
minikube delete
minikube start --nodes=3 --cpus=max
# minikube mount /sys/fs/cgroup:/host/sys/fs/cgroup

echo 2. Start Simple Servers
kubectl apply -f yaml/simple_server/simple_server.yaml
kubectl apply -f yaml/simple_server/simple_server2.yaml
sleep 10
kubectl expose deployment simple-server --type=LoadBalancer --port=8080 & 
kubectl expose deployment simple-server2 --type=LoadBalancer --port=8080 & 

sleep 20
minikube service simple-server & 
minikube service simple-server2 &
#minikube service list -n default -o json | jq '.[1].URLs[0]' > target_url.txt

echo 3. Start Metrics-server
kubectl delete -n kube-system deployments.apps metrics-server
minikube addons enable metrics-server
# image=`kubectl get deployments.apps -n kube-system metrics-server -o wide --no-headers | awk '/[[:space:]]/ {print $7}'`
# echo $image
# kubectl patch -n kube-system deployments.apps metrics-server -p '{"spec":{"template":{"spec":{"containers":[{"name":"metrics-server", "image": "'${image}'", "args":["--cert-dir=/tmp","--secure-port=4443","--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname","--kubelet-use-node-status-port","--metric-resolution=15s","--kubelet-insecure-tls"]}]}}}}'
# sleep 10
# kubectl rollout restart -n kube-system deployments.apps metrics-server
# kubectl get deployments.apps -n kube-system metrics-server --template='{{range $k := .spec.template.spec.containers}}{{$k.args}}{{"\n"}}{{end}}' | grep -o 'metric-resolution=[^ ]*'
sleep 60

echo 4. Start HPA
kubectl apply -f yaml/simple_server/simple_server_hpa.yaml
kubectl apply -f yaml/simple_server/simple_server_hpa2.yaml

echo 5. Kubelet Auth
kubectl apply -f yaml/kubelet_auth/server_account.yaml
kubectl apply -f yaml/kubelet_auth/cluster_role_binding_auth.yaml

echo ====== Done ======