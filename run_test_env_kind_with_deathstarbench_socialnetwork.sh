#!/bin/sh

echo 1. Run Log
cd ./DeathStarBench/socialNetwork/benchmark_scripts
baseLogPath=./log
currentTime=`date +"%m-%d_%H%M%S"`
mkdir $baseLogPath/$currentTime
logPath=$baseLogPath/$currentTime
./log.sh $logPath & log=$!

echo 2. Run Benchmark
./run_social_benchmark.sh $logPath
cd -

echo 3. Kill Log Proc
kill -9 $log

# 파드 리스트 가져오기
./time_checker.sh

echo ====== Done ======