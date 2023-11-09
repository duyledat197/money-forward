package database

import "reflect"

// FieldMap returns an list field names and a list pointer values of an entity
func FieldMap(e any) ([]string, []any) {
	var fieldNames []string
	var fieldValues []any
	v := reflect.ValueOf(e).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldName := field.Tag.Get("db")
		fieldValue := v.Field(i).Addr().Interface()
		fieldNames = append(fieldNames, fieldName)
		fieldValues = append(fieldValues, fieldValue)
	}

	return fieldNames, fieldValues
}
