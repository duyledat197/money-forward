package reflect_utils

import (
	"encoding/json"
	"reflect"
)

// ConvertMapToStruct convert the map[string]any to a passing result of generic struct.
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

// CopyStruct copy value of source to destination by using json marshal and unmarshal json methods.
func CopyStruct[T, V any](source *V, destination *T) error {
	b, err := json.Marshal(source)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, destination); err != nil {
		return err
	}

	return nil
}
