package prioqueue

import (
	"container/list"
	"sync"
)

type OptionPrioQueue func(*PrioQueue) error

type prioQueue struct {
	*list.List
	prio  int
	mutex sync.Mutex
}

type PrioQueue struct {
	queues *list.List
	mutex  sync.RWMutex
}

// default queue priority 1, the higher the value, the higher the priority.
func NewPrioQueue(opts ...OptionPrioQueue) (*PrioQueue, error) {
	pq := &PrioQueue{
		queues: list.New(),
	}
	err := error(nil)

	for _, opt := range opts {
		err = opt(pq)
		if err != nil {
			return nil, err
		}
	}
	return pq, nil
}

func (pq *PrioQueue) PrioPush(prio int, data interface{}) {
	queue := (*prioQueue)(nil)
	pq.mutex.RLock()
	for elem := pq.queues.Front(); elem != nil; elem = elem.Next() {
		value, _ := elem.Value.(*prioQueue)
		if value.prio == prio {
			queue = value
			break
		}
	}
	pq.mutex.RUnlock()

	if queue == nil {
		pq.mutex.Lock()
		if queue == nil {
			queue = &prioQueue{
				prio: prio,
				List: list.New(),
			}
			pq.queues.PushFront(queue)
		}
		pq.mutex.Unlock()
	}

	queue.mutex.Lock()
	queue.PushFront(data)
	queue.mutex.Unlock()
}

func (pq *PrioQueue) Push(data interface{}) {
	queue := (*prioQueue)(nil)
	pq.mutex.RLock()
	for elem := pq.queues.Front(); elem != nil; elem = elem.Next() {
		value, _ := elem.Value.(*prioQueue)
		if value.prio == 1 {
			queue = value
			break
		}
	}
	pq.mutex.RUnlock()

	if queue == nil {
		// critical
		return
	}

	queue.mutex.Lock()
	queue.PushFront(data)
	queue.mutex.Unlock()
}

func (pq *PrioQueue) Pop() interface{} {
	queue := (*prioQueue)(nil)
	pq.mutex.RLock()
	for elem := pq.queues.Back(); elem != nil; elem = elem.Prev() {
		queue, _ = elem.Value.(*prioQueue)
		queue.mutex.Lock()
		if queue.Len() != 0 {
			data := queue.Back()
			queue.mutex.Unlock()
			return queue.Remove(data)
		}
		queue.mutex.Unlock()
	}
	pq.mutex.RUnlock()
	return nil
}
