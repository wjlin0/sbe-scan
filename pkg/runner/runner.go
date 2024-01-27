package runner

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/projectdiscovery/gologger"
	errorutil "github.com/projectdiscovery/utils/errors"
	proxyutils "github.com/projectdiscovery/utils/proxy"
	"github.com/remeh/sizedwaitgroup"
	"github.com/wjlin0/sbe-scan/pkg/method"
	"github.com/wjlin0/sbe-scan/pkg/method/agent/one"
	"github.com/wjlin0/sbe-scan/pkg/types"
	updateutils "github.com/wjlin0/sbe-scan/pkg/update"
	"net/url"
	"path/filepath"
	"strings"
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

	r.displayExecutionInfo()
	var (
		sessions = r.Sessions
		options  = r.Options
	)
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

					if envJsonMap, err := r.Sessions.GetEnvJson(domain, strings.Split(DefaultEnvURL, ",")...); err != nil {
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
					gologger.Info().Msg("find url " + domain + "  write to " + r.addOutputDirectory(md5String(domain)+".application.json"))
					// yaml序列化 到 application.properties 文件
					return envJson.WriteProperties(r.addOutputDirectory(md5String(domain) + ".application.json"))
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
						gologger.Debug().Msgf("Running method one: %s", one.Describe())
						agent := &one.Agent{}
						if err := agent.Run(fmt.Sprintf(domain), sessions, options, envJson); err == nil {
							gologger.Info().Msg("find url " + domain + "  write to " + r.addOutputDirectory(md5String(domain)+".application.json"))
							return agent.EnvJson.WriteProperties(r.addOutputDirectory(md5String(domain) + ".application.json"))
						} else {
							gologger.Debug().Msg(err.Error())
						}
					default:
						return errorutil.New("no such method")
					}
				}
				if envJson.PropertySources != nil {
					gologger.Info().Msg("find url " + domain + "  write to " + r.addOutputDirectory(md5String(domain)+".application.json"))
					return envJson.WriteProperties(r.addOutputDirectory(md5String(domain) + ".application.json"))
				}

				return nil
			}(domain)
			if err != nil {
				gologger.Debug().Msg(err.Error())
			}

		}(domain)
	}
	wg.Wait()
	return nil
}
func (r *Runner) displayExecutionInfo() {
	opts := r.Options

	latestVersion, err := updateutils.GetToolVersionCallback(repoName, repoName)()
	if err != nil {
		if opts.Debug {
			gologger.Error().Msgf("%s version check failed: %v", repoName, err.Error())
		}
	} else {
		gologger.Info().Msgf("Current %s version v%v %v", repoName, version, updateutils.GetVersionDescription(version, latestVersion))
	}
	// 展示代理
	parse, _ := url.Parse(types.ProxyURL)
	if parse.Scheme == proxyutils.HTTPS || parse.Scheme == proxyutils.HTTP {
		gologger.Info().Msgf("Using %s as proxy server", parse.String())
	}
	if parse.Scheme == proxyutils.SOCKS5 {
		gologger.Info().Msgf("Using %s as socket proxy server", parse.String())
	}
	// 展示 targets数量
	gologger.Info().Msgf("Loaded %d targets from input", opts.Count())
	// 展示 methods
	if opts.IsMethodAuto() {
		gologger.Info().Msg("Running all methods")
	} else {
		gologger.Info().Msgf("Running %s method(s)", opts.Methods)
	}

}

func (r *Runner) addOutputDirectory(path string) string {
	dir := filepath.Join(r.Options.OutputDir, path)
	abs, err := filepath.Abs(dir)
	if err != nil {
		return dir
	}
	return abs
}
