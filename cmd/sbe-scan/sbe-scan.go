package main

import (
	"github.com/projectdiscovery/gologger"
	"github.com/wjlin0/sbe-scan/pkg/runner"
)

func main() {
	options, err := runner.ParseOptions()
	if err != nil {
		gologger.Error().Msgf("parse options error: %s", err.Error())
		return
	}
	newRunner, err := runner.NewRunner(options)
	if err != nil {
		gologger.Error().Msgf("parse options error: %s", err.Error())
		return
	}
	if err := newRunner.RunEnumeration(); err != nil {
		gologger.Error().Msgf("run enumeration error: %s", err.Error())
		return
	}
}
