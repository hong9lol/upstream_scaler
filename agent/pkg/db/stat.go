package database

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/db/entity"
)

func (d *DB) GetStat(deploymentName string) (*entity.Deployment, error) {
	// var usages []entity.Usage
	// usages = append(usages, entity.Usage{
	// 	Usage: 80,
	// 	Timestamp: time.Now().Unix(),
	// })

	// var containers []entity.Container
	// containers = append(containers, entity.Container{
	// 	Name: "temp_container1",
	// 	CPURequest: 200,
	// 	Usages: usages,
	// },entity.Container{
	// 	Name: "temp_container2",
	// 	CPURequest: 100,
	// 	Usages: usages,
	// })

	// var pods []entity.Pod
	// pods = append(pods, entity.Pod{
	// 	Name: "temp_pod",
	// 	Containers: containers,
	// })

	// dummy := entity.Deployment {
	// 	Name: "temp_deployment",
	// 	Pods: pods,
	// }

	var stat *entity.Deployment
	stat = new(entity.Deployment)
	err := d.inst.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Deployment"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		value := bucket.Get([]byte(deploymentName))
		if value == nil {
			stat = nil
		} else {
			json.Unmarshal(value, &stat)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return stat, nil
	// return &dummy, nil
}

func (d *DB) UpdateStat(deployment *entity.Deployment) error {
	// deployment check
	// 키 존재 여부 확인
	err := d.inst.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Deployment"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		byteValue, _ := json.Marshal(deployment)
		err := bucket.Put([]byte(deployment.Name), byteValue)
		if err != nil {
			return err
		}
		// value := bucket.Get([]byte(deployment.Name))
		// if value == nil {
		// 	// put new deployment
		// 	fmt.Printf("Key '%s' does not exist\n", deployment.Name)
		// 	byteValue, _ := json.Marshal(deployment)
		// 	bucket.Put([]byte(deployment.Name), byteValue)
		// } else {
		// 	// update the deployment
		// 	// prev, _ := d.GetStats(deployment.Name)

		// 	fmt.Printf("Key '%s' exists with value: %s\n", deployment.Name, value)

		// }

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
