package database

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/db/entity"
)

func (d *DB) UpdateHPA(data []byte) error {
	hpas := map[string]entity.HPA{}
	e := json.Unmarshal(data, &hpas)
	if e != nil {
		fmt.Println(e.Error())
	}
	// remove old one
	err := _db.inst.Batch(func(tx *bolt.Tx) error {
		vhpa := entity.HPA{}
		b := tx.Bucket([]byte("HPA"))       // get all HPA bucket
		b.ForEach(func(k, v []byte) error { // find if it is included the new hpa list
			json.Unmarshal(v, &vhpa)
			flag := false
			for _, hpa := range hpas {
				if hpa.Name == vhpa.Name {
					flag = true
				}
			}
			if flag == false {
				if err := b.Delete(k); err != nil {
					return err
				}
			}
			return nil
		})
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// add new one
	err = _db.inst.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("HPA")) // get all HPA bucket
		// fmt.Println(hpas)
		for _, hpa := range hpas {
			if b.Get([]byte(hpa.Name)) == nil {
				byteValue, _ := json.Marshal(hpa)
				if err := b.Put([]byte(hpa.Name), byteValue); err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error occurred while saving shows: %s", err.Error())
	}

	return nil
}

func (d *DB) GetAllHPA() []entity.HPA {
	hpas := []entity.HPA{}
	hpa := entity.HPA{}

	err := _db.inst.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("HPA"))       // get all HPA bucket
		b.ForEach(func(k, v []byte) error { // find if it is included the new hpa list
			// fmt.Println(string(v))
			json.Unmarshal(v, &hpa)
			hpas = append(hpas, hpa)
			return nil
		})
		return nil
	})

	if err != nil {
		fmt.Printf("Error occurred while saving shows: %s", err.Error())
	}

	return hpas
}
