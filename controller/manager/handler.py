import logging
import threading
import time
import math

from db import hpa as hpa_db, agent as agent_db
from kube_client import client
from manager import collector

job_list = []
wait_list = []


def hpa_updater():
    # HPA data update (interval 10s)
    interval = 10
    while True:
        hpa_list = client.get_hpa_list()
        hpa_db.update_hpas(hpa_list)
        time.sleep(interval)


# TODO: update agent list with event based logic
def agent_updater():
    interval = 60  # every 1 min
    while True:
        agent_list = client.get_agent_list("upstream-system")
        agent_db.update_agents(agent_list)
        time.sleep(interval)


# after scale out, deployment need to wait a few second next handling
def remove_from_wait_list(deployment_name, waiting_time):
    time.sleep(waiting_time)  # 대기 시간 동안 일시 중지
    wait_list.remove(deployment_name)


def do_scale(deployment_name, current_cpu_usage_rate, target_cpu_utilization):
    # TODO: get the deployment everytime, make it more efficient
    deployment = client.get_deployment(deployment_name)
    replica_count = math.ceil(
        deployment["replicas"] * (current_cpu_usage_rate/target_cpu_utilization))
    client.set_replica(deployment_name, replica_count)
    wait_list.append(deployment_name)
    threading.Thread(target=remove_from_wait_list,
                     args=(deployment_name, 10)).start()
    print("Scale out " + deployment_name + ", replicas:" + str(replica_count))


def get_cpu_usage_rate_per_sec(pods):

    for pod in pods:
        _pod = pods[pod]
        pod_cpu_usage_per_sec = 0.0
        containers = _pod["containers"]
        for container in containers:
            container_cpu_usage_per_sec = 0.0
            _container = containers[container]
            last = _container["usages"][len(_container) - 1]
            prev = _container["usages"][len(_container) - 2]
            usage = last["usage"] - prev["usage"]
            timestamp = last["timestamp"] - prev["timestamp"]
            container_cpu_usage_per_sec += float(usage) / float(timestamp)
            pod_cpu_usage_per_sec += (container_cpu_usage_per_sec /
                                      _container["cpu_request"]) * 100
            # logging.warning(container_cpu_usage_per_sec,
            #              _container["cpu_request"])
        pod_cpu_usage_per_sec = pod_cpu_usage_per_sec / len(containers)
    total_cpu_usage_per_sec = pod_cpu_usage_per_sec / len(pods)

    return total_cpu_usage_per_sec


def job_handler():
    interval = 1
    # if it is needed, do scale changing min pods
    while True:
        if is_empty():
            time.sleep(interval)
            continue

        job = job_dequeue()
        deployment_name = job["deployment_name"]
        if wait_list.count(deployment_name) > 0:
            continue
        hpa_name = job["hpa_name"]
        data = collector.collect_all_resource_usage_of_deployment(
            deployment_name)
        _hpa = hpa_db.get_hpa(hpa_name)

        # 알고리즘
        usage_rate = 0.0
        involved_nodes = 0
        for d in data:
            if "pods" not in d:
                continue
            if d["pods"] == None:
                continue
            involved_nodes += 1
            usage_rate += get_cpu_usage_rate_per_sec(d["pods"])
        if involved_nodes < 1:
            continue

        target_cpu_utilization = 100
        metrics = _hpa["metrics"]
        for metric in metrics:
            if metric["name"] == "cpu":
                target_cpu_utilization = metric["target_utilization"]
                break
        current_cpu_usage_rate = usage_rate / involved_nodes
        logging.warning(
            f"usage_rate: {usage_rate}, involved_nodes: {involved_nodes}, current_cpu_usage_rate: {current_cpu_usage_rate}, target_cpu_utilization: {target_cpu_utilization}")
        if current_cpu_usage_rate > target_cpu_utilization:
            do_scale(deployment_name, current_cpu_usage_rate,
                     target_cpu_utilization)


def start():
    # 매 10초마다 hpa 정보 가져오기
    hpa_updater_thread = threading.Thread(target=hpa_updater)
    hpa_updater_thread.start()

    # 매 10초마다 hpa 정보 가져오기
    agent_updater_thread = threading.Thread(target=agent_updater)
    agent_updater_thread.start()

    # Job handler
    job_handler_thread = threading.Thread(target=job_handler)
    job_handler_thread.start()
    # collector.collect_all_resource_usage_of_deployment("deployment_a")
    # hpa_updater_thread.join()


def job_enqueue(job):
    if job_list.count(job) > 0:
        print("This job is already handling, skip it now")
        return

    job_list.insert(0, job)


def job_dequeue():
    return job_list.pop()


def is_empty():
    if len(job_list) < 1:
        return True
    else:
        return False
