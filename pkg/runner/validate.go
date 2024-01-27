package runner

import (
	"fmt"
	fileutil "github.com/projectdiscovery/utils/file"
	stringsutil "github.com/projectdiscovery/utils/strings"
	"github.com/wjlin0/sbe-scan/pkg/types"
	"github.com/wjlin0/sbe-scan/pkg/utils"
	"strings"
	"time"
)

const (
	DefaultThread           = 10
	DefaultTimeout          = 10
	DefaultInputReadTimeout = 3 * time.Minute
	SupportMethods          = "one"
	DefaultEnvURL           = "/env,/actuator/env,/actuator;..%2f..%2f/env"
)

func ValidateRunEnumeration(options *types.Options) error {

	for _, m := range options.Methods {
		if isSupportMethods(m) {
			return fmt.Errorf("not support method %s", m)
		}
	}

	// loading the proxy server list from file or cli and test the connectivity
	if err := loadProxyServers(options); err != nil {
		return err
	}
	// 判断 output 是否存在 , 判断 output 是否为目录 ,output 是否能够写入文件
	if !fileutil.FileOrFolderExists(options.OutputDir) {
		// 创建目录
		if err := fileutil.CreateFolders(options.OutputDir); err != nil {
			return err
		}
	} else {
		// 判断是否为目录
		if fileutil.FileExists(options.OutputDir) {
			return fmt.Errorf("%s is a file, not a folder", options.OutputDir)
		}

	}
	if !utils.IsWritableDirectory(options.OutputDir) {
		return fmt.Errorf("%s is not writable", options.OutputDir)
	}

	// set default input
	if options.Thread <= 0 {
		options.Thread = DefaultThread
	}
	if options.Timeout <= 0 {
		options.Timeout = DefaultTimeout
	}
	if options.InputReadTimeout <= 0 {
		options.InputReadTimeout = DefaultInputReadTimeout
	}
	return nil
}

// isSupportMethods 支持的方法
func isSupportMethods(m string) bool {
	return stringsutil.ContainsAny(m, strings.Split(SupportMethods, ",")...)
}
