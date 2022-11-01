//go:build integration
// +build integration

package integration

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/amido/mrbuild/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var version string

var projectDir = flag.String("projectdir", ".", "Project Directory")
var binaryCmd = flag.String("binarycmd", "mrbuild", "Name and path of the binary to use to run the tests")

type BaseIntegration struct {
	suite.Suite

	ProjectDir string

	// the name of the command to run
	BinaryCmd string

	// CmdOutput to be used for analysis
	CmdOutput string

	Assert *assert.Assertions
}

// RunCommand provides a way for all the Integration tests to run the CLI scaffold command
// in the same way
// The command and arguments are passed as strings, and the func will split up the arguments
// and then run accordingly. The output of the command is set on the suite struct so that
// the tests can analyse it
// A third option can be provided which is the ignore parameter. If set to true then the
// function will not err on a non 0 exit code. This is so that the output of the command
// can be check to make sure that the user has been informed as to why things have not worked
func (suite *BaseIntegration) RunCommand(command string, arguments string, ignore bool) {

	// use the util function to split the arguments
	cmd, args := util.BuildCommand(command, arguments)

	// write out the command thst ius being run
	cmdlogFile := filepath.Join(suite.ProjectDir, "cmdlog.txt")

	err := ioutil.WriteFile(cmdlogFile, []byte(fmt.Sprintf("%s %s", command, arguments)), 0666)

	if err != nil {
		suite.T().Fatalf("Error writing command to log file: %s", err.Error())
	}

	// configure the exec command to execute the command
	out, err := exec.Command(cmd, args...).Output()
	if err != nil && !ignore {
		suite.T().Errorf("Error running command: %v", err)
	}
	suite.CmdOutput = string(out)
}

// CheckCmdOutput checks the output of the command against the supplied pattern and
// returns a boolean stating if that pattern has been found
func (suite *BaseIntegration) CheckCmdOutput(pattern string) bool {
	var result bool

	suite.T().Logf("Looking for pattern: %s", pattern)

	re := regexp.MustCompile(pattern)
	result = re.MatchString(suite.CmdOutput)

	return result
}
