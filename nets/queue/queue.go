package queue

import (
	"container/list"
	"errors"
	"sync"
	"sync/atomic"
)

func New(capacity int) *Queue {
	queue := &Queue{capacity: capacity, list: list.New()}
	queue.cond = sync.NewCond(&queue.mutex)
	return queue
}

type Queue struct {
	list     *list.List
	closed   bool
	capacity int
	cond     *sync.Cond
	sph      int64
	spp      int64
	mutex    sync.Mutex
}

func (queue *Queue) Push(event interface{}) error {

	queue.mutex.Lock()
	err := queue.push(event)
	queue.mutex.Unlock()

	if err == nil {
		atomic.AddInt64(&queue.sph, 1)
		queue.cond.Signal()
	}
	return nil
}

func (queue *Queue) push(event interface{}) error {

	if queue.closed {
		return errClosed
	}

	list := queue.list
	capacity := queue.capacity
	if capacity > 0 && list.Len() >= capacity {
		return errOverflow
	}

	list.PushBack(event)
	return nil
}

func (queue *Queue) Pop() (interface{}, bool) {

	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	list := queue.list

	for {
		if queue.closed || list.Len() > 0 {
			break
		}
		queue.cond.Wait()
	}

	if list.Len() > 0 {
		return list.Remove(list.Front()), false
	}

	return nil, true
}

func (queue *Queue) Close() error {

	queue.mutex.Lock()
	err := queue.close()
	queue.mutex.Unlock()

	if err == nil {
		queue.cond.Signal()
	}
	return err
}

func (queue *Queue) close() error {

	if queue.closed {
		return errClosed
	}

	queue.closed = true
	return nil
}

var errClosed = errors.New("queue closed")
var errOverflow = errors.New("queue overflow")
