#!/bin/bash

echo ====== TEST Start ======

# limits=( 20 30 40 50 60 )
limits=( 60 75 )
#limits=( 100 125 150 200 )
for limit in ${limits[@]}
    do    
        python3 limits.py $limit
        SET=$(seq 0 2)
        for i in $SET
            do  
	        echo $limit $i
	       	VAL1=60
		if [ ${limit} != ${VAL1} ] ; then
#			./init_test_env_gcp_with_deathstarbench_socialnetwork.sh default
			./init_test_env_gcp_with_deathstarbench_socialnetwork.sh fast
			./init_test_env_gcp_with_deathstarbench_socialnetwork.sh
#			echo 1
		else
			./init_test_env_gcp_with_deathstarbench_socialnetwork.sh
#			echo 2
		fi
#               ./init_test_env_gcp_with_deathstarbench_socialnetwork.sh default
#              ./init_test_env_gcp_with_deathstarbench_socialnetwork.sh fast
#              ./init_test_env_gcp_with_deathstarbench_socialnetwork.sh
           done
    done

echo ====== TEST Done ======i
