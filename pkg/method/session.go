package method

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/corpix/uarand"
	"github.com/pkg/errors"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/ratelimit"
	"github.com/projectdiscovery/retryablehttp-go"
	errorutil "github.com/projectdiscovery/utils/errors"
	"github.com/wjlin0/sbe-scan/pkg/types"
	"github.com/wjlin0/sbe-scan/pkg/utils"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Session handles session agent sessions
type Session struct {
	Options    *types.Options
	Client     *retryablehttp.Client
	RetryMax   int
	RateLimits *ratelimit.MultiLimiter
}

func NewSession(options *types.Options) (sessions *Session, err error) {
	timeout := 30
	retryMax := 3
	Transport := &http.Transport{
		MaxIdleConns:        -1,
		MaxIdleConnsPerHost: -1,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		ResponseHeaderTimeout: time.Duration(timeout) * time.Second,
		Proxy:                 http.ProxyFromEnvironment,
	}
	httpclient := &http.Client{
		Transport: Transport,
		Timeout:   time.Duration(timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	retryablehttpOptions := retryablehttp.Options{RetryMax: retryMax}
	retryablehttpOptions.RetryWaitMax = time.Duration(timeout) * time.Second
	client := retryablehttp.NewWithHTTPClient(httpclient, retryablehttpOptions)

	sessions = &Session{
		Client:   client,
		RetryMax: retryMax,
	}
	var rateLimit *ratelimit.Options
	if options.RateLimit > 0 {
		rateLimit = &ratelimit.Options{MaxCount: uint(options.RateLimit), Key: "default", Duration: time.Second}

	} else {
		rateLimit = &ratelimit.Options{IsUnlimited: true, Key: "default"}
	}

	sessions.RateLimits, err = ratelimit.NewMultiLimiter(context.Background(), rateLimit)
	if err != nil {
		return nil, err
	}
	sessions.Options = options
	return sessions, nil
}

func (s *Session) Get(url string) (*http.Response, error) {
	_ = s.RateLimits.Take("default")
	request, err := s.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return s.Do(request)
}
func (s *Session) Do(req *retryablehttp.Request) (*http.Response, error) {
	_ = s.RateLimits.Take("default")
	return s.Client.Do(req)
}

func ReadBody(resp *http.Response) (*bytes.Buffer, error) {
	defer resp.Body.Close()
	body := bytes.Buffer{}
	_, err := io.Copy(&body, resp.Body)
	if err != nil {
		if !strings.ContainsAny(err.Error(), "tls: user canceled") {
			return nil, err
		}
	}
	return &body, nil
}

func (s *Session) GetEnvJson(domain string, urls ...string) (map[string]*Configuration, error) {

	var envJsons []map[string]*Configuration

	var wg sync.WaitGroup
	var lock sync.Mutex
	var sessions = s
	wg.Add(len(urls))
	for _, url := range urls {
		go func(url string) {
			err := func(url string) error {
				envURL := fmt.Sprintf("%s%s", domain, url)
				defer wg.Done()
				// 获取env的内容
				resp, err := sessions.Get(envURL)
				if err != nil {

					return errorutil.NewWithErr(err).Msgf("request env url failed")
				}
				if resp.StatusCode != 200 {
					sprintf := "request env " + envURL + " failed, status code is " + strconv.Itoa(resp.StatusCode) + " not 200"

					return errors.New(sprintf)
				}

				//// 判断是不是 200 且 header 中有 Content-Type: application \ json
				//if !stringsutil.ContainsAny(resp.Header.Get("Content-Type"), "application", "json") {
				//	return fmt.Errorf("request env %s failed ,  content-type is not application/json", envURL)
				//}

				body, err := ReadBody(resp)
				if err != nil {
					return errorutil.NewWithErr(err).Msgf("read env body failed")
				}
				// 判断内容是否能被json序列化
				if !utils.IsJsonData(body.Bytes()) {
					return errors.New(envURL + " is not json data")
				}
				// json 反序列化
				envJson := &Configuration{}
				if err := envJson.Unmarshal(body.Bytes()); err != nil {
					return errorutil.NewWithErr(err).Msgf("unmarshal env json failed")
				}
				if envJson != nil && len(envJson.PropertySources) > 0 {
					lock.Lock()
					envJsons = append(envJsons, map[string]*Configuration{envURL: envJson})
					lock.Unlock()
				}
				return nil
			}(url)
			if err != nil {
				gologger.Debug().Msg(err.Error())
			}
		}(url)
	}
	wg.Wait()
	if len(envJsons) == 0 {
		return nil, errors.New("request " + domain + fmt.Sprintf(" %d times, but no env json found", len(urls)))
	}

	rand.Seed(time.Now().UnixNano())
	return envJsons[rand.Intn(len(envJsons))], nil
}

func (s *Session) GetJolokiaList(domain string, urls ...string) (map[string]*JolokiaList, error) {
	var jolokiaLists []map[string]*JolokiaList

	var wg sync.WaitGroup
	var lock sync.Mutex
	var sessions = s

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			err := func(url string) error {
				jolokiaURL := fmt.Sprintf("%s%s", domain, url)

				// 获取env的内容
				resp, err := sessions.Get(jolokiaURL)
				if err != nil {

					return errorutil.NewWithErr(err).Msgf("request jolokia url failed")
				}
				if resp.StatusCode != 200 {
					errorsStr := "request jolokia " + jolokiaURL + " failed, status code is " + strconv.Itoa(resp.StatusCode) + " not 200"
					return errors.New(errorsStr)
				}

				body, err := ReadBody(resp)
				if err != nil {
					return errors.Wrap(err, "read jolokia body failed")
				}
				if !utils.IsJsonData(body.Bytes()) {
					return errors.New(jolokiaURL + " is not json data")
				}
				// json 反序列化
				jolokiaList := &JolokiaList{}
				if err := json.Unmarshal(body.Bytes(), jolokiaList); err != nil {
					return errors.Wrap(err, "unmarshal jolokia json failed")
				}
				if jolokiaList != nil && len(jolokiaList.Value) > 0 {
					lock.Lock()
					jolokiaLists = append(jolokiaLists, map[string]*JolokiaList{jolokiaURL: jolokiaList})
					lock.Unlock()
				}
				return nil
			}(url)
			if err != nil {
				gologger.Debug().Msg(err.Error())
			}
		}(url)
	}
	wg.Wait()
	if len(jolokiaLists) == 0 {
		return nil, fmt.Errorf("request %d jolokiaURL ,but no jolokia json found", len(urls))
	}

	rand.Seed(time.Now().UnixNano())
	ri := rand.Intn(len(jolokiaLists))

	return jolokiaLists[ri], nil
}

func (s *Session) NewRequest(method, url string, body interface{}) (*retryablehttp.Request, error) {
	req, err := retryablehttp.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	switch s.Options.IsNeedHeaderAdd() {
	case false:
	case true:
		for _, header := range s.Options.Headers {
			if split := strings.SplitN(header, ":", 2); len(split) == 2 {
				req.Header.Set(strings.TrimSpace(split[0]), strings.TrimSpace(split[1]))
			}
		}
	}
	// 判断 是否存在了 User-Agent
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", uarand.GetRandom())
	}
	return req, nil
}
