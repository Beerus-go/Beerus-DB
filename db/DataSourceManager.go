package db

import (
	"database/sql"
	"github.com/yuyenews/Beerus-DB/pool"
)

var dataSource = make(map[string]*pool.DbPool)

func AddDataSource(name string, dbPool *pool.DbPool) {
	dataSource[name] = dbPool
}

func GetDataSource(name string) (*sql.DB, error) {
	return dataSource[name].GetConn()
}
