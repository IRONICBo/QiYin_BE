package utils

import (
	"encoding/json"
	"errors"
	"reflect"
)

func StructToMapWithJSONTag(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil, errors.New("data parameter must be a struct or a pointer to struct")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func StructToMapString(data interface{}) (map[string]string, error) {
	result := make(map[string]string)

	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil, errors.New("data parameter must be a struct or a pointer to struct")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func FlattenMap(input map[string]interface{}, infix string) map[string]interface{} {
	flattened := make(map[string]interface{})
	for k, v := range input {
		switch v := v.(type) {
		case map[string]interface{}:
			dfsMap(flattened, v, k, infix)
		default:
			flattened[k] = v
		}
	}
	return flattened
}

func dfsMap(flattened, data map[string]interface{}, path string, infix string) {
	for k, v := range data {
		newPath := path + infix + k
		switch v := v.(type) {
		case map[string]interface{}:
			dfsMap(flattened, v, newPath, infix)
		default:
			flattened[newPath] = v
		}
	}
}
