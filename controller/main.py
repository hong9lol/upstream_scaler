import json

from flask import Flask, request
from kube_client import api
from manager import handler
from db import hpa

app = Flask(__name__)


@app.route('/api/v1/node', methods=['GET'])
def node():
    ret = api.get_node_list()
    print(ret)
    return json.dumps(ret)


@app.route('/api/v1/deployment', methods=['GET'])
def get_deployment():
    ret = api.get_deployment_list()
    print(ret)
    return json.dumps(ret)


@app.route('/api/v1/hpa', methods=['GET'])
def get_hpa():
    ret = hpa.get_all_hpas()
    if ret.count != 0:
        return json.dumps(ret)
    else:
        return json.dumps([])


@app.route("/api/v1/notify", methods=['POST'])
def notify():
    job = json.loads(request.get_data())
    try:
        handler.job_enqueue(job)
        return "Success", 200
    except Exception:
        return "Internal Server Error", 500


@app.route('/api/v1/health')
def health_check():
    return "pong"


@app.route('/')
def root():
    return "Upstream Horizontal Scaler"


if __name__ == '__main__':
    handler.start()
    app.run(host="0.0.0.0", debug=False, port=5001)
