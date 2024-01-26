package utils

import "encoding/json"

func IsJsonData(data []byte) bool {
	var d interface{}
	if err := json.Unmarshal(data, &d); err != nil {
		return false
	}
	return true
}
