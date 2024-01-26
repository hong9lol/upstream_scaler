import threading, time, logging

from kubernetes import client, config

# Kubernetes API 클라이언트 생성
try:
    config.load_kube_config()  # 또는 config.load_incluster_config() 사용
except:
    logging.info("config for in-cluster")
    config.load_incluster_config()

api_client = client.ApiClient()
core_v1 = client.CoreV1Api()
apps_v1 = client.AppsV1Api(api_client)
autoscaling_v2 = client.AutoscalingV2Api(api_client)
namespace = "default"  # 원하는 네임스페이스로 수정


def get_agent_list(_namespace):
    agents = []
    try:
        # 파드 목록 가져오기
        pod_list = core_v1.list_namespaced_pod(_namespace)
        for pod in pod_list.items:
            if "controller" in pod.metadata.name:
                continue
            p = dict()
            p["name"] = pod.metadata.name
            p["pod_ip"] = pod.status.pod_ip
            p["host_ip"] = pod.status.host_ip
            agents.append(p)
            # 파드 목록 출력
        logging.info("파드 목록:")
        for pod in pod_list.items:
            logging.info(f"Pod Name: {pod.metadata.name}, IP: {pod.status.podIP}")
        return agents
    except Exception as e:
        logging.error("Get agent list fail" + str(e))
    finally:
        return agents


def get_node_list():
    core_v1 = client.CoreV1Api(api_client)
    node_list_raw = core_v1.list_node()
    node_list = list()

    for node in node_list_raw.items:
        _node = dict()
        _node["name"] = node.metadata.name
        node_addr_list = list()
        for address in node.status.addresses:
            _addr = dict()
            _addr["address"] = address.address
            _addr["type"] = address.type
            node_addr_list.append(_addr)

        _node["addresses"] = node_addr_list
        node_list.append(_node)
    return node_list


def get_deployment_list():
    # 모든 Deployments를 가져오기
    deployment_list_raw = apps_v1.list_namespaced_deployment(namespace)
    deployment_list = list()

    for deployment in deployment_list_raw.items:
        _deployment = dict()
        _deployment["name"] = deployment.metadata.name
        _deployment["replicas"] = deployment.spec.replicas
        deployment_list.append(_deployment)
    return deployment_list


def get_deployment(deployment_name):
    deployment_raw = apps_v1.read_namespaced_deployment(
        name=deployment_name, namespace=namespace
    )
    deployment = dict()
    deployment["name"] = deployment_raw.metadata.name
    deployment["replicas"] = deployment_raw.status.available_replicas

    return deployment


def get_hpa_list():
    hpa_list_raw = autoscaling_v2.list_namespaced_horizontal_pod_autoscaler(namespace)
    hpa_list = list()
    # HPA 목록 출력
    for hpa in hpa_list_raw.items:
        _hpa = dict()
        _hpa["name"] = hpa.metadata.name
        _hpa["namespace"] = hpa.metadata.namespace
        _hpa["min_replicas"] = hpa.spec.min_replicas
        _hpa["max_replicas"] = hpa.spec.max_replicas
        _hpa["target"] = hpa.spec.scale_target_ref.name
        metric_list = list()
        for metric in hpa.spec.metrics:
            _metric = dict()
            _metric["name"] = metric.resource.name
            _metric["target_utilization"] = metric.resource.target.average_utilization
            _metric["type"] = metric.resource.target.type
            metric_list.append(_metric)

        _hpa["metrics"] = metric_list
        hpa_list.append(_hpa)

    return hpa_list


def reset_min_replica(hpa_name, waiting_time, init_min_replicas):
    time.sleep(waiting_time)  # 대기 시간 동안 일시 중지
    change_hpa_min_replicas(hpa_name, init_min_replicas)


def change_hpa_min_replicas(hpa_name, new_min_replicas):
    # HPA 객체 가져오기
    hpa_raw = autoscaling_v2.read_namespaced_horizontal_pod_autoscaler(
        hpa_name, namespace
    )

    prev_min_replicas = hpa_raw.spec.min_replicas
    # 최소 복제 수 변경
    hpa_raw.spec.min_replicas = new_min_replicas

    # 변경 사항 적용
    autoscaling_v2.replace_namespaced_horizontal_pod_autoscaler(
        hpa_name, namespace, hpa_raw
    )
    logging.info("change {hpa_name}'s min replicas to {new_min_replicas}")

    return prev_min_replicas


def set_replica(deployment_name, replica_count, hpa_name):
    try:
        # Deployment 조회
        deployment = apps_v1.read_namespaced_deployment(
            name=deployment_name, namespace=namespace
        )
        # 새로운 레플리카 개수 설정
        deployment.spec.replicas = replica_count
        # Deployment 업데이트
        apps_v1.replace_namespaced_deployment(
            name=deployment_name, namespace=namespace, body=deployment
        )  # , field_validation="PATCH")
        init_min_replicas = change_hpa_min_replicas(hpa_name, replica_count)
        threading.Thread(
            target=reset_min_replica, args=(hpa_name, init_min_replicas, 60)
        )
        logging.info(f"Scale out [{deployment_name}] replicas: {replica_count}")

    except Exception as e:
        logging.info(f"Fail update deployment: {str(e)}")
