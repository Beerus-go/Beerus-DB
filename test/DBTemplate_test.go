package test

import (
	"github.com/yuyenews/Beerus-DB/commons/builder"
	"github.com/yuyenews/Beerus-DB/commons/dbutil"
	"github.com/yuyenews/Beerus-DB/db"
	"github.com/yuyenews/Beerus-DB/operation"
	"github.com/yuyenews/Beerus-DB/operation/entity"
	"github.com/yuyenews/Beerus-DB/pool"
	"log"
	"strconv"
	"testing"
)

func TestUpdate(t *testing.T) {
	initDbPool()
	snowflake, err := dbutil.New(5)
	snowflakeId, err := snowflake.Generate()

	param := make([]interface{}, 2)
	param[0] = strconv.FormatUint(snowflakeId, 10)
	param[1] = 1

	operation.GetDBTemplate("dbPoolTest").Exec("update xt_message_board set user_name = ? where id = ?", param)

	param2 := make([]interface{}, 1)
	param2[0] = param[1]
	result, err := operation.GetDBTemplate("dbPoolTest").SelectOne("select * from xt_message_board where id = ?", param2)

	if err != nil {
		t.Error("TestUpdate: " + err.Error())
		return
	}

	if result["user_name"] != param[0] {
		t.Error("TestUpdate: Failed to modify, after modification, user_name is not equal to the specified value")
	}
}

func TestUpdateByMap(t *testing.T) {
	initDbPool()

	snowflake, err := dbutil.New(5)
	snowflakeId, err := snowflake.Generate()

	res := ResultStruct{Id: 1, UserName: strconv.FormatUint(snowflakeId, 10)}

	operation.GetDBTemplate("dbPoolTest").ExecByMap("update xt_message_board set user_name = {user_name} where id = {id}", dbutil.StructToMap(&res))

	param2 := make([]interface{}, 1)
	param2[0] = res.Id
	result, err := operation.GetDBTemplate("dbPoolTest").SelectOne("select * from xt_message_board where id = ?", param2)

	if err != nil {
		t.Error("TestUpdateByMap: " + err.Error())
		return
	}

	if result["user_name"] != res.UserName {
		t.Error("TestUpdateByMap: Failed to modify, after modification, user_name is not equal to the specified value")
	}
}

func TestUpdateTx(t *testing.T) {
	initDbPool()

	// ------------------------------- test rollback -------------------------------

	id, err := db.Transaction()
	if err != nil {
		t.Error("TestUpdateTx: " + err.Error())
		return
	}

	snowflake, err := dbutil.New(5)
	snowflakeId, err := snowflake.Generate()

	res := ResultStruct{Id: 1, UserName: strconv.FormatUint(snowflakeId, 10)}

	ss, err := operation.GetDBTemplateTx(id, "dbPoolTest").ExecByTxMap("update xt_message_board set user_name = {user_name} where id = {id}", dbutil.StructToMap(&res))
	if err != nil {
		db.Rollback(id)
		t.Error("TestUpdateTx: " + err.Error())
		return
	}
	log.Println(ss.RowsAffected())

	db.Rollback(id)

	param2 := make([]interface{}, 1)
	param2[0] = res.Id
	result, err := operation.GetDBTemplate("dbPoolTest").SelectOne("select * from xt_message_board where id = ?", param2)

	if err != nil {
		t.Error("TestUpdateTx: " + err.Error())
		return
	}

	if result["user_name"] == res.UserName {
		t.Error("TestUpdateTx: Transaction rollback failure")
		return
	}

	// ------------------------------- test commit -------------------------------

	snowflakeId, err = snowflake.Generate()

	res.UserName = strconv.FormatUint(snowflakeId, 10)

	id, err = db.Transaction()
	if err != nil {
		t.Error("TestUpdateTx: " + err.Error())
		return
	}

	ss, err = operation.GetDBTemplateTx(id, "dbPoolTest").ExecByTxMap("update xt_message_board set user_name = {user_name} where id = {id}", dbutil.StructToMap(&res))
	if err != nil {
		db.Rollback(id)
		t.Error("TestUpdateTx: " + err.Error())
		return
	}
	log.Println(ss.RowsAffected())

	param2 = make([]interface{}, 1)
	param2[0] = res.Id
	result, err = operation.GetDBTemplate("dbPoolTest").SelectOne("select * from xt_message_board where id = ?", param2)

	if err != nil {
		t.Error("TestUpdateTx: " + err.Error())
		return
	}

	if result["user_name"] == res.UserName {
		t.Error("TestUpdateTx: The transaction has taken effect before it has been committed")
		return
	}

	db.Commit(id)

	result, err = operation.GetDBTemplate("dbPoolTest").SelectOne("select * from xt_message_board where id = ?", param2)

	if err != nil {
		t.Error("TestUpdateTx: " + err.Error())
		return
	}

	if result["user_name"] != res.UserName {
		t.Error("TestUpdateTx: Transaction commit failure")
		return
	}
}

func TestSelectList(t *testing.T) {
	initDbPool()

	param := make([]interface{}, 1)
	param[0] = 1

	resultMap, err := operation.GetDBTemplate("dbPoolTest").SelectList("select * from xt_message_board where id = ?", param)
	if err != nil {
		t.Error("TestSelectList: " + err.Error())
		return
	}
	for _, row := range resultMap {
		res := ResultStruct{}
		dbutil.MapToStruct(row, &res)

		print(res.Id)
		print(" | ")
		print(res.UserName)
		print(" | ")
		print(res.UpdateTime)
		print(" | ")
		println(res.UserEmail)
	}
}

