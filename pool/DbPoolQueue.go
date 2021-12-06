package pool

import (
	"database/sql"
	"sync"
)

type DbPoolQueueItem struct {
	Conn *sql.DB
	Next *DbPoolQueueItem
}

type DbPoolQueue struct {
	First *DbPoolQueueItem
	Last  *DbPoolQueueItem
	Size  int64
	mutex sync.Mutex
}

func (dp *DbPoolQueue) Add(item *sql.DB) {
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

func (dp *DbPoolQueue) Pop() *sql.DB {
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
