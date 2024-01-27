package types

import (
	"bufio"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	readerutil "github.com/projectdiscovery/utils/reader"
	stringsutil "github.com/projectdiscovery/utils/strings"
	"os"
	"strings"
	"time"
)

type Options struct {
	URL              goflags.StringSlice
	List             goflags.StringSlice
	EnvURL           goflags.StringSlice
	JolokiaURL       goflags.StringSlice
	JolokiaListURL   goflags.StringSlice
	EnvName          goflags.StringSlice
	Headers          goflags.StringSlice
	ProxyURL         goflags.StringSlice
	Debug            bool
	DisableStdin     bool
	RateLimit        int
	Thread           int
	Methods          goflags.StringSlice
	OutputDir        string
	Targets          []string
	InputReadTimeout time.Duration
	Stdin            bool
	Timeout          int
}

// Count 计算目标有多少个
func (opt *Options) Count() int {
	return len(opt.Targets)
}

// InitTargets 初始化 targets
func (opt *Options) InitTargets() {
	var targetMap = make(map[string]struct{})
	var target string
	// 写一个内部函数用于调整 target
	var adjustTarget = func(target string) string {
		target = strings.TrimSpace(target)
		if target == "" {
			return ""
		}
		target = strings.TrimSuffix(target, "/")
		if !stringsutil.HasPrefixAny(target, "http://", "https://") {
			target = "http://" + target
		}
		return target
	}

	for _, target = range opt.URL {
		if target = adjustTarget(target); target == "" {
			continue
		}

		targetMap[target] = struct{}{}
	}
	for _, target = range opt.List {
		if target = adjustTarget(target); target == "" {
			continue
		}
		targetMap[target] = struct{}{}
	}
	// 从标准输入中读取
	if opt.Stdin {
		reader := readerutil.TimeoutReader{Reader: os.Stdin, Timeout: opt.InputReadTimeout}
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			if target = adjustTarget(scanner.Text()); target == "" {
				continue
			}
			targetMap[target] = struct{}{}
		}
	}
	for target = range targetMap {
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
