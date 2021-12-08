package operation

import (
	"database/sql"
	"errors"
	"github.com/yuyenews/Beerus-DB/db"
)

type DBTemplate struct {
	DataSourceName string
	TxId           string
}

// GetDBTemplate Get templates for database operations that do not require transactions
func GetDBTemplate(dataSourceName string) *DBTemplate {
	template := new(DBTemplate)
	template.DataSourceName = dataSourceName
	return template
}

// GetDBTemplateTx Get templates for database operations that require transactions
func GetDBTemplateTx(txId string, dataSourceName string) *DBTemplate {
	template := new(DBTemplate)
	template.TxId = txId
	template.DataSourceName = dataSourceName
	return template
}

// ------------------------ Query operation ----------------------

// SelectList Query and return a list
func (template *DBTemplate) SelectList(sql string, params []interface{}) ([]map[string]string, error) {
	conn, err := db.GetConnection(template.DataSourceName)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	if conn.Status == false {
		return nil, errors.New("the connection has been closed")
	}
	result, err := db.SelectList(conn.DB, sql, params)
	return result, err
}

// SelectListNoParameters Query and return a list
func (template *DBTemplate) SelectListNoParameters(sql string) ([]map[string]string, error) {
	return template.SelectList(sql, make([]interface{}, 0))
}

// SelectListByMap Query by structure parameters
func (template *DBTemplate) SelectListByMap(sql string, paramsStruct map[string]interface{}) ([]map[string]string, error) {
	sqlStr, params := SqlConvert(sql, paramsStruct)
	return template.SelectList(sqlStr, params)
}

// SelectOne Query and return a piece of data
func (template *DBTemplate) SelectOne(sql string, params []interface{}) (map[string]string, error) {
	result, err := template.SelectList(sql, params)
	if err != nil {
		return nil, err
	}

	if result != nil && len(result) > 0 {
		return result[0], nil
	}

	return nil, errors.New("this condition queries more than one data")
}

// SelectOneNoParameters Query and return a piece of data
func (template *DBTemplate) SelectOneNoParameters(sql string) (map[string]string, error) {
	return template.SelectOne(sql, make([]interface{}, 0))
}

// SelectOneByMap Query and return a piece of data by structure parameters
func (template *DBTemplate) SelectOneByMap(sql string, paramsStruct map[string]interface{}) (map[string]string, error) {
	sqlStr, params := SqlConvert(sql, paramsStruct)
	return template.SelectOne(sqlStr, params)
}

// ----------------------- Add, delete and change without transactions ------------------------

// Update Add, delete and change operations
func (template *DBTemplate) Update(sql string, params []interface{}) (sql.Result, error) {
	conn, err := db.GetConnection(template.DataSourceName)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	result, err := db.Update(conn.DB, sql, params)
	return result, err
}

// UpdateNoParameters Add, delete and change operations
func (template *DBTemplate) UpdateNoParameters(sql string) (sql.Result, error) {
	return template.Update(sql, make([]interface{}, 0))
}

// UpdateByMap Add, delete and change operations based on the structure parameters
func (template *DBTemplate) UpdateByMap(sql string, paramsStruct map[string]interface{}) (sql.Result, error) {
	sqlStr, params := SqlConvert(sql, paramsStruct)
	return template.Update(sqlStr, params)
}

// ----------------------- Add, delete and change operations with transactions ------------------------

// UpdateByTx Add, delete and change operations with transactions
func (template *DBTemplate) UpdateByTx(sql string, params []interface{}) (sql.Result, error) {
	tx, err := db.GetTx(template.TxId, template.DataSourceName)
	if err != nil {
		return nil, err
	}
	result, err := db.UpdateByTx(tx, sql, params)
	return result, err
}

// UpdateByTxNoParameters Add, delete and change operations with transactions
func (template *DBTemplate) UpdateByTxNoParameters(sql string) (sql.Result, error) {
	return template.UpdateByTx(sql, make([]interface{}, 0))
}

// UpdateByTxMap Add, delete and change operations with transactions based on structure parameters
func (template *DBTemplate) UpdateByTxMap(sql string, paramsStruct map[string]interface{}) (sql.Result, error) {
	sqlStr, params := SqlConvert(sql, paramsStruct)
	return template.UpdateByTx(sqlStr, params)
}
