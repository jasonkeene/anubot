package api

import "fmt"

func getStringFromMap(key string, data map[string]interface{}) (string, error) {
	d, ok := data[key]
	if !ok {
		return "", fmt.Errorf("Key %s was not provided", key)
	}
	value, ok := d.(string)
	if !ok {
		return "", fmt.Errorf("Unable to assert string type for key %s", key)
	}
	return value, nil
}
