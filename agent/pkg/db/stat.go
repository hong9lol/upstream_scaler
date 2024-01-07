package database

import (
	"errors"

	"github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/db/entity"
)

var _deployment map[string]entity.Deployment = map[string]entity.Deployment{}

// func (d *DB) GetStat(deploymentName string) (*entity.Deployment, error) {
// 	// var usages []entity.Usage
// 	// usages = append(usages, entity.Usage{
// 	// 	Usage: 80,
// 	// 	Timestamp: time.Now().Unix(),
// 	// })

// 	// var containers []entity.Container
// 	// containers = append(containers, entity.Container{
// 	// 	Name: "temp_container1",
// 	// 	CPURequest: 200,
// 	// 	Usages: usages,
// 	// },entity.Container{
// 	// 	Name: "temp_container2",
// 	// 	CPURequest: 100,
// 	// 	Usages: usages,
// 	// })

// 	// var pods []entity.Pod
// 	// pods = append(pods, entity.Pod{
// 	// 	Name: "temp_pod",
// 	// 	Containers: containers,
// 	// })

// 	// dummy := entity.Deployment {
// 	// 	Name: "temp_deployment",
// 	// 	Pods: pods,
// 	// }

// 	var stat *entity.Deployment
// 	stat = new(entity.Deployment)
// 	err := d.inst.View(func(tx *bolt.Tx) error {
// 		bucket := tx.Bucket([]byte("Deployment"))
// 		if bucket == nil {
// 			return fmt.Errorf("bucket not found")
// 		}

// 		value := bucket.Get([]byte(deploymentName))
// 		if value == nil {
// 			stat = nil
// 		} else {
// 			json.Unmarshal(value, &stat)
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}
// 	return stat, nil
// 	// return &dummy, nil
// }

// func (d *DB) UpdateStat(deployment *entity.Deployment) error {
// 	err := d.inst.Update(func(tx *bolt.Tx) error {
// 		bucket := tx.Bucket([]byte("Deployment"))
// 		if bucket == nil {
// 			return fmt.Errorf("bucket not found")
// 		}

// 		byteValue, _ := json.Marshal(deployment)
// 		err := bucket.Put([]byte(deployment.Name), byteValue)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (d *DB) GetStat(deploymentName string) (entity.Deployment, error) {
	if _, ok := _deployment[deploymentName]; !ok {
		// no data
		return _deployment[deploymentName], errors.New("NOT FOUND")
	}

	return _deployment[deploymentName], nil
}

func (d *DB) UpdateStat(deployment entity.Deployment) error {
	_deployment[deployment.Name] = deployment

	return nil
}
