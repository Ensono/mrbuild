package config

type InputConfig struct {
	Config string `mapstructure:"config"`

	// version of the application
	Version string `yaml:"-"`

	Directory Directory `mapstructure:"directory"`
	Log       Log       `mapstructure:"log"`
	Projects  []Project `mapstructure:"projects"`
	Pool      Pool      `mapstructure:"pool"`
	Branch    string    `mapstructure:"branch"` // Branch that changes should be measured against
	Options   Options   `mapstructure:"options"`
	Datafile  string    `mapstructure:"datafile"`
}
