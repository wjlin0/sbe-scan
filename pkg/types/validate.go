package types

import (
	"fmt"
	"github.com/pkg/errors"
	stringsutil "github.com/projectdiscovery/utils/strings"
)

const DefaultThread = 10

func (opt *Options) Validate() error {
	if opt.URL == nil && opt.List == nil {
		return errors.New("no targets specified")
	}

	for _, m := range opt.Methods {
		if !opt.IsSupportMethods(m) {
			return fmt.Errorf("not support method %s", m)
		}
	}
	if opt.Thread <= 0 {
		opt.Thread = DefaultThread
	}

	return nil
}
func SupportMethods() []string {
	return []string{"one"}
}

// IsSupportMethods 支持的方法
func (opt *Options) IsSupportMethods(m string) bool {
	return stringsutil.ContainsAny(m, SupportMethods()...)
}
