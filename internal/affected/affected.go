package affected

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/amido/mrbuild/internal/config"
	"github.com/amido/mrbuild/internal/models"
	"github.com/amido/mrbuild/internal/util"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type Affected struct {
	App    *models.App
	Config *config.Config
	Logger *logrus.Logger
}

// New allocates a new AffectedPointer to the given config
func New(app *models.App, conf *config.Config, logger *logrus.Logger) *Affected {
	return &Affected{
		App:    app,
		Config: conf,
		Logger: logger,
	}
}

// Run performs the operations of the affected command which is to
// look decide which folder have changed and then perform a build
// within those folders according\ to the config file
func (a *Affected) Run() error {
	var err error

	// check the runtime configuration and set defaults
	err = a.Config.Check()
	if err != nil {
		a.Logger.Fatalln(err.Error())
	}

	// if a datafile has been specified, read in the data
	// otherwise run the git command to get a list of the changed files
	list, err := a.getFiles()

	// determine if any of the files match the patterns specified for the project
	affectedProjects := a.getProjects(list)

	a.App.Logger.Debugf("Analysing %d projects", len(affectedProjects))

	// Configure the worker pool
	// As each project will have its own build mechanism a pool of workers is setup to run
	// each build on a concurrent thread
	a.App.ConfigureWorkers(a.Config.Input.Pool.Workers)

	// iterate around the affected projects and spawn each build process
	for _, p := range affectedProjects {

		// Output the command that is to be run along with the directory it will be run in
		a.App.Logger.WithFields(
			log.Fields{
				"workingDir": p.Directory,
				"project":    p.Name,
				"command":    p.GetCommand(),
			},
		).Info("Executing command")

		if a.Config.IsDryRun() {
			a.App.Logger.Warn("Not running command as in DryRun mode")
		} else {

			// get the command parts
			cmd, args := p.GetCommandParts()

			// submit the command to be run
			a.App.Workers.Submit(func() {
				output, err := a.Config.ExecuteCommand(
					p.Directory,
					a.Logger,
					cmd,
					args,
					true,
					false,
				)

				if err == nil {
					a.App.Logger.Info(output)
				} else {
					a.App.Logger.Error(err.Error())
				}
			})
		}
	}

	// wait for all the jobs to complete
	a.App.Workers.StopWait()

	return err
}

// getFiles returns a list of files that are affected in this branch
// this can be done by reading the datafile, if it has been specified or by
// running the git command to get the list
func (a *Affected) getFiles() (string, error) {

	var content []byte
	var files string
	var err error

	// if running in pipe mode get the data from stdnin
	if util.IsInputFromPipe() {

		data, err := io.ReadAll(os.Stdin)

		if err != nil {
			a.Logger.Errorf("Unable to read content from pipe")
		} else {
			files = string(data)
		}
	} else if a.Config.Input.Datafile == "" {

		// execute the command
		files, err = a.Config.ExecuteCommand(
			a.Config.Input.Directory.WorkingDir,
			a.Logger,
			"git",
			fmt.Sprintf("--no-pager diff --name-only %s", a.Config.Input.Branch),
			false,
			false,
		)

		if err != nil {
			a.Logger.Errorf("Issue running command: %s", err.Error())
		}
	} else {

		// attempt to read in the file
		content, err = ioutil.ReadFile(a.Config.Input.Datafile)
		files = string(content)
	}

	return files, err
}

// getProjects iterates around the projects that have been defined in the configuration
// file and determine if any of the them have been changed
// If they have then find the command for the project and add to an array along
// with the project directory
func (a *Affected) getProjects(list string) []models.SpawnBuild {

	var re *regexp.Regexp
	var spawns []models.SpawnBuild

	// determine the path to the project, this is based on the location of the

	// iterate around the projects
	for _, project := range a.Config.Input.Projects {

		// check to see if hte project is to be ignored
		if a.Config.Input.Options.IgnoreProject(project.Name) {
			a.App.Logger.Warnf("Ignoring project: %s", project.Name)
			continue
		}

		// use the project folder and the patterns to try and match with the data
		for _, pattern := range project.Patterns {

			_pattern := fmt.Sprintf("(?m)%s/%s", project.Folder, pattern)

			re = regexp.MustCompile(_pattern)

			// determine if the regular expression matches the list of files
			matches := re.MatchString(list)
			if matches {

				// determine the path that the build should be run in
				folder := project.Build.Folder

				// if folder is . then set as the path to the configuration file
				// if it is null then use the project folder
				if folder == "" {
					folder = project.Folder
				} else if folder == "." {
					folder = a.Config.Self.GetDir()
				}

				spawns = append(spawns, models.SpawnBuild{
					Name:      project.Name,
					Command:   project.Build.Cmd,
					Directory: folder,
					Env:       project.Env,
				})

				// as a match has been found, exit out of the inner loop and move
				// onto the next project
				break
			}
		}
	}

	return spawns
}
