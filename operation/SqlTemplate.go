package operation

import (
	"regexp"
	"strings"
)

// SqlConvert Replace the {} placeholder in sql with the ? placeholder, and convert the parameter to []interface{}
func SqlConvert(sql string, params map[string]interface{}) (string, []interface{}) {

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
