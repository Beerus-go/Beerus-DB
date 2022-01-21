<h1> 
    <a href="https://beeruscc.com">Beerus-DB</a> ·
    <img src="https://img.shields.io/badge/licenes-MIT-brightgreen.svg"/> 
    <img src="https://img.shields.io/badge/golang-1.17.3-brightgreen.svg"/> 
    <img src="https://img.shields.io/badge/release-tags-brightgreen.svg"/>
</h1>

Beerus-DB is a database operation framework, currently only supports Mysql,
Use [go-sql-driver/mysql] to do database connection and basic operations, 
based on this do a lot of extensions, such as, connection pool management, 
multiple data sources, transaction management, single table no sql operation, 
multiple tables and complex operations can write their own sql, sql support {} placeholder, 
can use struct as parameters to operate the database, etc.

## Installation

```shell
go get github.com/yuyenews/Beerus-DB@v1.1.3

go get github.com/go-sql-driver/mysql
```

## Documentation

[https://beeruscc.com/beerusdb](https://beeruscc.com/beerusdb)

## Examples

###  No sql additions, deletions, update and select example

***Query specified table data based on custom conditions***
```go
conditions := builder.Create().
	Add("id > ?", 10).
	Add("and (user_name = ? or age > ?)", "bee", 18).
	Add("order by create_time desc", entity.NotWhere).
	Build()

resultMap, err := operation.GetDBTemplate("Data source name").Select("table name", conditions)
```

***Update data according to conditions***

```go
// Conditions set
conditions := builder.Create().
	Add("id = ?", 1).
	Build()

// Data settings to be modified
data := ResultStruct{UserName: "TestNoSqlUpdate"}

// Execute the modification operation
result, err := operation.GetDBTemplate("Data source name").Update("table name", dbutil.StructToMapIgnore(&data, true),conditions)

```

***Deleting data based on conditions***
```go
// Set delete conditions
conditions := builder.Create().
	Add("id = ?", 2).
	Build()

// Perform a delete operation
_, err := operation.GetDBTemplate("Data source name").Delete("table name", conditions)
```

***Insert data***

```go
data := ResultStruct{
		UserName: "TestNoSqlInsert",
		UserEmail: "xxxxx@163.com",
		UpdateTime: "2021-12-09 13:50:00",
	}

result, err := operation.GetDBTemplate("Data source name").Insert("table name", dbutil.StructToMapIgnore(&data, true))

```

### Using sql to add, delete, update and select

***Add, delete, update***

sql can be any one of add, delete, modify
```go
// with struct as parameter
res := ResultStruct{Id: 1, UserName: "TestUpdateByMap"}
operation.GetDBTemplate("Data source name").ExecByMap("update xt_message_board set user_name = {user_name} where id = {id}", dbutil.StructToMap(&res))

// Using arrays as parameters
param := make([]interface{}, 2)
param[0] = "TestUpdate"
param[1] = 1

operation.GetDBTemplate("Data source name").Exec("update xt_message_board set user_name = ? where id = ?", param)

```

***Select***

Support any query sql
```go
// Using arrays as parameters
param := make([]interface{}, 1)
param[0] = 1

resultMap, err := operation.GetDBTemplate("Data source name").SelectList("select * from xt_message_board where id = ?", param)

// with struct as parameter
res := ResultStruct{Id: 1}
resultMap, err := operation.GetDBTemplate("Data source name").SelectListByMap("select * from xt_message_board where id < {id}", dbutil.StructToMap(&res))
```

***Paging queries***

```go
data := ResultStruct{
    UserName: "TestNoSqlInsert",
    UserEmail: "xxxxx@163.com",
}

param := entity.PageParam{CurrentPage: 1, PageSize: 20, Params: dbutil.StructToMap(&data)}
result, err := operation.GetDBTemplate("Data source name").SelectPage("select * from xt_message_board where user_name = {user_name} and user_email = {user_email}", param)
```

***Transaction Management***

```go
id, err := db.Transaction()
if err != nil {
    t.Error("TestUpdateTx: " + err.Error())
    return
}

res := ResultStruct{Id: 1, UserName: "TestUpdateTx"}

ss, err := operation.GetDBTemplateTx(id, "Data source name").ExecByTxMap("update xt_message_board set user_name = {user_name} where id = {id}", dbutil.StructToMap(&res))
if err != nil {
    db.Rollback(id)
    t.Error("Data source name: " + err.Error())
    return
}
log.Println(ss.RowsAffected())

db.Commmit(id)
```

## License

Beerus is [MIT licensed](https://github.com/yuyenews/Beerus-DB/blob/master/LICENSE)
