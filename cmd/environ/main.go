// environ is an executable that will read a file containing config struct(s) and create env files for each struct encountered.
// exit error 1 is related to setup and opening of files
// exit error 2 is related to closing and saving of files.
package main

import (
	"flag"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

const (
	// flags
	inputFlagName      = "input"
	outputPathFlagName = "output"
	targetFlagName     = "target"
	logLevelFlagName   = "v"

	// regex pattern for structs
	startStructPattern = "type .*? struct {"
	endStructPattern   = "}"
)

var (
	inputFileName    = flag.String(inputFlagName, "./config/config.go", "input file containing config struct(s)")
	outputDir        = flag.String(outputPathFlagName, ".env/", "output directory to save env files for config struct(s)")
	targetConfigName = flag.String(targetFlagName, "", "the name of a specific config structure to generate an env file for")
	logLevel         = flag.String(logLevelFlagName, "false", "verbosity level: debug, info, warn, error. Default: error")
)

func main() {
	var (
		inputFile       *os.File
		envFile         *os.File
		hasTarget       bool
		envFileOpen     bool
		envFilesCreated []string
		err             error
		logger          *slog.Logger

		startStructRegex = regexp.MustCompile(startStructPattern)
	)

	flag.Parse()
	logger = setupLogger(*logLevel)
	// if the required flags are empty, print usage and exit
	if outputDir == nil || inputFileName == nil {
		flag.Usage()
		os.Exit(1)
	}
	// flag if we are generating a specific env file
	hasTarget = *targetConfigName != ""

	// open input file to read structs from
	inputFile, err = loadFile(*inputFileName)
	if err != nil {
		logger.Error("failed to load input file with config struct(s)",
			"error", err,
		)
		os.Exit(1)
	}

	// read the lines of the file into memory
	lines, err := readFile(inputFile)
	inputFile.Close()
	if err != nil {
		logger.Error("failed to read the file into memory",
			"error", err,
		)
		os.Exit(1)
	}
	logger.Debug("read input file",
		"# of lines", len(lines),
	)

	// process the file
	for i := range lines {
		// trim the space off the line
		line := strings.TrimSpace(lines[i])
		logger.Debug("processing",
			"index", i,
			"line", line,
			"env file open?", envFileOpen,
		)
		switch {
		case startStructRegex.MatchString(line):
			// check if this is the desired target provided
			configStructName := getConfigNameFromLine(line)
			logger.Debug("start of a struct detected",
				"struct name", configStructName,
			)
			if hasTarget {
				if configStructName != *targetConfigName {
					// skip until we find the target struct
					continue
				}
			}
			// start handling a new env file for the config struct
			envFile, err = openNewEnvFile(configStructName, outputDir)
			if err != nil {
				logger.Error("failed to open the new env file",
					"config struct", configStructName,
					"error", err,
				)
				break
			}
			envFileOpen = true
		case envFileOpen && line == endStructPattern:
			// record the name of this file and close it
			envFilesCreated = append(envFilesCreated, envFile.Name())
			err = envFile.Close()
			if err != nil {
				logger.Error("failed to close the env file",
					"env file", envFilesCreated[len(envFilesCreated)-1],
					"error", err,
				)
				break
			}
			envFileOpen = false
		case envFileOpen:
			// complete writing a new env file for a config struct
			tag := getEnvTagFromLine(line)
			value := getEnvValueFromLine(line)
			err = writeEnvFileLine(envFile, tag, value)
			if err != nil {
				logger.Error("failed to write new envfile line",
					"tag", tag,
					"value", value,
				)
				break
			}
		}
	}

	// ensure the envfile opened is closed
	if envFileOpen && envFile != nil {
		closeErr := envFile.Close()
		if closeErr != nil {
			logger.Error("failed to close the env file",
				"error", closeErr,
			)
		}
		os.Exit(2)
	}
	logger.Info("completed processing",
		slog.String("files created", strings.Join(envFilesCreated, ",")),
	)
	// exit if we errored, instead of logging success
	if err != nil {
		os.Exit(2)
	}

	logger.Info("successfully processed all configs")
}
