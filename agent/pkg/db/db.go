package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/boltdb/bolt"
)

type DB struct {
	inst *bolt.DB
}

var _db *DB
var lock = &sync.Mutex{}

func newDB() *DB {
	boltDB, err := bolt.Open("upstream_scaler.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = boltDB.Update(func(tx *bolt.Tx) error {
		if err := tx.DeleteBucket([]byte("Deployment")); err != nil {
			// return err
		}
		if err := tx.DeleteBucket([]byte("HPA")); err != nil {
			// return err
		}
		_, err := tx.CreateBucketIfNotExists([]byte("Deployment"))
		if err != nil {
			return fmt.Errorf("could not create deployment bucket: %v", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte("HPA"))
		if err != nil {
			return fmt.Errorf("could not create hpa bucket: %v", err)
		}
		return nil
	})

	if err != nil {
		fmt.Println("DB Setup Failed")
		panic(err)
	}

	fmt.Println("DB Setup Done")
	_db = new(DB)
	_db.inst = boltDB
	return _db
}

func GetInstance() *DB {
	if _db == nil {
		lock.Lock()
		defer lock.Unlock()
		if _db == nil {
			return newDB()
		} else {
			fmt.Println("_db instance already created.")
		}
	}
	return _db
}
