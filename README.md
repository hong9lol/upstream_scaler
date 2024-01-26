# Upstream Scaler

## Environment

> Minikube

```
./init_test_env_minikube_with_simple_api_servers.sh
```

> GCP

```
TBD
```

## What it is for

This is a custom HPA(hotizontal pod auto-scaler). The Generic HPA incurs overhead due to frequent polling of contaienr resources. This overhead causes poor performance and reliability.
The goal of "Upstream Scaler" is to reduce unnecessary work during inactivity by replacing periodic polling.

## Agent

The Agent is a daemonset running on every node to monitor container resources which the user is interested in. This continuously supervises all container resources of Pods defined in the HPA operating within Node, and upon exceeding the threshold, it reports to the controller.

## Controller

The Controller is a deployment relaying between all Agents on each node. It is also an api client of the Kubernetes api-server. it collects resources from agents and sends scaling commands when necessary to api-server.

## How to run

> Agent

```bash
$ cd ./agent/scripts
$ ./build.sh d
$ ./run.sh k
```

> Controller

```bash
$ cd ./agent/scripts
$ ./build.sh d

# minikube
$ ./run.sh k m
# others
$ ./run.sh k
```
