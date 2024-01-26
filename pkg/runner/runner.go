package runner

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/projectdiscovery/gologger"
	errorutil "github.com/projectdiscovery/utils/errors"
	"github.com/remeh/sizedwaitgroup"
	"github.com/wjlin0/sbe-scan/pkg/method"
	"github.com/wjlin0/sbe-scan/pkg/method/agent/one"
	"github.com/wjlin0/sbe-scan/pkg/types"
)

type Runner struct {
	Options  *types.Options
	Sessions *method.Session
}

func NewRunner(options *types.Options) (*Runner, error) {
	sessions, err := method.NewSession(options)
	if err != nil {
		return nil, errorutil.New(err.Error())
	}
	r := &Runner{
		Sessions: sessions,
		Options:  options,
	}

	return r, nil

}
func (r *Runner) RunEnumeration() error {

	var (
		sessions = r.Sessions
		options  = r.Options
	)
	if options.Count() == 0 {
		return errorutil.New("no targets specified")
	}
	// 获取输入的env的URL，如果选择默认则 /actuator/env 与 domain 拼接
	var (
		envJson *method.Configuration
		//envURL  string
	)

	wg := sizedwaitgroup.New(options.Thread)
	for _, domain := range options.Targets {
		wg.Add()
		go func(domain string) {
			defer wg.Done()
			err := func(domain string) error {
				if options.IsEnvAuto() {

					if envJsonMap, err := r.Sessions.GetEnvJson(domain, "/env", "/actuator/env", "/actuator;..%2f..%2f/env"); err != nil {
						return err
					} else {
						for url, env := range envJsonMap {
							_ = url
							envJson = env
						}
					}

				} else {
					if envJsonMap, err := r.Sessions.GetEnvJson(domain, options.EnvURL...); err != nil {
						return err
					} else {
						for url, env := range envJsonMap {
							_ = url
							envJson = env
						}
					}
				}
				// 对domain进行md5值计算
				var md5String = func(data interface{}) string {
					hasher := md5.New()
					if v, ok := data.([]byte); ok {
						hasher.Write(v)
					} else if v, ok := data.(string); ok {
						hasher.Write([]byte(v))
					} else {
						return ""
					}

					return hex.EncodeToString(hasher.Sum(nil))
				}

				// 判断是否被屏蔽敏感数据
				if envJson.ProfilesActive() {
					// yaml序列化 到 application.properties 文件
					return envJson.WriteProperties(md5String(domain) + ".application.json")
				}

				var methods = make(map[string]struct{})
				if options.IsMethodAuto() {
					methods["one"] = struct{}{}
				} else {
					for _, m := range options.Methods {
						methods[m] = struct{}{}
					}
				}
				for m := range methods {
					switch m {
					case "one":
						// 运行方法一
						gologger.Info().Msgf("Running method one: %s", one.Describe())
						agent := &one.Agent{}
						if err := agent.Run(fmt.Sprintf(domain), sessions, options, envJson); err == nil {
							return agent.EnvJson.WriteProperties(md5String(domain) + ".application.json")
						} else {
							gologger.Debug().Msg(err.Error())
						}
					default:
						return errorutil.New("no such method")
					}
				}
				if envJson.PropertySources != nil {
					return envJson.WriteProperties(md5String(domain) + ".application.json")
				}
				return nil
			}(domain)
			if err != nil {
				gologger.Error().Msg(err.Error())
			}
		}(domain)
	}
	wg.Wait()
	return nil
}
