package operation

import (
	"database/sql"
	"errors"
	"github.com/Beerus-go/Beerus-DB/db"
	"github.com/Beerus-go/Beerus-DB/operation/entity"
	"strconv"
	"strings"
)

type DBTemplate struct {
	DataSourceName string
	TxId           uint64
}

// GetDBTemplate Get templates for database operations that do not require transactions
func GetDBTemplate(dataSourceName string) *DBTemplate {
	template := new(DBTemplate)
	template.DataSourceName = dataSourceName
	return template
}

// GetDBTemplateTx Get templates for database operations that require transactions
func GetDBTemplateTx(txId uint64, dataSourceName string) *DBTemplate {
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

	if result != nil {
		if len(result) == 1 {
			return result[0], nil
		} else if len(result) > 1 {
			return nil, errors.New("this condition queries more than one data")
		}
	}

	return nil, nil
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

// Exec Add, delete and change operations
func (template *DBTemplate) Exec(sql string, params []interface{}) (sql.Result, error) {
	conn, err := db.GetConnection(template.DataSourceName)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	result, err := db.Update(conn.DB, sql, params)
	return result, err
}

// ExecNoParameters Add, delete and change operations
func (template *DBTemplate) ExecNoParameters(sql string) (sql.Result, error) {
	return template.Exec(sql, make([]interface{}, 0))
}

// ExecByMap Add, delete and change operations based on the structure parameters
func (template *DBTemplate) ExecByMap(sql string, paramsStruct map[string]interface{}) (sql.Result, error) {
	sqlStr, params := SqlConvert(sql, paramsStruct)
	return template.Exec(sqlStr, params)
}

// ----------------------- Add, delete and change operations with transactions ------------------------

// ExecByTx Add, delete and change operations with transactions
func (template *DBTemplate) ExecByTx(sql string, params []interface{}) (sql.Result, error) {
	tx, err := db.GetTx(template.TxId, template.DataSourceName)
	if err != nil {
		return nil, err
	}
	result, err := db.UpdateByTx(tx, sql, params)
	return result, err
}

// ExecByTxNoParameters Add, delete and change operations with transactions
func (template *DBTemplate) ExecByTxNoParameters(sql string) (sql.Result, error) {
	return template.ExecByTx(sql, make([]interface{}, 0))
}

// ExecByTxMap Add, delete and change operations with transactions based on structure parameters
func (template *DBTemplate) ExecByTxMap(sql string, paramsStruct map[string]interface{}) (sql.Result, error) {
	sqlStr, params := SqlConvert(sql, paramsStruct)
	return template.ExecByTx(sqlStr, params)
}

// ----------------------- Paging queries ------------------------

// SelectPage Paging queries
func (template *DBTemplate) SelectPage(sql string, pageParam entity.PageParam) (*entity.PageResult, error) {
	countSql := "select count(0) total from(" + sql + ") tb"
	return template.SelectPageCustomCount(sql, countSql, pageParam)
}

// SelectPageCustomCount Paging queries, custom countSql
func (template *DBTemplate) SelectPageCustomCount(sql string, countSql string, pageParam entity.PageParam) (*entity.PageResult, error) {
	if pageParam.Params == nil {
		pageParam.Params = make(map[string]interface{})
	}
	querySql := sql + " limit {pageOffset},{pageLength}"

	countResult, err := template.SelectOneByMap(countSql, pageParam.Params)
	if err != nil {
		return nil, err
	}

	totalStr := countResult["total"]
	if totalStr == "" {
		totalStr = "0"
	}

	total, err := strconv.Atoi(totalStr)
	if err != nil {
		return nil, err
	}

	pageParam.Params["pageOffset"] = (pageParam.CurrentPage - 1) * pageParam.PageSize
	pageParam.Params["pageLength"] = pageParam.PageSize

	resultMap, err := template.SelectListByMap(querySql, pageParam.Params)
	if err != nil {
		return nil, err
	}

	result := new(entity.PageResult)
	result.PageCount = total
	result.DataList = resultMap
	result.PageTotal = entity.CalcPageTotal(pageParam.PageSize, total)
	result.PageSize = pageParam.PageSize
	result.CurrentPage = pageParam.CurrentPage

	return result, nil
}

// ----------------------- single table no sql add delete update select ------------------------

// Select No sql query
func (template *DBTemplate) Select(tableName string, params []*entity.Condition) ([]map[string]string, error) {

	sql := new(strings.Builder)
	sql.WriteString("select * from ")
	sql.WriteString(tableName)

	if params == nil || len(params) <= 0 {
		return template.SelectListNoParameters(sql.String())
	}

	sql.WriteString(" where ")
	sqlStr, paramsArray := GetSql(sql, params)

	if sqlStr != "" && paramsArray != nil && len(paramsArray) > 0 {
		return template.SelectList(sqlStr, paramsArray)
	}

	return nil, nil
}

// SelectNoParameters No sql query no parameters
func (template *DBTemplate) SelectNoParameters(tableName string) ([]map[string]string, error) {
	return template.Select(tableName, nil)
}

// Update No sql update
func (template *DBTemplate) Update(tableName string, data map[string]interface{}, params []*entity.Condition) (sql.Result, error) {
	sqlStr, paramsArray, err := getUpdate(tableName, data, params)
	if err != nil {
		return nil, err
	}

	return template.Exec(sqlStr, paramsArray)
}

// Insert No sql insert
func (template *DBTemplate) Insert(tableName string, data map[string]interface{}) (sql.Result, error) {
	sqlStr, paramsArray, err := getInsert(tableName, data)
	if err != nil {
		return nil, err
	}

	return template.Exec(sqlStr, paramsArray)
}

// Delete No sql delete
func (template *DBTemplate) Delete(tableName string, params []*entity.Condition) (sql.Result, error) {
	sqlStr, paramsArray, err := getDelete(tableName, params)
	if err != nil {
		return nil, err
	}

	return template.Exec(sqlStr, paramsArray)
}

// UpdateTx No sql update with transactions
func (template *DBTemplate) UpdateTx(tableName string, data map[string]interface{}, params []*entity.Condition) (sql.Result, error) {
	sqlStr, paramsArray, err := getUpdate(tableName, data, params)
	if err != nil {
		return nil, err
	}

	return template.ExecByTx(sqlStr, paramsArray)
}

// InsertTx No sql insert with transactions
func (template *DBTemplate) InsertTx(tableName string, data map[string]interface{}) (sql.Result, error) {
	sqlStr, paramsArray, err := getInsert(tableName, data)
	if err != nil {
		return nil, err
	}

	return template.ExecByTx(sqlStr, paramsArray)
}

// DeleteTx No sql delete with transactions
func (template *DBTemplate) DeleteTx(tableName string, params []*entity.Condition) (sql.Result, error) {
	sqlStr, paramsArray, err := getDelete(tableName, params)
	if err != nil {
		return nil, err
	}

	return template.ExecByTx(sqlStr, paramsArray)
}

// ----------------------- Splice add, delete and update sql ------------------------

// getUpdate Return the sql of update according to the parameters
func getUpdate(tableName string, data map[string]interface{}, params []*entity.Condition) (string, []interface{}, error) {
	if params == nil || len(params) <= 0 {
		return "", nil, errors.New("in order to prevent accidents, please write your own sql to implement the update without conditions")
	}

	if data == nil || len(data) <= 0 {
		return "", nil, errors.New("data must not be empty")
	}
	sql := new(strings.Builder)
	sql.WriteString("update ")
	sql.WriteString(tableName)
	sql.WriteString(" set ")

	sqlStr, paramsArray := GetUpdateSql(sql, data, params)

	return sqlStr, paramsArray, nil
}

// getInsert Return the sql of insert according to the parameters
func getInsert(tableName string, data map[string]interface{}) (string, []interface{}, error) {
	if data == nil || len(data) <= 0 {
		return "", nil, errors.New("data must not be empty")
	}

	sql := new(strings.Builder)
	sql.WriteString("insert into ")
	sql.WriteString(tableName)

	sqlStr, paramsArray := getInsertSql(sql, data)

	return sqlStr, paramsArray, nil
}

// getDelete Return the sql of delete according to the parameters
func getDelete(tableName string, params []*entity.Condition) (string, []interface{}, error) {
	if params == nil || len(params) <= 0 {
		return "", nil, errors.New("in order to prevent accidents, please write your own sql to implement the deletion without conditions")
	}

	sql := new(strings.Builder)
	sql.WriteString("delete from ")
	sql.WriteString(tableName)
	sql.WriteString(" where ")

	sqlStr, paramsArray := GetSql(sql, params)

	return sqlStr, paramsArray, nil
}
