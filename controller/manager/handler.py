import threading
import time
from db import hpa
from kube_client import api

job_list = []


def hpa_updater():
    # HPA data update (interval 10s)
    interval = 10
    while True:
        time.sleep(interval)
        hpa_list = api.get_hpa_list()
        hpa.update_hpas(hpa_list)


def job_handler():
    job = job_dequeue()

    # find deployment name

    # reqeust all resource from agent
    #collector.collect_all_resource_usage_of_deployment("deployment_a")

    # get hpa condition
    # hpa.get_hpa(deployment_name)

    # if it is needed, do scale changeing min pods

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
