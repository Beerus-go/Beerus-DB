package db

import (
	"github.com/yuyenews/Beerus-DB/pool"
)

var dataSource = make(map[string]*pool.DbPool)

// AddDataSource Add DbPool as a data source
func AddDataSource(name string, dbPool *pool.DbPool) {
	dataSource[name] = dbPool
}

// GetConnection Obtaining a connection from a specified data source
func GetConnection(name string) (*pool.Connection, error) {
	return dataSource[name].GetConn()
}
