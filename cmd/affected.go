package cmd

import (
	"github.com/amido/mrbuild/internal/affected"
	"github.com/amido/mrbuild/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	affectedCmd = &cobra.Command{
		Use:   "affected",
		Short: "Detect the files that have been changed",
		Long:  "",
		Run:   executeAffectedRun,

		// Execute prerun function to ensure that the datafile exists if specified
		PreRun: affectedPreRun,
	}
)

func init() {

	// declare command variables

	// - number of workers that should be in the pool
	var workers int

	// - path to data file containing sample data
	var datafile string

	// - list of projects to ignore
	var ignore string

	// add the command
	rootCmd.AddCommand(affectedCmd)

	affectedCmd.Flags().StringVar(&ignore, "ignore", "", "List of projects that should not be processed (command delimited).")
	affectedCmd.Flags().StringVar(&datafile, "datafile", "", "Path to file containing git file data to work with")
	affectedCmd.Flags().IntVar(&workers, "workers", 1, "Number of workers to spawn jobs to")

	viper.BindPFlag("options.ignore", affectedCmd.Flags().Lookup("ignore"))
	viper.BindPFlag("datafile", affectedCmd.Flags().Lookup("datafile"))
	viper.BindPFlag("workers", affectedCmd.Flags().Lookup("workers"))

}

func affectedPreRun(ccmd *cobra.Command, args []string) {

	// check to see if the data file exists, error if not
	if Config.Input.Datafile != "" && !util.Exists(Config.Input.Datafile) {
		App.Logger.Fatalf("Specified data file cannot be found: %s", Config.Input.Datafile)
	}
}

func executeAffectedRun(ccmd *cobra.Command, args []string) {

	// if running in DryRun mode then set the workers to 1
	if Config.IsDryRun() {
		App.Logger.Warn("Set workers to 1 as running in DryRun mode")
		Config.Input.Pool.Workers = 1
	}

	// Call the affected method
	affected := affected.New(&App, &Config, App.Logger)
	err := affected.Run()
	if err != nil {
		App.Logger.Fatalf("Error running command: %s", err.Error())
	}
}
