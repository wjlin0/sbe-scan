package main

import (
	"github.com/projectdiscovery/gologger"
	"github.com/wjlin0/sbe-scan/pkg/runner"
)

func main() {

	newRunner, err := runner.NewRunner(runner.ParseOptions())
	if err != nil {
		gologger.Fatal().Msgf("new runner error: %s", err.Error())
		return
	}
	if err := newRunner.RunEnumeration(); err != nil {
		gologger.Fatal().Msgf("run enumeration error: %s", err.Error())
	}
}
