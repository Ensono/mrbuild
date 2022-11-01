package config

type Project struct {
	Name     string            `mapstructure:"name"`
	Folder   string            `mapstructure:"folder"`
	Patterns []string          `mapstructure:"patterns"`
	Build    Build             `mapstructure:"build"` // Command to run if the directory contents have changed
	Env      map[string]string `mapstructure:"env"`   // list of environment variables that should be set when the command is executed
}
