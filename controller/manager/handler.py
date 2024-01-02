import threading
import time
from db import hpa
from kube_client import api

job_list = []


def hpa_updater():
    # HPA data update (interval 10s)
    interval = 10
    while True:
        hpa_list = api.get_hpa_list()
        hpa.update_hpas(hpa_list)
        time.sleep(interval)


def job_handler():


    # find deployment name

    # reqeust all resource from agent
    #collector.collect_all_resource_usage_of_deployment("deployment_a")

    # get hpa condition
    # hpa.get_hpa(deployment_name)

    # if it is needed, do scale changeing min pods
    while True:
        if is_empty():
            time.sleep(1)
            continue

        job = job_dequeue()


def start():
    # 매 10초마다 hpa 정보 가져오기
    hpa_updater_thread = threading.Thread(target=hpa_updater)
    hpa_updater_thread.start()

    # Job handler
    job_handler_thread = threading.Thread(target=job_handler)
    job_handler_thread.start()
    # collector.collect_all_resource_usage_of_deployment("deployment_a")
    # hpa_updater_thread.join()


def job_enqueue(job):
    job_list.insert(0, job)
    print(job_list)


def job_dequeue():
    return job_list.pop()


def is_empty():
    if len(job_list) < 1:
        return True
    else:
        return False
