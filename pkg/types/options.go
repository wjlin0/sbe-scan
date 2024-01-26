package types

import (
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	stringsutil "github.com/projectdiscovery/utils/strings"
	"strings"
)

type Options struct {
	URL            goflags.StringSlice
	List           goflags.StringSlice
	EnvURL         goflags.StringSlice
	JolokiaURL     goflags.StringSlice
	JolokiaListURL goflags.StringSlice
	EnvName        goflags.StringSlice
	Headers        goflags.StringSlice
	Debug          bool
	RateLimit      int
	Thread         int
	Methods        goflags.StringSlice
	Targets        []string
}

// Count 计算目标有多少个
func (opt *Options) Count() int {
	return len(opt.Targets)
}

// InitTargets 初始化 targets
func (opt *Options) InitTargets() {
	var targetMap = make(map[string]struct{})
	for _, url := range opt.URL {
		// 除去URL末尾的 /
		url = strings.TrimSuffix(url, "/")
		// 添加上 http://
		if !stringsutil.HasPrefixAny(url, "http://", "https://") {
			url = "http://" + url
		}
		targetMap[url] = struct{}{}
	}
	for _, list := range opt.List {
		// 除去URL末尾的 /
		list = strings.TrimSuffix(list, "/")
		// 添加上 http://
		if !stringsutil.HasPrefixAny(list, "http://", "https://") {
			list = "http://" + list
		}
		targetMap[list] = struct{}{}
	}
	for target := range targetMap {
		opt.Targets = append(opt.Targets, target)
	}
}

// IsEnvAuto Env自动获取
func (opt *Options) IsEnvAuto() bool {
	return len(opt.EnvURL) == 0
}

// IsMethodAuto Method自动获取
func (opt *Options) IsMethodAuto() bool {
	return len(opt.Methods) == 0
}

// IsJolokiaAuto Jolokia自动获取
func (opt *Options) IsJolokiaAuto() bool {
	return len(opt.JolokiaURL) == 0
}

// IsJolokiaListAuto JolokiaList自动获取
func (opt *Options) IsJolokiaListAuto() bool {
	return len(opt.JolokiaListURL) == 0
}

// IsEnvNameAuto EnvName自动获取
func (opt *Options) IsEnvNameAuto() bool {
	return len(opt.EnvName) == 0
}

// IsNeedHeaderAdd Header是否需要添加
func (opt *Options) IsNeedHeaderAdd() bool {
	return len(opt.Headers) != 0
}
func (opt *Options) SetOutput() {
	switch {
	case opt.Debug:
		gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
	}
}

func (opt *Options) CheckVersion(callback func()) {
	callback()
}
