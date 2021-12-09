package operation

import (
	"github.com/yuyenews/Beerus-DB/operation/entity"
	"log"
	"regexp"
	"strings"
)

// SqlConvert Replace the {} placeholder in sql with the ? placeholder, and convert the parameter to []interface{}
func SqlConvert(sql string, params map[string]interface{}) (string, []interface{}) {
	if params == nil {
		return "", nil
	}

	paramArray := make([]interface{}, 0)

	r := regexp.MustCompile(`{[^}]+}`)

	for {
		zw := r.FindString(sql)
		if zw == "" {
			break
		}
		sql = strings.ReplaceAll(sql, zw, "?")

		paramName := strings.ReplaceAll(zw, "{", "")
		paramName = strings.ReplaceAll(paramName, "}", "")
		paramArray = append(paramArray, params[paramName])
	}

	return sql, paramArray
}

// GetSql Splice sql with conditions
func GetSql(sql *strings.Builder, params []*entity.Condition) (string, []interface{}) {
	if params == nil {
		return "", nil
	}

	paramArray := make([]interface{}, 0)

	for _, item := range params {
		key := item.Key
		val := item.Val

		if val == nil || val == "" {
			log.Println("val is empty, already skipped")
			continue
		}

		if key == "" {
			log.Println("key is empty, already skipped")
			continue
		}

		sql.WriteString(" ")
		sql.WriteString(key)

		if val != entity.NotWhere {
			paramArray = append(paramArray, val)
		}
	}

	return sql.String(), paramArray
}

// GetUpdateSql Get the sql of update
func GetUpdateSql(sql *strings.Builder, data map[string]interface{}, params []*entity.Condition) (string, []interface{}) {
	paramArray := make([]interface{}, 0)

	first := true
	for key, value := range data {
		if value == nil || value == "" {
			continue
		}
		if first == false {
			sql.WriteString(",")
		}
		sql.WriteString(key)
		sql.WriteString("= ?")

		paramArray = append(paramArray, value)

		first = false
	}

	sql.WriteString(" where ")
	sqlStr, param := GetSql(sql, params)

	for _, val := range param {
		paramArray = append(paramArray, val)
	}

	return sqlStr, paramArray
}

// getInsertSql Get the sql of insert
func getInsertSql(sql *strings.Builder, data map[string]interface{}) (string, []interface{}) {
	paramArray := make([]interface{}, 0)

	values := new(strings.Builder)
	values.WriteString("(")

	sql.WriteString("(")

	first := true
	for key, value := range data {
		if value == nil || value == "" {
			continue
		}

		if first == false {
			sql.WriteString(",")
			values.WriteString(",")
		}
		sql.WriteString(key)
		values.WriteString("?")

		paramArray = append(paramArray, value)

		first = false
	}
	values.WriteString(")")
	sql.WriteString(")")

	sql.WriteString("values")
	sql.WriteString(values.String())

	return sql.String(), paramArray
}
