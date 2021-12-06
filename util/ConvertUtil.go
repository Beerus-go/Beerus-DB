package util

import (
	"github.com/yuyenews/Beerus-DB/commons"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const (
	Field = "field"
)

// MapToStruct map to struct
func MapToStruct(rows map[string]string, pointResult interface{}, result interface{}) {
	var paramType = reflect.TypeOf(result)
	var paramElem = reflect.ValueOf(pointResult).Elem()

	fieldNum := paramType.NumField()
	for i := 0; i < fieldNum; i++ {
		setValue(paramType, paramElem, rows, i)
	}
}

// setValue Assigning values to fields
func setValue(paramType reflect.Type, paramElem reflect.Value, rows map[string]string, i int) {
	var structField = paramType.Field(i)
	fieldName := structField.Name
	fieldTag := structField.Tag
	fieldType := structField.Type.Name()

	field := paramElem.FieldByName(fieldName)
	paramValue := rows[fieldName]

	if paramValue == "" {
		if fieldTag != "" {
			fieldParamName := fieldTag.Get(Field)
			if fieldParamName != "" {
				paramValue = rows[fieldParamName]
			}
		}
		if paramValue == "" {
			return
		}
	}

	// Unify the handling of numeric variable types to remove the bit identifiers and facilitate the following judgments
	var fType = GetFieldType(fieldType)
	if fType != "" {
		fieldType = fType
	}

	switch fieldType {
	case commons.Int:
		val, err := strconv.ParseInt(paramValue, 10, 64)
		if err != nil {
			errorPrint(fieldName, err)
			return
		}
		field.SetInt(val)
	case commons.Uint:
		val, err := strconv.ParseUint(paramValue, 10, 64)
		if err != nil {
			errorPrint(fieldName, err)
			return
		}
		field.SetUint(val)
		break
	case commons.Float:
		val, err := strconv.ParseFloat(paramValue, 64)
		if err != nil {
			errorPrint(fieldName, err)
			return
		}
		field.SetFloat(val)
		break
	case commons.Bool:
		val, err := strconv.ParseBool(paramValue)
		if err != nil {
			errorPrint(fieldName, err)
			return
		}
		field.SetBool(val)
		break
	case commons.String:
		field.SetString(paramValue)
		break
	}
}

// GetFieldType 获取字段类型
func GetFieldType(fieldType string) string {
	if strings.HasPrefix(fieldType, commons.Int) {
		return commons.Int
	}

	if strings.HasPrefix(fieldType, commons.Float) {
		return commons.Float
	}

	if strings.HasPrefix(fieldType, commons.Uint) {
		return commons.Uint
	}

	return ""
}

// errorPrint
func errorPrint(fieldName string, err error) {
	if err != nil {
		log.Println("field:" + fieldName + "Setting value Exception occurs, " + err.Error())
	}
}
