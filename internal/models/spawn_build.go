package models

import (
	"strings"
)

type SpawnBuild struct {
	Name      string // Name of the project in the mono repo
	Directory string // Directory in which the the command should be run
	Command   string // Command to run
	Env       map[string]string
	Order     int
}

// GetCommand returns a single string containing the command and the arguments that should be executed
func (s *SpawnBuild) GetCommand() string {
	return s.Command
}

// GetCommandParts returns the command and the arguments to the calling function
func (s *SpawnBuild) GetCommandParts() (string, string) {
	cmdParts := strings.SplitN(s.Command, " ", 2)

	return cmdParts[0], cmdParts[1]
}
