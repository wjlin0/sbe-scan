package method

import "fmt"

type JolokiaList struct {
	Value map[string]map[string]info `json:"value"`
}

type info struct {
	Op map[string]interface{} `json:"op"`
}

// IsGetProperty 判断是否存在 getProperty 的 mbean
func (j *JolokiaList) IsGetProperty() bool {
	for _, values := range j.Value {
		for _, infos := range values {
			if _, ok := infos.Op["getProperty"]; ok {
				return true
			}
		}
	}
	return false
}

// GetGetProperty 获取 getProperty 的 mbean
func (j *JolokiaList) GetGetProperty() string {
	for key, values := range j.Value {
		for infosKey, infos := range values {
			if _, ok := infos.Op["getProperty"]; ok {
				return fmt.Sprintf("%s:%s", key, infosKey)
			}
		}
	}
	return ""
}
