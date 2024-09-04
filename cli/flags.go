package cli

import "flag"

const (
	// GenerateEnvFilesCommand holds the value for the command
	GenerateEnvFilesCommand = "gen"

	// cli flags
	inputFlagName      = "input"
	outputPathFlagName = "output"
	targetFlagName     = "target"
	logLevelFlagName   = "v"
)

var (
	// flag set for generating env-files from a go file containg configu structs, continues on error so we return our own errors
	envFileGeneratorFlagSet = flag.NewFlagSet(GenerateEnvFilesCommand, flag.ContinueOnError)
	// vars to hold parsed flag values
	inputFileName    = envFileGeneratorFlagSet.String(inputFlagName, "./config/config.go", "input file containing config struct(s)")
	outputDir        = envFileGeneratorFlagSet.String(outputPathFlagName, ".env/", "output directory to save env files for config struct(s)")
	targetConfigName = envFileGeneratorFlagSet.String(targetFlagName, "", "the name of a specific config structure to generate an env file for")
	logLevel         = envFileGeneratorFlagSet.String(logLevelFlagName, "error", "logging level: debug, info, warn, error")
)
