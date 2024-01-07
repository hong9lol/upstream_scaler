import threading
import time
from kubernetes import client, config

# Kubernetes API 클라이언트 생성
config.load_kube_config()  # 또는 config.load_incluster_config() 사용
api_client = client.ApiClient()
apps_v1 = client.AppsV1Api(api_client)
autoscaling_v2 = client.AutoscalingV2Api(api_client)
namespace = "default"  # 원하는 네임스페이스로 수정


def get_node_list():
    core_v1 = client.CoreV1Api(api_client)
    nodes = core_v1.list_node()
    result = list()
    # HPA 목록 출력
    for node in nodes.items:
        d = dict()
        d["name"] = node.metadata.name
        l = list()
        for address in node.status.addresses:
            _d = dict()
            _d["address"] = address.address
            _d["type"] = address.type
            l.append(_d)

        d["addresses"] = l
        result.append(d)
    return result


def get_deployment_list():
    # 모든 Deployments를 가져오기
    deployments = apps_v1.list_namespaced_deployment(namespace)
    result = list()
    # HPA 목록 출력
    for deployment in deployments.items:
        d = dict()
        d["name"] = deployment.metadata.name
        d["replicas"] = deployment.spec.replicas
        result.append(d)
    return result


def get_deployment(deployment_name):
    deployment = apps_v1.read_namespaced_deployment(name=deployment_name, namespace=namespace)
    d = dict()
    # print(deployment)
    d["name"] = deployment.metadata.name
    d["replicas"] = deployment.status.available_replicas

    return d


def get_hpa_list():
    hpas = autoscaling_v2.list_namespaced_horizontal_pod_autoscaler(namespace)
    result = list()
    # HPA 목록 출력
    for hpa in hpas.items:
        d = dict()
        d["name"] = hpa.metadata.name
        d["namespace"] = hpa.metadata.namespace
        d["min_replicas"] = hpa.spec.min_replicas
        d["max_replicas"] = hpa.spec.max_replicas
        d["target"] = hpa.spec.scale_target_ref.name
        l = list()
        for metric in hpa.spec.metrics:
            _d = dict()
            _d["name"] = metric.resource.name
            _d["target_utilization"] = metric.resource.target.average_utilization
            _d["type"] = metric.resource.target.type
            l.append(_d)

        d["metrics"] = l
        result.append(d)

    return result


def update_min_replica(deployment_name, waiting_time):
    time.sleep(waiting_time)  # 대기 시간 동안 일시 중지
    print("-1 replica")
    deployment = apps_v1.read_namespaced_deployment(name=deployment_name, namespace=namespace)
    if deployment.spec.replicas > 1:
        deployment.spec.replicas = deployment.spec.replicas - 1
        apps_v1.replace_namespaced_deployment(name=deployment_name, namespace=namespace, body=deployment)


def set_replica(deployment_name, replica_count):
    try:
        # Deployment 조회
        deployment = apps_v1.read_namespaced_deployment(name=deployment_name, namespace=namespace)
        # 새로운 레플리카 개수 설정
        deployment.spec.replicas = replica_count
        # Deployment 업데이트
        apps_v1.replace_namespaced_deployment(name=deployment_name, namespace=namespace, body=deployment)
        # -- min replica after a min.. (it is not work if we handle deployment not hpa)
        # threading.Thread(target=update_min_replica, args=(deployment_name, 60))
    except Exception as e:
        print(f"Fail update deployment: {str(e)}")
