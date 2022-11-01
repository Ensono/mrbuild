package config

import "path/filepath"

type SelfConfig struct {
	CmdLogPath string
	Path       string // path to the configuration file that was read in
}

// GetDir returns the directory of the configuration files used
// It is assumed that this is in the root of the repository
func (sc *SelfConfig) GetDir() string {
	abs, _ := filepath.Abs(filepath.Dir(sc.Path))

	return abs
}
