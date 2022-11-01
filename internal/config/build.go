package config

type Build struct {
	Cmd    string `mapstructure:"cmd"`
	Folder string `mapstructure:"folder"`
}
