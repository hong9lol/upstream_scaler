#!/bin/bash

echo ====== TEST Start ======

# limits=( 20 30 40 50 60 )
limits=( 30 45 60 75 )
for limit in ${limits[@]}
    do    
        python3 limits.py $limit
        SET=$(seq 0 5)
        for i in $SET
            do  
                echo $limit $i
                #./init_test_env_gcp_with_deathstarbench_socialnetwork.sh default
                #./init_test_env_gcp_with_deathstarbench_socialnetwork.sh fast
                ./init_test_env_gcp_with_deathstarbench_socialnetwork.sh
            done
    done

echo ====== TEST Done ======
