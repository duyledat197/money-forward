package reflect_utils

import (
	"reflect"
)

func ConvertMapToStruct[T any](m map[string]any, s *T) error {
	stValue := reflect.ValueOf(s).Elem()
	sType := stValue.Type()
	for i := 0; i < sType.NumField(); i++ {
		field := sType.Field(i)
		name := field.Tag.Get("json")
		if value, ok := m[name]; ok {
			stValue.Field(i).Set(reflect.ValueOf(value))
		}
	}

	return nil
}
