package method

import (
	"encoding/json"
	"github.com/wjlin0/sbe-scan/pkg/types"
	"os"
)

type PropertySources struct {
	Name  string                   `json:"name"`
	Props map[string]PropertyValue `json:"properties"`
}
type PropertyValue struct {
	Value  interface{} `json:"value"`
	Origin string      `json:"origin"`
}

type Configuration struct {
	ActiveProfiles  []interface{}     `json:"activeProfiles"`
	PropertySources []PropertySources `json:"PropertySources"`
}

func (e *Configuration) ProfilesActive() bool {
	// 判断 properties 中的 value 有没有 ***,遇到第一个就返回false
	for _, source := range e.PropertySources {
		for _, props := range source.Props {
			if types.ToString(props.Value) == "******" {
				return false
			}
		}
	}

	return true
}

func (e *Configuration) GetProperties() map[string]string {
	properties := make(map[string]string)
	for _, source := range e.PropertySources {
		for key, value := range source.Props {
			properties[key] = types.ToString(value.Value)
		}
	}
	return properties
}

func (e *Configuration) WriteProperties(path string) error {
	marshal, err := e.Marshal()
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(marshal), 0644)

}

func (e *Configuration) Clone() *Configuration {
	configuration := &Configuration{}
	var propertySources []PropertySources
	for _, source := range e.PropertySources {
		var props = make(map[string]PropertyValue)
		for key, value := range source.Props {
			props[key] = value
		}
		propertySources = append(propertySources, PropertySources{
			Name:  source.Name,
			Props: props,
		})
	}
	configuration.PropertySources = propertySources
	var activeProfiles []interface{}
	for _, profile := range e.ActiveProfiles {
		activeProfiles = append(activeProfiles, profile)
	}
	configuration.ActiveProfiles = activeProfiles
	return configuration
}
func (e *Configuration) UpdatePropertiesValue(maps map[string]string) {
	for _, source := range e.PropertySources {
		for key, _ := range source.Props {
			if value, ok := maps[key]; ok {
				source.Props[key] = PropertyValue{
					Value: value,
				}
			}
		}
	}
}

// Unmarshal json反序列化
func (e *Configuration) Unmarshal(data []byte) error {

	return json.Unmarshal(data, e)
}
func (e *Configuration) Marshal() ([]byte, error) {
	if e.ActiveProfiles == nil {
		e.ActiveProfiles = make([]interface{}, 0)
	}
	if e.PropertySources == nil {
		e.PropertySources = make([]PropertySources, 0)
	}

	return json.MarshalIndent(e, "", "    ")
}
