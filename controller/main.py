import json

from flask import Flask, request
from kube_client import client
from manager import handler
from db import hpa, agent
import logging

app = Flask(__name__)

logging.basicConfig(
    format="%(asctime)s %(levelname)s:%(message)s",
    level=logging.INFO,
    datefmt="%m/%d/%Y %I:%M:%S %p",
)


_port = 5001


@app.route("/api/v1/nodes", methods=["GET"])
def node():
    ret = client.get_node_list()
    logging.info(ret)
    return json.dumps(ret)


@app.route("/api/v1/deployments", methods=["GET"])
def get_deployment():
    ret = client.get_deployment_list()
    logging.info(ret)
    return json.dumps(ret)


@app.route("/api/v1/hpas", methods=["GET"])
def get_hpa():
    ret = hpa.get_all_hpas()
    if len(ret) != 0:
        return json.dumps(ret)
    else:
        return json.dumps([])


@app.route("/api/v1/agents", methods=["GET"])
def get_agent():
    ret = agent.get_all_agents()
    if len(ret) != 0:
        return json.dumps(ret)
    else:
        return json.dumps([])


@app.route(
    "/api/v1/agent/", methods=["GET"]
)  # should be post, but now we only need notification
def added_agent():
    handler.update_agents()
    return "Success", 200


@app.route("/api/v1/notify", methods=["POST"])
def notify():
    job = json.loads(request.get_data())
    logging.error(job)
    try:
        handler.job_enqueue(job)
        return "Success", 200
    except Exception:
        return "Internal Server Error", 500


@app.route("/api/v1/health")
def health_check():
    return "pong"


@app.route("/")
def root():
    return "Upstream Horizontal Scaler Controller"


if __name__ == "__main__":
    handler.start()
    app.run(host="0.0.0.0", debug=False, port=_port)
