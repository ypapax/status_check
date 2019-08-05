package queue

import (
	"sync"
)

type Queue struct {
	mtx           sync.Mutex
	lastUsedIndex int
	arr           []int
}

func New(arr []int) Queue {
	return Queue{
		arr: arr,
	}
}

func (q *Queue) Next() int {
	if len(q.arr) == 1 {
		return q.arr[0]
	}
	q.mtx.Lock()
	defer func() {
		q.lastUsedIndex++
		if q.lastUsedIndex >= len(q.arr) {
			q.lastUsedIndex = 0
		}
		q.mtx.Unlock()
	}()
	v := q.arr[q.lastUsedIndex]
	return v
}
