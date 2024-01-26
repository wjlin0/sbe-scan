package runner

import "github.com/projectdiscovery/gologger"

const banner = `
            __                                              
   _____   / /_   ___           _____  _____  ____ _   ____ 
  / ___/  / __ \ / _ \ ______  / ___/ / ___/ / __  /  / __ \
 (__  )  / /_/ //  __//_____/ (__  ) / /__  / /_/ /  / / / /
/____/  /_.___/ \___/        /____/  \___/  \__,_/  /_/ /_/
`
const (
	version  = `0.0.1`
	repoName = `sbe-scan`
)

// showBanner is used to show the banner to the user
func showBanner() {
	gologger.Print().Msgf("%s\n", banner)
	gologger.Print().Msgf("\t\t\t\twjlin0.com\n\n")
	gologger.Print().Msgf("慎用。你要为自己的行为负责\n")
	gologger.Print().Msgf("开发者不承担任何责任，也不对任何误用或损坏负责.\n")
}