func TestSelectListNoParameters(t *testing.T) {
	initDbPool()

	resultMap, err := operation.GetDBTemplate("dbPoolTest").SelectListNoParameters("select * from xt_message_board")
	if err != nil {
		t.Error("TestSelectListNoParameters: " + err.Error())
		return
	}
	for _, row := range resultMap {
		res := ResultStruct{}
		dbutil.MapToStruct(row, &res)

		print(res.Id)
		print(" | ")
		print(res.UserName)
		print(" | ")
		print(res.UpdateTime)
		print(" | ")
		println(res.UserEmail)
	}
}

func TestSelectListByMap(t *testing.T) {
	initDbPool()

	res := ResultStruct{Id: 1}

	resultMap, err := operation.GetDBTemplate("dbPoolTest").SelectListByMap("select * from xt_message_board where id < {id}", dbutil.StructToMap(&res))
	if err != nil {
		t.Error("TestSelectListByMap: " + err.Error())
		return
	}

	for _, row := range resultMap {
		res := ResultStruct{}
		dbutil.MapToStruct(row, &res)

		print(res.Id)
		print(" | ")
		print(res.UserName)
		print(" | ")
		print(res.UpdateTime)
		print(" | ")
		println(res.UserEmail)
	}
}

func TestNoSqlSelect(t *testing.T) {
	initDbPool()

	conditions := builder.Create().
		Add("id > ? and user_name = ?", 10, "testTx").
		Add("and id < ?", 23).
		Add("order by id desc", entity.NotWhere).
		Build()

	resultMap, err := operation.GetDBTemplate("dbPoolTest").Select("xt_message_board", conditions)

	if err != nil {
		t.Error("TestSelectListByMap: " + err.Error())
		return
	}

	for _, row := range resultMap {
		res := ResultStruct{}
		dbutil.MapToStruct(row, &res)

		print(res.Id)
		print(" | ")
		print(res.UserName)
		print(" | ")
		print(res.UpdateTime)
		print(" | ")
		println(res.UserEmail)
	}
}

func TestNoSqlUpdate(t *testing.T) {
	initDbPool()

	snowflake, err := dbutil.New(5)
	snowflakeId, err := snowflake.Generate()

	conditions := builder.Create().
		Add("id = ?", 1).
		Build()

	data := ResultStruct{UserName: strconv.FormatUint(snowflakeId, 10)}
	operation.GetDBTemplate("dbPoolTest").Update("xt_message_board", dbutil.StructToMapIgnore(&data, true), conditions)

	param2 := make([]interface{}, 1)
	param2[0] = 1
	result, err := operation.GetDBTemplate("dbPoolTest").SelectOne("select * from xt_message_board where id = ?", param2)

	if err != nil {
		t.Error("TestNoSqlUpdate: " + err.Error())
		return
	}

	if result["user_name"] != data.UserName {
		t.Error("TestNoSqlUpdate: Failed to modify, after modification, user_name is not equal to the specified value")
	}
}

func TestNoSqlInsert(t *testing.T) {
	initDbPool()

	data := ResultStruct{
		UserName:   "TestNoSqlInsert",
		UserEmail:  "xxxxx@163.com",
		UpdateTime: "2021-12-09 13:50:00",
	}

	result, err := operation.GetDBTemplate("dbPoolTest").Insert("xt_message_board", dbutil.StructToMapIgnore(&data, true))
	if err != nil {
		t.Error("TestNoSqlInsert: " + err.Error())
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Error("TestNoSqlInsert: " + err.Error())
		return
	}

	param2 := make([]interface{}, 1)
	param2[0] = id
	result2, err := operation.GetDBTemplate("dbPoolTest").SelectOne("select * from xt_message_board where id = ?", param2)

	if err != nil {
		t.Error("TestNoSqlInsert: " + err.Error())
		return
	}

	if result2["user_name"] != data.UserName {
		t.Error("TestNoSqlInsert: Failed to insert, Unsuccessful data insertion")
	}
}

func TestNoSqlDelete(t *testing.T) {
	initDbPool()

	conditions := builder.Create().
		Add("id = ?", 14).
		Build()

	_, err := operation.GetDBTemplate("dbPoolTest").Delete("xt_message_board", conditions)

	if err != nil {
		t.Error("TestNoSqlDelete: " + err.Error())
		return
	}

	param2 := make([]interface{}, 1)
	param2[0] = 2
	result2, err := operation.GetDBTemplate("dbPoolTest").SelectOne("select * from xt_message_board where id = ?", param2)

	if err != nil {
		t.Error("TestNoSqlInsert: " + err.Error())
		return
	}

	if result2 != nil {
		t.Error("TestNoSqlInsert: Data deletion failure")
		return
	}
}

func TestSelectPage(t *testing.T) {
	initDbPool()

	param := entity.PageParam{CurrentPage: 2, PageSize: 5}
	result, err := operation.GetDBTemplate("dbPoolTest").SelectPage("select * from xt_message_board", param)

	if err != nil {
		t.Error("TestSelectPage: " + err.Error())
		return
	}

	for _, row := range result.DataList {
		res := ResultStruct{}
		dbutil.MapToStruct(row, &res)

		print(res.Id)
		print(" | ")
		print(res.UserName)
		print(" | ")
		print(res.UpdateTime)
		print(" | ")
		println(res.UserEmail)
	}
}

func initDbPool() {
	dbPool := new(pool.DbPool)
	dbPool.InitialSize = 1
	dbPool.ExpandSize = 1
	dbPool.MaxOpen = 1
	dbPool.MinOpen = 0
	dbPool.Url = "root:123456@(127.0.0.1:3306)/test"

	db.AddDataSource("dbPoolTest", dbPool)
}

type ResultStruct struct {
	Id         int    `field:"id" ignore:"true"`
	UserName   string `field:"user_name"`
	UserEmail  string `field:"user_email"`
	UpdateTime string `field:"update_time"`
}
