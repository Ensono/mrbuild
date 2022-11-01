package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/amido/mrbuild/internal/constants"
	"github.com/amido/mrbuild/internal/util"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Input InputConfig
	Self  SelfConfig

	// allow environment variables to be set
	envvars map[string]string
}

// Check ensures that there are sensible defaults for values
// that may not have been set in the input, e.g. the branch to compare against
func (c *Config) Check() error {
	var err error

	// iterate around the projects and check that branch has a default value
	if strings.TrimSpace(c.Input.Branch) == "" {
		c.Input.Branch = "main"
	}

	// set necessary default values
	c.SetDefaultValues()

	return err
}

// IsDryRun returns the boolean value of the dryrun option
func (c *Config) IsDryRun() bool {
	return c.Input.Options.DryRun
}

// UseCmdLog states of the command log should be used
func (c *Config) UseCmdLog() bool {
	return c.Input.Options.CmdLog
}

func (c *Config) GetEnvironmentVars(envvars map[string]string) {
	c.envvars = envvars
}

// GetVersion returns the current version of the application
// It will check to see uif the Version is empty, if it is, it will
// set and identifiable local build version
func (config *Config) GetVersion() string {
	var version string

	version = config.Input.Version

	if version == "" {
		version = constants.DefaultVersion
	}

	return strings.TrimSpace(version)
}

// ExecuteCommand executes the command and arguments that have been supplied to the function
func (config *Config) ExecuteCommand(path string, logger *logrus.Logger, command string, arguments string, show bool, force bool) (string, error) {

	var result bytes.Buffer
	var err error
	var mwriter io.Writer
	var writers []io.Writer

	// get the command and arguments
	cmd, args := util.BuildCommand(command, arguments)

	// output the command being run if in debug mode
	logger.Debugf("Command: %s %s", command, arguments)

	// Write out the command log
	if path != "" {
		err = config.WriteCmdLog(path, fmt.Sprintf("%s %s", command, arguments))
		if err != nil {
			logger.Warnf("Unable to write command to log: %s", err.Error())
		}
	}

	// add the result to the writers
	writers = append(writers, &result)

	// add the stdout to the multiwriter if being displayed
	if show {
		writers = append(writers, os.Stdout)
	}

	// add stderr to the mwriter, if running in loglevel greater than info
	levels := []string{"debug", "trace"}
	if util.SliceContains(levels, strings.ToLower(logger.GetLevel().String())) {

		// set the logger as a writer, this is so that errors from any commands that are
		// run are added to the file as well (if one has been set)
		w := logger.WriterLevel(logrus.DebugLevel)
		defer w.Close()
		writers = append(writers, w)
	}

	mwriter = io.MultiWriter(writers...)

	// set the command that needs to be executed
	cmdLine := exec.Command(cmd, args...)
	cmdLine.Stdout = mwriter
	cmdLine.Stderr = mwriter

	// determine if any environment variables need to be set
	if len(config.envvars) > 0 {

		// ensure that the current envvars are preserved
		cmdLine.Env = os.Environ()

		// iterate around the envvars
		for name, value := range config.envvars {
			cmdLine.Env = append(cmdLine.Env, fmt.Sprintf("%s=%s", name, value))
		}
	}

	// set the path for the command, if it exists
	if util.Exists(path) {
		cmdLine.Dir = path
	}

	// only run the command if not in dryrun mode
	// or if the force option has been set, this is for non-destructive commands such as checking the version of
	// a command
	if !config.IsDryRun() || force {
		if err = cmdLine.Run(); err != nil {
			logger.Errorf("Error running command: %s", err.Error())
			return strings.TrimSpace(result.String()), err
		}
	}

	return strings.TrimSpace(result.String()), err
}

func (config *Config) WriteCmdLog(path string, cmd string) error {

	var err error

	// return empty error if not logging commands
	if !config.UseCmdLog() {
		return err
	}

	// get a reference to the file, either to create or append to the file
	f, err := os.OpenFile(config.Self.CmdLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	// write out the cmd to the file
	if _, err := f.WriteString(fmt.Sprintf("[%s] %s\n", path, cmd)); err != nil {
		return err
	}

	return err
}

// SetDefaultValues sets values in the config object that are based off other values in the
// config object
func (config *Config) SetDefaultValues() {

	// Set the currentdirectory to the path that the CLI is currently running in
	cwd, _ := os.Getwd()
	config.Input.Directory.WorkingDir = cwd

	config.Self.CmdLogPath = filepath.Join(config.Input.Directory.WorkingDir, "cmdlog.txt")
}
