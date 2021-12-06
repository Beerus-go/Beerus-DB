package test

import (
	"github.com/yuyenews/Beerus-DB/db"
	"github.com/yuyenews/Beerus-DB/pool"
)

func main() {

	dbPool := new(pool.DbPool)

	dbPool.InitialSize = 3
	dbPool.ExpSize = 2
	dbPool.MaxActive = 100
	dbPool.Url = "root:123456@(127.0.0.1:3306)/xt-manager"

	conn, err := dbPool.GetConn()
	if err != nil {
		println(err.Error())
		return
	}

	resultMap, err := db.SelectList(conn, "", nil)
	if err != nil {
		println(err.Error())
		return
	}

	for _, row := range resultMap {
		for key, val := range row {
			print(key)
			print(": ")
			print(val)
			print("\t")
		}
		print("\n")
	}
}
