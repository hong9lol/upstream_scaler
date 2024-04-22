import os

# 디렉토리 경로 설정
directory_path = "./DeathStarBench/socialNetwork/benchmark_scripts/log/temp3"

# 해당 디렉토리 내의 모든 폴더를 찾습니다.
folders = [
    d
    for d in os.listdir(directory_path)
    if os.path.isdir(os.path.join(directory_path, d))
]


import yaml
import sys


# def get_limit_cpu():
#     # The path to your YAML file
#     yaml_file_path = "DeathStarBench/socialNetwork/helm-chart/socialnetwork/values.yaml"

#     key_to_change = "limitCpu"

#     ret = ""
#     with open(yaml_file_path, "r") as file:
#         yaml_content = yaml.safe_load(file)

#     ret = yaml_content["global"][key_to_change]
#     return ret.split("m")[0]


# 각 폴더를 순회하며 'out.log' 파일을 찾아 해당 파일 내용을 읽습니다.
i = 1
_type = ""
sorted_folders = sorted(folders)
result = dict()
test_cnt = 0
init_limit_cpu = 30
for folder in sorted_folders:
    print("[Test #" + str(test_cnt) + "]")
    test_cnt += 1
    # print("folder name:", folder)
    if i % 3 == 1:
        _type = "default"
    elif i % 3 == 2:
        _type = "fast"
    else:
        _type = "upstream"
    i += 1
    total_sent = 0
    success_sent = 0
    result[folder] = [test_cnt]
    file_path = os.path.join(directory_path, folder, "output.log")
    if os.path.exists(file_path):
        # with open(file_path, "r") as file:
        #     for line in file:
        #         if " requests " in line:
        #             if int(line.split(" requests ")[0]) < 800:
        #                 total_sent = 1
        #                 success_sent = 1

        with open(file_path, "r") as file:
            result[folder].append(_type)
            # read output.log and get total requests and success requests
            k = 0
            t_sent = 0
            s_sent = 0
            for line in file:
                if "Sent " in line:
                    total_sent += int(line.split(" ")[1])
                    t_sent = int(line.split(" ")[1])
                elif " requests " in line:
                    success_sent += int(line.split(" requests ")[0])
                    s_sent = int(line.split(" requests ")[0])
                    print(t_sent, s_sent)
                    # if int(line.split(" requests ")[0]) < 700:
                    #     total_sent = 1
                    #     success_sent = 1
                    #     break
                    k += 1
                    # if k > 18:
                    #     break
        print(total_sent)
        print(
            _type, total_sent, success_sent, round(success_sent / total_sent * 100, 2)
        )
        result[folder].append(
            [total_sent, success_sent, round(success_sent / total_sent * 100, 2)]
        )

    d = dict()
    file_path = os.path.join(directory_path, folder, "podcnt.txt")
    if os.path.exists(file_path):
        with open(file_path, "r") as file:
            # read output.log and get total requests and success requests
            k = 0
            for line in file:
                key = line.split("	")[0]
                # print(line, key)
                if "redis" in key:
                    continue
                if "mongo" in key:
                    continue
                if "memcached" in key:
                    continue
                if "nginx" in key:
                    continue
                if "media-front" in key:
                    continue
                if "jaeger" in key:
                    continue
                if key in d:
                    d[key][0] += 1
                    d[key].append(line.split("	")[2])
                else:
                    d[key] = [1, line.split("	")[2]]
        _l = []
        for key, item in d.items():
            # print(key, item)
            _l.append([key, item[0]])
        result[folder].append(_l)

    # print("=======================")

import pandas as pd
import os

# print(result)
writer = pd.ExcelWriter("data.xlsx")
i = 0
idx = 0
repeat = 1 * 3
limit_add = 0
for key, item in result.items():
    print("Test Case #", item[0])
    # print("Type:", item[1])
    print("request", item[2])

    deployments = sorted(item[3])
    deployments.insert(0, ["type", item[1]])
    deployments.insert(0, ["total_request", item[2][0]])
    deployments.insert(0, ["200_request", item[2][1]])
    deployments.insert(0, ["limitCPU", init_limit_cpu + limit_add])

    for deployment in deployments:
        print(deployment[0], deployment[1])

    df1 = pd.DataFrame(deployments, columns=["deployment", "Pods"])
    i += 1
    if item[1] == "default":
        df1.to_excel(
            writer,
            sheet_name="Data",
            # index_label="index",
            index=False,
            header=None,
            startrow=idx,
            startcol=0,
        )

    elif item[1] == "fast":
        df1.to_excel(
            writer,
            sheet_name="Data",
            # index_label="index",
            index=False,
            header=None,
            startrow=idx,
            startcol=3,
        )

    else:
        df1.to_excel(
            writer,
            sheet_name="Data",
            # index_label="index",
            index=False,
            header=None,
            startrow=idx,
            startcol=6,
        )
        idx += len(deployments) + 1
        if i % repeat == 0:
            limit_add += 15

    print("=======================")
writer.close()
