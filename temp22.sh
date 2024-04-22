#!/bin/bash
RANDOM_NUMBER=$((RANDOM % 10))
echo "Sleeping for $RANDOM_NUMBER seconds..."


_start=$(date +%s%3N)
echo $_start >> $_start.txt
echo 4
# Loop indefinitely
