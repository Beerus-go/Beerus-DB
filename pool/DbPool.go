package pool

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

// DbPool Connection Pool
type DbPool struct {
	Url          string
	InitialSize  int
	ExpandSize   int
	MaxOpen      int
	numOpen      int
	MinOpen      int
	isInitialize bool
	PoolQueue    DbPoolQueue
	mutex        sync.Mutex
}

// GetConn Get a connection
func (pool *DbPool) GetConn() (*Connection, error) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	// Verify connection pool configuration parameters and create initialized connections
	err := pool.validAndInit()
	if err != nil {
		return nil, err
	}

	// Get a connection from the queue
	conn := pool.PoolQueue.Poll()

	// If there are no more connections in the queue, expand the connections and fetch them from the queue once more
	if conn == nil {
		err := pool.expansion()
		if err != nil {
			return nil, err
		}
		conn = pool.PoolQueue.Poll()
	}
	if conn != nil {
		conn.Status = true
	}

	// If idle connections are already <= minimum connections, expand connections
	if pool.PoolQueue.Size <= pool.MinOpen {
		pool.expansion()
	}

	return conn, nil
}

// Close Returning used connections to the connection pool
func (pool *DbPool) Close(conn *Connection) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	conn.Status = false
	pool.numOpen--
	pool.PoolQueue.Add(conn)
}

// expansion Expanded Connectivity
func (pool *DbPool) expansion() error {
	if pool.numOpen < pool.MaxOpen {
		remaining := pool.MaxOpen - pool.numOpen
		if remaining >= pool.ExpandSize {
			return pool.createConn(pool.ExpandSize)
		}

		return pool.createConn(remaining)
	}

	return errors.New("the current number of active connections has exceeded MaxOpen")
}

// createConn Create the specified number of connections to the queue
func (pool *DbPool) createConn(size int) error {
	if size <= 0 {
		return errors.New("the number of expansions must be greater than 0")
	}
	for i := 0; i < size; i++ {
		db, err := sql.Open("mysql", pool.Url)
		if err != nil {
			return err
		}

		conn := new(Connection)
		conn.DB = db
		conn.Pool = pool
		conn.Status = true

		pool.PoolQueue.Add(conn)
	}

	pool.numOpen += size
	return nil
}

// validAndInit Verify connection pool configuration parameters and create initialized connections
func (pool *DbPool) validAndInit() error {
	if pool.isInitialize == false {
		if pool.Url == "" {
			return errors.New("url must not be empty")
		}

		if pool.ExpandSize <= 0 {
			return errors.New("ExpandSize must > 0")
		}

		if pool.InitialSize <= 0 {
			return errors.New("InitialSize must > 0")
		}

		if pool.InitialSize < pool.MinOpen {
			return errors.New("InitialSize cannot be smaller than MinOpen")
		}

		if pool.MaxOpen <= 0 {
			return errors.New("MaxOpen must > 0")
		}

		if pool.MinOpen < 0 {
			return errors.New("MaxOpen must >= 0")
		}

		if pool.MaxOpen < pool.MinOpen {
			return errors.New("MinOpen must not be larger than MaxOpen")
		}

		pool.createConn(pool.InitialSize)

		pool.isInitialize = true
	}
	return nil
}
