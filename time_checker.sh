#!/bin/bash

suffix="upstream"
if [ "$1" = "default" ]; then
    suffix="default"
fi

currentTime=`date +"%m-%d_%H%M%S"`
# 파드 리스트 가져오기
pod_list=$(kubectl get pods -n default -o=jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.metadata.creationTimestamp}{"\t"}{.metadata.labels.app}{"\n"}{end}' | sort -k2)

# 가져온 파드 리스트 출력
echo -e "Deployment 이름\t파드 이름\t생성 시간"
while IFS=$'\t' read -r pod_name creation_timestamp deployment_name; do
    echo -e "$deployment_name\t$pod_name\t$creation_timestamp"
# done <<< "$pod_list" > ${currentTime}_${suffix}.txt
done <<< "$pod_list" > podcnt.txt