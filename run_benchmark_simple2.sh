#!/bin/bash

# Target URL
targetIP=`cat target_url.txt`
url="http://$targetIP:8080/api/v1/mid"
# Set the deployment name and namespace
DEPLOYMENT_NAME="simple-server"
NAMESPACE="default"

_start=$(date +%s%3N)
echo $_start >> $_start.txt
# Loop indefinitely
while true; do
  # Start time measurement
  start=$(date +%s.%N)
  
  # Make 25 requests
  for i in {1..80}; do
    curl -s "$url" > /dev/null
  done

  # Wait for the next second to start a new batch of 20 requests
  # end=$(date +%s.%N)
  # elapsed=$(echo "$end - $start" | bc)
  # sleep=$(echo "0.2 - $elapsed" | bc)
  
  # # Only sleep if elapsed time is less than 1 second
  # if (( $(echo "$sleep > 0" | bc -l) )); then
  #   sleep $sleep
  # fi

CURRENT_REPLICAS=$(kubectl get deployment "$DEPLOYMENT_NAME" -n "$NAMESPACE" -o=jsonpath='{.status.replicas}')
if [ "$CURRENT_REPLICAS" -gt 2 ]; then
  break
fi

done
_end=$(date +%s%3N)
echo $_end >> $_start.txt
duration=$((_end - _start))
echo $duration >> $_start.txt