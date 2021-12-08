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

// StructToMap struct to map
func StructToMap(pointData interface{}, data interface{}) map[string]interface{} {

	result := make(map[string]interface{})

	var paramType = reflect.TypeOf(data)
	var paramElem = reflect.ValueOf(pointData).Elem()

	fieldNum := paramType.NumField()
	for i := 0; i < fieldNum; i++ {
		var structField = paramType.Field(i)
		fieldName := structField.Name
		fieldTag := structField.Tag
		fieldType := structField.Type.Name()
		field := paramElem.FieldByName(fieldName)

		if fieldTag != "" {
			fieldTagName := fieldTag.Get(Field)
			if fieldTagName != "" {
				fieldName = fieldTagName
			}
		}

		result[fieldName] = getValue(field, fieldType)
	}

	return result
}

// getValue get the value of the field
func getValue(field reflect.Value, fieldType string) interface{} {
	// Unify the handling of numeric variable types to remove the bit identifiers and facilitate the following judgments
	var fType = GetFieldType(fieldType)
	if fType != "" {
		fieldType = fType
	}

	switch fieldType {
	case commons.Int:
		return field.Int()
	case commons.Uint:
		return field.Uint()
	case commons.Float:
		return field.Float()
	case commons.Bool:
		return field.Bool()
	case commons.String:
		return field.String()
	}

	return nil
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

// GetFieldType get field type
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