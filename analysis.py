import os

# 디렉토리 경로 설정
directory_path = "./DeathStarBench/socialNetwork/benchmark_scripts/log/tt"

# 해당 디렉토리 내의 모든 폴더를 찾습니다.
folders = [
    d
    for d in os.listdir(directory_path)
    if os.path.isdir(os.path.join(directory_path, d))
]

# 각 폴더를 순회하며 'out.log' 파일을 찾아 해당 파일 내용을 읽습니다.
i = 1
_type = ""
sorted_folders = sorted(folders)
result = dict()
for folder in sorted_folders:
    # print("folder name:", folder)
    if i % 2 == 1:
        _type = "default"
    else:
        _type = "upstream"
    result[folder] = [_type]
    i += 1
    total_sent = 0
    success_sent = 0
    file_path = os.path.join(directory_path, folder, "output.log")
    if os.path.exists(file_path):
        with open(file_path, "r") as file:
            # read output.log and get total requests and success requests
            k = 0
            for line in file:
                if "Sent " in line:
                    total_sent += int(line.split(" ")[1])
                elif " requests " in line:
                    success_sent += int(line.split(" requests ")[0])
                    k += 1
                    if k > 18:
                        break

        # print(
        #     _type, total_sent, success_sent, round(success_sent / total_sent * 100, 2)
        # )
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
idx = 0
for key, item in result.items():
    print("Type:", item[0])
    print("request", item[1])
    deployments = sorted(item[2])
    for deployment in deployments:
        print(deployment[0], deployment[1])

    df1 = pd.DataFrame(deployments, columns=["deployment", "Pods"])

    if item[0] == "default":
        df1.to_excel(
            writer,
            sheet_name="Data",
            # index_label="index",
            index=False,
            header=None,
            startrow=idx,
            startcol=0,
        )

    else:
        df1.to_excel(
            writer,
            sheet_name="Data",
            # index_label="index",
            index=False,
            header=None,
            startrow=idx,
            startcol=3,
        )
        idx += len(deployments) + 1

    print("=======================")
writer.close()
