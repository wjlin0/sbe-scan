package runner

import (
	"fmt"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	fileutil "github.com/projectdiscovery/utils/file"
	"github.com/wjlin0/sbe-scan/pkg/types"
	updateutils "github.com/wjlin0/sbe-scan/pkg/update"
	"time"
)

func ParseOptions() *types.Options {
	options := &types.Options{}
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`sbe-scan is a tool to scan spring boot env.`)

	flagSet.CreateGroup("input", "Input",
		flagSet.StringSliceVarP(&options.URL, "u", "url", nil, "URL to scan", goflags.FileCommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&options.List, "list", nil, "File containing list of URLs to scan", goflags.FileCommaSeparatedStringSliceOptions),
	)
	flagSet.CreateGroup("output", "Output",
		flagSet.StringVarP(&options.OutputDir, "output", "o", "output/", "Directory to save results"),
	)
	flagSet.CreateGroup("config", "Config",
		flagSet.StringSliceVarP(&options.EnvURL, "env-url", "eu", nil, "URL to get env", goflags.FileCommaSeparatedStringSliceOptions),
		flagSet.StringSliceVarP(&options.JolokiaURL, "jolokia-url", "ju", nil, "URL to get jolokia", goflags.FileCommaSeparatedStringSliceOptions),
		flagSet.StringSliceVarP(&options.JolokiaListURL, "jolokia-list-url", "jlu", nil, "URL to get jolokia list", goflags.FileCommaSeparatedStringSliceOptions),
		flagSet.StringSliceVarP(&options.EnvName, "env-name", "en", nil, "env name to get env", goflags.FileCommaSeparatedStringSliceOptions),
		flagSet.StringSliceVarP(&options.Methods, "method", "m", nil, fmt.Sprintf("method to get env (support methods %s)", SupportMethods), goflags.FileCommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&options.Headers, "header", nil, "Headers to use for enumeration", goflags.FileCommaSeparatedStringSliceOptions),
	)
	flagSet.CreateGroup("limit", "Limit",
		flagSet.IntVar(&options.Timeout, "timeout", 10, "time to wait in seconds before timeout"),
		flagSet.IntVarP(&options.Thread, "thread", "t", 10, "Number of concurrent threads (default 10)"),
		flagSet.IntVarP(&options.RateLimit, "rate-limit", "rl", 0, "Rate limit for enumeration speed (n req/sec)"),
	)
	flagSet.CreateGroup("debug", "Debug",
		flagSet.BoolVar(&options.Debug, "debug", false, "Enable debugging"),
		flagSet.StringSliceVarP(&options.ProxyURL, "proxy", "p", nil, "list of http/socks5 proxy to use (comma separated or file input)", goflags.FileCommaSeparatedStringSliceOptions),
		flagSet.DurationVarP(&options.InputReadTimeout, "input-read-timeout", "irt", time.Duration(3*time.Minute), "timeout on input read"),
		flagSet.BoolVar(&options.DisableStdin, "no-stdin", false, "disable stdin processing"),
	)
	flagSet.CreateGroup("update", "Update",
		flagSet.CallbackVar(updateutils.GetUpdateToolCallback(repoName, version), "update", "Update tool"),
	)
	flagSet.SetCustomHelpText(`Examples:
Run sbe-scan on a single targets
	$ sbe-scan -url https://example.com
Run sbe-scan on a list of targets
	$ sbe-scan -list list.txt
Run sbe-scan on a single targets with env-url
	$ sbe-scan -url https://example.com -eu /actuator/env
Run sbe-scan on a single targets with jolokia-list-url
	$ sbe-scan -url https://example.com -jlu /actuator/jolokia/list
Run sbe-scan on a single targets a proxy server
	$ sbe-scan -url https://example.com  -proxy http://127.0.0.1:7890
	`)
	_ = flagSet.Parse()
	showBanner()
	options.SetOutput()
	options.Stdin = !options.DisableStdin && fileutil.HasStdin()
	options.InitTargets()
	err := ValidateRunEnumeration(options)
	if err != nil {
		gologger.Fatal().Msgf("Program exiting: %s", err.Error())
	}
	return options
}
