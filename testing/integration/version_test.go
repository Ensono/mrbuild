//go:build integration
// +build integration

package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// VersionSuite creates a suite of tests that check the version
// of the binary being tested is correct
type VersionSuite struct {
	BaseIntegration
}

// TestVersionSuite runs the suite of tests to check that the CLI is the correct
// version for the tests
func TestVersionSuite(t *testing.T) {

	s := new(VersionSuite)
	s.BinaryCmd = *binaryCmd
	s.ProjectDir = *projectDir

	s.Assert = assert.New(t)

	suite.Run(t, s)
}

// TestVersionNumber checks the output of the version command of the
// CLI and matches it against the version held in the integration test
func (suite *VersionSuite) TestVersionNumber() {

	// run the command and then check the output
	arguments := "version"
	suite.BaseIntegration.RunCommand(suite.BinaryCmd, arguments, false)

	suite.T().Run("CLI is the correct version", func(t *testing.T) {

		// escape the . in the version number
		escaped := strings.Replace(version, ".", `\.`, -1)
		pattern := fmt.Sprintf(`Version:\s+%s`, escaped)

		matched := suite.CheckCmdOutput(pattern)

		suite.Assert.Equal(true, matched, "Version number should be the same as the CLI being tested")
	})
}
