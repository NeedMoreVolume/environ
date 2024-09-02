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

	// regex pattern for structs
	startStructPattern = "type .*? struct {"
	endStructPattern   = "}"
)

var (
	inputFileName    = flag.String(inputFlagName, "./config/config.go", "input file containing config struct(s)")
	outputDir        = flag.String(outputPathFlagName, ".env/", "output directory to save env files for config struct(s)")
	targetConfigName = flag.String(targetFlagName, "", "the name of a specific config structure to generate an env file for")
)

func main() {
	var (
		inputFile       *os.File
		envFile         *os.File
		hasTarget       bool
		envFileOpen     bool
		envFilesCreated []string
		err             error

		startStructRegex = regexp.MustCompile(startStructPattern)
		endStructRegex   = regexp.MustCompile(endStructPattern)
	)

	flag.Parse()
	// if the required flags are empty, print usage and exit
	if outputDir == nil || inputFileName == nil {
		flag.Usage()
		os.Exit(1)
	}
	// flag if we are generating a specific env file
	hasTarget = *targetConfigName == ""

	// open input file to read structs from
	inputFile, err = loadFile(*inputFileName)
	if err != nil {
		slog.Error("failed to load input file with config struct(s)", "error", err)
		os.Exit(1)
	}

	// read the lines of the file into memory
	lines, err := readFile(inputFile)
	inputFile.Close()
	if err != nil {
		slog.Error("failed to read the file into memory", "error", err)
		os.Exit(1)
	}

	// process the file
	for i := range lines {
		// trim the space off the line
		line := strings.TrimSpace(lines[i])
		switch {
		case startStructRegex.MatchString(line):
			// check if this is the desired target provided
			configStructName := getConfigNameFromLine(line)
			if hasTarget {
				if configStructName != *targetConfigName {
					// skip until we find the target struct
					continue
				}
			}
			// start handling a new env file for the config struct
			envFile, err = openNewEnvFile(configStructName, outputDir)
			if err != nil {
				slog.Error("failed to open the new env file", "config struct", configStructName, "error", err)
				break
			}
			envFileOpen = true
		case envFileOpen && endStructRegex.MatchString(line):
			// record the name of this file and close it
			envFilesCreated = append(envFilesCreated, envFile.Name())
			err = envFile.Close()
			if err != nil {
				slog.Error("failed to close the env file", "env file", envFilesCreated[len(envFilesCreated)-1], "error", err)
				break
			}
			envFileOpen = false
		case envFileOpen:
			// complete writing a new env file for a config struct
			tag := getEnvTagFromLine(line)
			value := getEnvValueFromLine(line)
			err = writeEnvFileLine(envFile, tag, value)
			if err != nil {
				slog.Error("failed to write new envfile line", "tag", tag, "value", value)
				break
			}
		}
	}

	// ensure the envfile opened is closed
	if envFileOpen {
		closeErr := envFile.Close()
		if closeErr != nil {
			slog.Error("failed to close the env file", "error", closeErr)
		}
		os.Exit(2)
	}
	slog.Info("completed processing", "env-files created", envFilesCreated)
	// exit if we errored, instead of logging success
	if err != nil {
		os.Exit(2)
	}

	slog.Info("successfully processed all configs")
}
