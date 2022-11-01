package config

import (
	"strings"

	"github.com/amido/mrbuild/internal/util"
)

// Options holds the options for the CLI, such as turning on cmd logging
type Options struct {
	CmdLog bool   `mapstructure:"cmdlog"`
	DryRun bool   `mapstructure:"dryrun"`
	Ignore string `mapstructure:"ignore"`
}

func (o *Options) IgnoreProject(project string) bool {
	var ignore bool

	// if the ignore string is not empty slit the string using a comma
	// and see if the slice contains the project
	if o.Ignore != "" {
		ignored := strings.Split(o.Ignore, ",")

		if util.SliceContains(ignored, project) {
			ignore = true
		}
	}

	return ignore
}
