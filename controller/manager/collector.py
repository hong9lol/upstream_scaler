import threading

import requests

from kube_client import api
from utils import regex

lock = threading.Lock()
data = []


def requester(url, deployment):
    # for test
    url = "127.0.0.1:3001"

    response = requests.get('http://' + url + "/api/v1/metrics/" + deployment)
    print(response.json())
    lock.acquire()  # 작업이 끝나기 전까지 다른 쓰레드가 공유데이터 접근을 금지
    data.append(response.json())
    lock.release()  # lock 해제


# need to be thread safe... don't call multiple times at the same time
def collect_all_resource_usage_of_deployment(deployment):
    data.clear()
    threads = []
    nodes = api.get_node_list()
    for node in nodes:
        for ip in node["addresses"]:
            if regex.ip_validation_check(ip["address"]):
                threads.append(threading.Thread(
                    target=requester, args=(ip["address"], deployment)))

    for thread in threads:
        thread.start()

    for thread in threads:
        thread.join()

    return data
