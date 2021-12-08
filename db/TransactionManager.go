package db

import (
	"database/sql"
	"errors"
	"github.com/yuyenews/Beerus-DB/pool"
)

type TxData struct {
	Tx   *sql.Tx
	Conn *pool.Connection
}

var txMaps = make(map[string]map[string]*TxData)

// Transaction Open a transaction
func Transaction() (string, error) {
	connMaps := make(map[string]*TxData)
	for key, val := range dataSource {
		conn, err := val.GetConn()
		if err != nil {
			return "", err
		}

		tx, err := conn.DB.Begin()
		if err != nil {
			return "", err
		}

		txData := new(TxData)
		txData.Conn = conn
		txData.Tx = tx
		connMaps[key] = txData
	}

	txMaps[""] = connMaps

	return "", nil
}

// Commit transaction
func Commit(id string) error {
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
func Rollback(id string) error {
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
func GetTx(id string, dataSourceName string) (*sql.Tx, error) {
	txMap := txMaps[id]
	if txMap == nil {
		return nil, errors.New("no transaction operation with id " + id + " exists")
	}

	dataMap := txMap[dataSourceName]

	if dataMap == nil {
		return nil, errors.New("no data source with the name " + dataSourceName + " exists")
	}

	return dataMap.Tx, nil
}
