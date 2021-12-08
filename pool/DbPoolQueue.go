package pool

import (
	"sync"
)

// DbPoolQueueItem Elements inside the queue
type DbPoolQueueItem struct {
	Conn *Connection
	Next *DbPoolQueueItem
}

// DbPoolQueue Queue
type DbPoolQueue struct {
	First *DbPoolQueueItem
	Last  *DbPoolQueueItem
	Size  int
	mutex sync.Mutex
}

// Add connections to the queue
func (dp *DbPoolQueue) Add(item *Connection) {
	dp.mutex.Lock()

	defer dp.mutex.Unlock()

	dpItem := new(DbPoolQueueItem)
	dpItem.Conn = item
	dpItem.Next = nil

	if dp.First == nil {
		dp.First = dpItem
	}

	if dp.Last != nil {
		dp.Last.Next = dpItem
	}
	dp.Last = dpItem
	dp.Size++
}

// Poll Taking elements out of the queue
func (dp *DbPoolQueue) Poll() *Connection {
	dp.mutex.Lock()

	defer dp.mutex.Unlock()

	first := dp.First
	if first == nil {
		dp.Last = nil
		return nil
	}
	dp.First = dp.First.Next
	dp.Size--
	return first.Conn
}
