package db

import (
	"database/sql"
	"errors"
)

var txMaps = make(map[string]map[string]*sql.Tx)

func Transaction() (string, error) {
	connMaps := make(map[string]*sql.Tx)
	for key, val := range dataSource {
		conn, err := val.GetConn()
		if err != nil {
			return "", err
		}

		tx, err := conn.Begin()
		if err != nil {
			return "", err
		}

		connMaps[key] = tx
	}

	txMaps[""] = connMaps

	return "", nil
}

func Commit(id string) error {
	connMaps := txMaps[id]
	if connMaps == nil || len(connMaps) <= 0 {
		return errors.New("")
	}

	for _, val := range connMaps {
		val.Commit()
	}

	delete(txMaps, id)
	return nil
}

func Rollback(id string) error {
	connMaps := txMaps[id]
	if connMaps == nil || len(connMaps) <= 0 {
		return errors.New("")
	}

	for _, val := range connMaps {
		val.Rollback()
	}

	delete(txMaps, id)
	return nil
}

func GetTx(id string, dataSourceName string) *sql.Tx {
	return txMaps[id][dataSourceName]
}
