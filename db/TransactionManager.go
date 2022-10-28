package db

import (
	"database/sql"
	"errors"
	"github.com/Beerus-go/Beerus-DB/commons/dbutil"
	"github.com/Beerus-go/Beerus-DB/pool"
	"strconv"
)

type TxData struct {
	Tx   *sql.Tx
	Conn *pool.Connection
}

var txMaps = make(map[uint64]map[string]*TxData)

var snowflake *dbutil.SnowFlake

// Transaction Open a transaction
func Transaction() (uint64, error) {
	var err error
	var snowflakeId uint64

	if snowflake == nil {
		snowflake, err = dbutil.New(1)
		if err != nil {
			return 0, err
		}
	}

	snowflakeId, err = snowflake.Generate()

	connMaps := make(map[string]*TxData)
	for key, val := range dataSource {
		conn, err := val.GetConn()
		if err != nil {
			return 0, err
		}

		tx, err := conn.DB.Begin()
		if err != nil {
			return 0, err
		}

		txData := new(TxData)
		txData.Conn = conn
		txData.Tx = tx
		connMaps[key] = txData
	}

	txMaps[snowflakeId] = connMaps

	return snowflakeId, nil
}

// Commit transaction
func Commit(id uint64) error {
	connMaps := txMaps[id]
	if connMaps == nil || len(connMaps) <= 0 {
		return errors.New("")
	}

	for _, val := range connMaps {
		val.Tx.Commit()
		val.Conn.Close()
	}

	delete(txMaps, id)
	return nil
}

// Rollback transactions
func Rollback(id uint64) error {
	connMaps := txMaps[id]
	if connMaps == nil || len(connMaps) <= 0 {
		return errors.New("")
	}

	for _, val := range connMaps {
		val.Tx.Rollback()
		val.Conn.Close()
	}

	delete(txMaps, id)
	return nil
}

// GetTx Get the corresponding transaction manager based on ID and data source name
func GetTx(id uint64, dataSourceName string) (*sql.Tx, error) {
	txMap := txMaps[id]
	if txMap == nil {
		return nil, errors.New("no transaction operation with id " + strconv.FormatUint(id, 10) + " exists")
	}

	dataMap := txMap[dataSourceName]

	if dataMap == nil {
		return nil, errors.New("no data source with the name " + dataSourceName + " exists")
	}

	return dataMap.Tx, nil
}
