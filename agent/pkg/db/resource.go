package database

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/common"
)

type DB struct {
	inst *bolt.DB
}

var _db *DB
var lock = &sync.Mutex{}

func newDB() *DB {
	boltDB, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = boltDB.Update(func(tx *bolt.Tx) error {
		if err := tx.DeleteBucket([]byte("Deployment")); err != nil {
			return err
		}
		if err := tx.DeleteBucket([]byte("HPA")); err != nil {
			return err
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

func (d *DB) UpdateHPA(data []byte) error {
	hpas := []common.HPAObject{}
	json.Unmarshal(data, &hpas)

	// remove old one
	err := _db.inst.Batch(func(tx *bolt.Tx) error {
		vhpa := common.HPAObject{}
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

func (d *DB) GetAllHPA() []common.HPAObject {
	hpas := []common.HPAObject{}
	hpa := common.HPAObject{}

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

func (d *DB) UpdateStats(deploymentName string) error {

	// _usage := []Usage{
	// 	{Usage: 1, Timestamp: 2},
	// 	{Usage: 3, Timestamp: 4},
	// 	{Usage: 5, Timestamp: 6},
	// }

	// _pod := []Pod{
	// 	{Name: "pod1", Usages: _usage},
	// 	{Name: "pod2", Usages: _usage},
	// }

	// entryBytes, _ := json.Marshal(_pod)

	// err := _db.inst.Update(func(tx *bolt.Tx) error {
	// 	err := tx.Bucket([]byte("Deployment")).Put([]byte("deployment1"), []byte(entryBytes))
	// 	if err != nil {
	// 		return fmt.Errorf("could not insert weight: %v", err)
	// 	}
	// 	return nil
	// })
	// fmt.Println("Added deployment1")

	// var pod []Pod
	// err = _db.inst.View(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte("Deployment"))
	// 	json.Unmarshal(b.Get([]byte("deployment1")), &pod)
	// 	fmt.Println(pod[0].Name)
	// 	// fmt.Println(string(b.Get([]byte("deployment1"))))
	// 	b.ForEach(func(k, v []byte) error {
	// 		fmt.Println(string(k), string(v))
	// 		return nil
	// 	})
	// 	return nil
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return nil

	// fmt.Println("Added Weight")

	// txn := db.inst.Txn(true)

	// _usage := []Usage{
	// 	{usage: 1, timestamp: 2},
	// 	{usage: 3, timestamp: 4},
	// 	{usage: 5, timestamp: 6},
	// }

	// _pod := []Pod{
	// 	{name: "pod1", usages: _usage},
	// 	{name: "pod2", usages: _usage},
	// }

	// deployment := []Deployment{
	// 	{name: "deployment1", pods: _pod},
	// 	{name: "deployment2", pods: _pod},
	// }

	// for _, p := range deployment {
	// 	if err := txn.Insert("deployment", p); err != nil {
	// 		panic(err)
	// 	}
	// }

	// // Commit the transaction
	// txn.Commit()

	// return db
}

// func NewDB() DB {
// 	schema := &memdb.DBSchema{
// 		Tables: map[string]*memdb.TableSchema{
// 			"deployment": {
// 				Name: "deployment",
// 				Indexes: map[string]*memdb.IndexSchema{"id": {
// 					Name:   "pod",
// 					Unique: true,
// 					Indexer: &memdb.CompoundMultiIndex{
// 						Indexes: []memdb.Indexer{
// 							&memdb.StringFieldIndex{
// 								Field: "name",
// 							},
// 							&memdb.UintFieldIndex{
// 								Field: "timestamp",
// 							},
// 						},
// 					},
// 				},
// 				},
// 			},
// 		},
// 	}
// 	memdb, err := memdb.NewMemDB(schema)
// 	if err != nil {
// 		panic(err)
// 	}

// 	db := DB{}
// 	db.inst = memdb

// 	txn := db.inst.Txn(true)

// 	_usage := []Usage{
// 		{usage: 1, timestamp: 2},
// 		{usage: 3, timestamp: 4},
// 		{usage: 5, timestamp: 6},
// 	}

// 	_pod := []Pod{
// 		{name: "pod1", usages: _usage},
// 		{name: "pod2", usages: _usage},
// 	}

// 	deployment := []Deployment{
// 		{name: "deployment1", pods: _pod},
// 		{name: "deployment2", pods: _pod},
// 	}

// 	for _, p := range deployment {
// 		if err := txn.Insert("deployment", p); err != nil {
// 			panic(err)
// 		}
// 	}

// 	// Commit the transaction
// 	txn.Commit()

// 	return db
// }

// func (n *Npc) getSerialNumber() int64 {
// 	return n.NpcNumber
// }

// func (n *Npc) getName() string {
// 	return n.Name
// }

// Create a sample struct
// type Person struct {
// 	Email string
// 	Name  string
// 	Age   int
// }

// Create the DB schema

// Create a new data base
// db, err := memdb.NewMemDB(schema)
// if err != nil {
// 	panic(err)
// }

// // Create a write transaction
// txn := db.Txn(true)

// // Insert some people
// people := []*Person{
// 	&Person{"joe@aol.com", "Joe", 30},
// 	&Person{"lucy@aol.com", "Lucy", 35},
// 	&Person{"tariq@aol.com", "Tariq", 21},
// 	&Person{"dorothy@aol.com", "Dorothy", 53},
// }
// for _, p := range people {
// 	if err := txn.Insert("person", p); err != nil {
// 		panic(err)
// 	}
// }

// // Commit the transaction
// txn.Commit()

// // Create read-only transaction
// txn = db.Txn(false)
// defer txn.Abort()

// // Lookup by email
// raw, err := txn.First("person", "id", "joe@aol.com")
// if err != nil {
// 	panic(err)
// }

// // Say hi!
// fmt.Printf("Hello %s!\n", raw.(*Person).Name)

// // List all the people
// it, err := txn.Get("person", "id")
// if err != nil {
// 	panic(err)
// }

// fmt.Println("All the people:")
// for obj := it.Next(); obj != nil; obj = it.Next() {
// 	p := obj.(*Person)
// 	fmt.Printf("  %s\n", p.Name)
// }

// // Range scan over people with ages between 25 and 35 inclusive
// it, err = txn.LowerBound("person", "age", 25)
// if err != nil {
// 	panic(err)
// }

// fmt.Println("People aged 25 - 35:")
// for obj := it.Next(); obj != nil; obj = it.Next() {
// 	p := obj.(*Person)
// 	if p.Age > 35 {
// 		break
// 	}
// 	fmt.Printf("  %s is aged %d\n", p.Name, p.Age)
// }

// queue
// import "container/list"

// type Queue struct {
//   v *list.List
// }

// func NewQueue() *Queue {
//   return &Queue{list.New()}
// }

// func (q *Queue) Push(v interface{}) {
//   q.v.PushBack(v)
// }

// func (q *Queue) Pop() interface{} {
//   front := q.v.Front()
//   if front == nil {
//     return nil
//   }

//   return q.v.Remove(front)
// }
