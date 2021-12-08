package test

import (
	"github.com/yuyenews/Beerus-DB/commons/util"
	"github.com/yuyenews/Beerus-DB/db"
	"github.com/yuyenews/Beerus-DB/operation"
	"github.com/yuyenews/Beerus-DB/pool"
	"log"
	"testing"
)

func TestUpdate(t *testing.T) {
	initDbPool()

	param := make([]interface{}, 2)
	param[0] = "TestUpdate"
	param[1] = 1

	operation.GetDBTemplate("dbPoolTest").Update("update xt_message_board set user_name = ? where id = ?", param)

	param2 := make([]interface{}, 1)
	param2[0] = param[0]
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

	res := ResultStruct{}
	res.Id = 1
	res.UserName = "TestUpdateByMap"

	operation.GetDBTemplate("dbPoolTest").UpdateByMap("update xt_message_board set user_name = {user_name} where id = {id}", util.StructToMap(&res, res))

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

	res := ResultStruct{}
	res.Id = 1
	res.UserName = "TestUpdateTx"

	ss, err := operation.GetDBTemplateTx(id, "dbPoolTest").UpdateByTxMap("update xt_message_board set user_name = {user_name} where id = {id}", util.StructToMap(&res, res))
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

	id, err = db.Transaction()
	if err != nil {
		t.Error("TestUpdateTx: " + err.Error())
		return
	}

	ss, err = operation.GetDBTemplateTx(id, "dbPoolTest").UpdateByTxMap("update xt_message_board set user_name = {user_name} where id = {id}", util.StructToMap(&res, res))
	if err != nil {
		db.Rollback(id)
		t.Error("TestUpdateTx: " + err.Error())
		return
	}
	log.Println(ss.RowsAffected())

	db.Commit(id)

	param2 = make([]interface{}, 1)
	param2[0] = res.Id
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
		res := new(ResultStruct)
		util.MapToStruct(row, res, *res)

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
		res := new(ResultStruct)
		util.MapToStruct(row, res, *res)

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

	res := ResultStruct{}
	res.Id = 1
	resultMap, err := operation.GetDBTemplate("dbPoolTest").SelectListByMap("select * from xt_message_board where id < {id}", util.StructToMap(&res, res))
	if err != nil {
		t.Error("TestSelectListByMap: " + err.Error())
		return
	}

	for _, row := range resultMap {
		res := new(ResultStruct)
		util.MapToStruct(row, res, *res)

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
	dbPool.ExpSize = 1
	dbPool.MaxOpen = 1
	dbPool.MinOpen = 0
	dbPool.Url = "root:123456@(127.0.0.1:3306)/xt-manager"

	db.AddDataSource("dbPoolTest", dbPool)
}

type ResultStruct struct {
	Id         int    `field:"id"`
	UserName   string `field:"user_name"`
	UserEmail  string `field:"user_email"`
	UpdateTime string `field:"update_time"`
}
