package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/amido/mrbuild/internal/config"
	"github.com/amido/mrbuild/internal/constants"
	"github.com/amido/mrbuild/internal/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Variable to hold the path to the configuration file
	cfgFile string

	// App holds app level objects, such as logging
	App models.App

	// Config variable to hold the model after parsing
	Config config.Config

	// Set a variable to hold the version number of the application
	version string
)

var rootCmd = &cobra.Command{
	Use:     "mrbuild",
	Short:   "Provides tools when working with mono-repos",
	Long:    "",
	Version: version,

	// Call pre-run method to unmarshal the config into the app models
	PersistentPreRun: preRun,
}

// Execute the entry point for the application
func Execute() {

	// determine if there was an error in the application
	err := rootCmd.Execute()

	if err != nil {
		log.Fatalf("%v", err)
	}
}

func init() {

	// Declare variables to accept the flags from the command line
	// or environment variables
	var logLevel string
	var logFormat string
	var logColour bool
	var logFile string

	// - options
	var cmdlog bool
	var dryrun bool

	// call the initConfig to read in any configuration files
	cobra.OnInitialize(initConfig)

	// Add flags required for the command

	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "info", "Logging Level")
	rootCmd.PersistentFlags().StringVarP(&logFormat, "logformat", "f", "text", "Logging format, text or json")
	rootCmd.PersistentFlags().BoolVarP(&logColour, "logcolour", "", true, "State if colours should be used in the text output")
	rootCmd.PersistentFlags().StringVar(&logFile, "logfile", "", "File to write logs to")

	rootCmd.PersistentFlags().BoolVar(&cmdlog, "cmdlog", false, "Specify if commands should be logged")
	rootCmd.PersistentFlags().BoolVar(&dryrun, "dryrun", false, "Perform a dryrun of the CLI. No changes will be made on disk")

	// Bind command line arguments
	viper.BindPFlags(rootCmd.Flags())

	// Configure the logging options

	viper.BindPFlag("log.format", rootCmd.PersistentFlags().Lookup("logformat"))
	viper.BindPFlag("log.colour", rootCmd.PersistentFlags().Lookup("logcolour"))
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("log.file", rootCmd.PersistentFlags().Lookup("logfile"))

	viper.BindPFlag("options.cmdlog", rootCmd.PersistentFlags().Lookup("cmdlog"))
	viper.BindPFlag("options.dryrun", rootCmd.PersistentFlags().Lookup("dryrun"))

}

// initConfig read in a config file and environment variables, if they have been set
func initConfig() {

	// Allow configuration option to be set using Environment variables
	viper.SetEnvPrefix(constants.EnvVarPrefix)

	// The configuration settings are nested
	// Change the `.` delimiter to a `_` when accessing from an environment variable
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// read in environment variables that match
	viper.AutomaticEnv()

	cmd, _, _ := rootCmd.Find(os.Args[1:])
	if cmd.Use == affectedCmd.Use {
		// set the cfgfile from Viper
		cfgFile = viper.GetString("config")

		if cfgFile != "" {

			// Set the config file parameters
			viper.SetConfigFile(cfgFile)
		}

		// Read in the configruation file
		err := viper.ReadInConfig()
		if err != nil && viper.ConfigFileUsed() != "" {
			fmt.Printf("Unable to read in configuration file: %s\n", err.Error())
			os.Exit(1)
			return
		}
	}
}

// preRun is used to ensure that dependencies are in place, such as git
func preRun(ccmd *cobra.Command, args []string) {

	err := viper.Unmarshal(&Config.Input)
	if err != nil {
		log.Fatalf("Unable to read configuration into models: %v", err)
	}

	// Ensure that the path to the configuration file is set
	Config.Self.Path = viper.ConfigFileUsed()

	// Configure application logging
	// This is done after unmarshalling of the configuration so that the
	// model values can be used rather than the strings from viper
	App.ConfigureLogging(Config.Input.Log)

	// Set the version of the app in the configuration
	Config.Input.Version = version

	//TODO: Add in check to make sure that git can be found or that the path has been
	// passed on the command line
}
