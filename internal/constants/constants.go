package constants

const (
	// AppName states the name of the applicatiojn
	AppName string = "Mono-repo Build"

	// Set the timestamp format for logging
	LoggingTimestamp = "Mon, 02 Jan 2006 15:04:05 -0700"

	// Set the prefix that needs to be used when setting the configuration
	// using environment variables
	EnvVarPrefix = "mrbuild"

	// DefaultErrorMessage defines the default error message if one has not been set
	DefaultErrorMessage = "An error occurred in the application"

	// Set a default version to use if one is not specified at build time
	DefaultVersion = "0.0.1-workstation"
)
