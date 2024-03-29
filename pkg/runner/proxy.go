package runner

import (
	"bufio"
	"fmt"
	errorutil "github.com/projectdiscovery/utils/errors"
	fileutil "github.com/projectdiscovery/utils/file"
	proxyutils "github.com/projectdiscovery/utils/proxy"
	"github.com/wjlin0/sbe-scan/pkg/types"
	"net/url"
	"os"
	"strings"
)

func loadProxyServers(options *types.Options) error {
	// TODO - Add your code here
	if len(options.ProxyURL) == 0 {
		return nil
	}
	proxyList := []string{}
	for _, p := range options.ProxyURL {
		if fileutil.FileExists(p) {
			file, err := os.Open(p)
			if err != nil {
				return fmt.Errorf("could not open proxy file: %w", err)
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				proxy := scanner.Text()
				if strings.TrimSpace(proxy) == "" {
					continue
				}
				proxyList = append(proxyList, proxy)
			}
		} else {
			proxyList = append(proxyList, p)
		}
	}
	aliveProxy, err := proxyutils.GetAnyAliveProxy(options.Timeout, proxyList...)
	if err != nil {
		return err
	}
	proxyURL, err := url.Parse(aliveProxy)
	if err != nil {
		return errorutil.WrapfWithNil(err, "failed to parse proxy got %v", err)
	}
	types.ProxyURL = proxyURL.String()
	return nil
}
