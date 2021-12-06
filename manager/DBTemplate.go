package manager

import (
	"database/sql"
	"github.com/yuyenews/Beerus-DB/db"
)

// -----------------------------------------------

func SelectList(dataSourceName string, sql string, params []interface{}) ([]map[string]string, error) {
	conn, err := db.GetDataSource(dataSourceName)
	if err != nil {
		return nil, err
	}
	return db.SelectList(conn, sql, params)
}

func SelectListByMap(dataSourceName string, sql string, paramsMap map[string]string) ([]map[string]string, error) {
	sqlStr, params := SqlConvert(sql, paramsMap)
	return SelectList(dataSourceName, sqlStr, params)
}

// -----------------------------------------------

func SelectListByTx(txId string, dataSourceName string, sql string, params []interface{}) ([]map[string]string, error) {
	tx := db.GetTx(txId, dataSourceName)
	return db.SelectListByTx(tx, sql, params)
}

func SelectListByTxMap(txId string, dataSourceName string, sql string, paramsMap map[string]string) ([]map[string]string, error) {
	sqlStr, params := SqlConvert(sql, paramsMap)
	return SelectListByTx(txId, dataSourceName, sqlStr, params)
}

// -----------------------------------------------

func Update(dataSourceName string, sql string, params []interface{}) (sql.Result, error) {
	conn, err := db.GetDataSource(dataSourceName)
	if err != nil {
		return nil, err
	}
	return db.Update(conn, sql, params)
}

func UpdateByMap(dataSourceName string, sql string, paramsMap map[string]string) (sql.Result, error) {
	sqlStr, params := SqlConvert(sql, paramsMap)
	return Update(dataSourceName, sqlStr, params)
}

// -----------------------------------------------

func UpdateByTx(txId string, dataSourceName string, sql string, params []interface{}) (sql.Result, error) {
	tx := db.GetTx(txId, dataSourceName)
	return db.UpdateByTx(tx, sql, params)
}

func UpdateByTxMap(txId string, dataSourceName string, sql string, paramsMap map[string]string) (sql.Result, error) {
	sqlStr, params := SqlConvert(sql, paramsMap)
	return UpdateByTx(txId, dataSourceName, sqlStr, params)
}
