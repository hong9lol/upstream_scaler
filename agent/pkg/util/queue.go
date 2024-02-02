package util

import (
	"container/list"
	"sync"
)

var (
	Mutex = &sync.Mutex{}
)

type Queue struct {
	v *list.List
}

func NewQueue() *Queue {
	return &Queue{list.New()}
}

func (q *Queue) Enqueue(v interface{}) {
	Mutex.Lock()

	q.v.PushBack(v)
	// remove oldest data
	if q.v.Len() > 10 {
		front := q.v.Front()
		q.v.Remove(front)
	}

	Mutex.Unlock()
}
