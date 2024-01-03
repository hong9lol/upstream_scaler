from kubernetes import client, config

# Kubernetes API 클라이언트 생성
config.load_kube_config()  # 또는 config.load_incluster_config() 사용
api_client = client.ApiClient()
apps_v1 = client.AppsV1Api(api_client)
autoscaling_v2 = client.AutoscalingV2Api(api_client)


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
    namespace = "default"  # 원하는 네임스페이스로 수정
    deployments = apps_v1.list_namespaced_deployment(namespace)
    result = list()
    # HPA 목록 출력
    for deployment in deployments.items:
        print(deployment.metadata.annotations["kubectl.kubernetes.io/last-applied-configuration"])
        d = dict()
        d["name"] = deployment.metadata.name
        d["replicas"] = deployment.spec.replicas
        result.append(d)
    return result


def get_hpa_list():
    namespace = "default"  # 원하는 네임스페이스로 수정
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
