import threading
import logging
import requests

from utils import regex

from db import agent as agent_db
lock = threading.Lock()
data = []


def requester(url, deployment):
    # for test
    # url = "127.0.0.1:3001"

    response = requests.get(
        'http://' + url + ":3001/api/v1/metrics/" + deployment)
    logging.warning(f"[Collect resource data in node({url})")

    lock.acquire()  # 작업이 끝나기 전까지 다른 쓰레드가 공유데이터 접근을 금지
    data.append(response.json())
    lock.release()  # lock 해제


# need to be thread safe... don't call multiple times at the same time
def collect_all_resource_usage_of_deployment(deployment):
    data.clear()
    threads = []
    agents = agent_db.get_all_agents()
    for agent in agents:
        _agent = agents[agent]
        _ip = _agent["pod_ip"]
        if regex.ip_validation_check(_ip):
            threads.append(threading.Thread(
                target=requester, args=(_ip, deployment)))
    # nodes = client.get_node_list()
    # for node in nodes:
    #     for ip in node["addresses"]:
    #         if regex.ip_validation_check(ip["address"]):
    #             threads.append(threading.Thread(
    #                 target=requester, args=(ip["address"], deployment)))

    for thread in threads:
        thread.start()

    for thread in threads:
        thread.join()

    return data
