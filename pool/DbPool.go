package pool

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

type DbPool struct {
	Url          string
	InitialSize  int
	ExpSize      int
	MaxActive    int
	nowActive    int
	isInitialize bool
	PoolQueue    DbPoolQueue
	mutex        sync.Mutex
}

func (pool *DbPool) GetConn() (*sql.DB, error) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	err := pool.valid()
	if err != nil {
		return nil, err
	}

	conn := pool.PoolQueue.Pop()

	if conn == nil {
		err := pool.expansion()
		if err != nil {
			return nil, err

		}
		conn = pool.PoolQueue.Pop()
	}
	return conn, nil
}

func (pool *DbPool) Close(conn *sql.DB) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	conn.Close()
	pool.nowActive--
}

func (pool *DbPool) expansion() error {
	if pool.nowActive < pool.MaxActive {
		remaining := pool.MaxActive - pool.nowActive
		if remaining >= pool.ExpSize {
			return pool.createConn(pool.ExpSize)
		}

		return pool.createConn(remaining)
	}

	return errors.New("the current number of active connections has exceeded MaxActive")
}

func (pool *DbPool) createConn(size int) error {
	for i := 0; i < size; i++ {
		db, err := sql.Open("mysql", pool.Url)
		if err != nil {
			return err
		}
		pool.PoolQueue.Add(db)
	}

	pool.nowActive += size
	return nil
}

func (pool *DbPool) valid() error {
	if pool.isInitialize == false {
		if pool.Url == "" {
			return errors.New("url must not be empty")
		}

		if pool.ExpSize <= 0 {
			return errors.New("ExpSize must not be empty")
		}

		if pool.InitialSize <= 0 {
			return errors.New("InitialSize must not be empty")
		}

		if pool.MaxActive <= 0 {
			return errors.New("MaxActive must not be empty")
		}

		pool.createConn(pool.InitialSize)

		pool.isInitialize = true
	}
	return nil
}
