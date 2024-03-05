#!/bin/bash

echo ====== TEST Start ======
SET=$(seq 0 20)
for i in $SET
    do          
        ./init_test_env_kind_with_simple_api_servers.sh default
        ./init_test_env_kind_with_simple_api_servers.sh fast
        ./init_test_env_kind_with_simple_api_servers.sh 
    done    

echo ====== TEST Done ======

# ./test_commands.sh