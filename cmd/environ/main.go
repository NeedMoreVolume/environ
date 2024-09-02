// This is an executable that will read a config struct file and create a default.env file
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"
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
	inputFileName = flag.String(inputFlagName, "./config.go", "input file containing config struct(s)")
	outputDir     = flag.String(outputPathFlagName, ".env/", "output directory to save env files for config struct(s)")
	targetEnv     = flag.String(targetFlagName, "Config", "name a specific config to generate an env file for")
)

func main() {
	var (
		inputFile        *os.File
		envFile          *os.File
		inConfig         = false
		startStructRegex = regexp.MustCompile(startStructPattern)
		endStructRegex   = regexp.MustCompile(endStructPattern)
		envFilesCreated  []string
		err              error
	)
	// if the required flags are empty, print usage and exit
	if outputDir == nil || inputFileName == nil {
		flag.Usage()
		return
	}

	// open input file to read structs from
	inputFile, err = loadFile(*inputFileName)
	if err != nil {
		slog.Error("failed to load input file with config struct(s)", "error", err)
		return
	}
	// read the lines of the file into memory
	lines, err := readFile(inputFile)
	if err != nil {
		slog.Error("failed to read the file into memory", "error", err)
		return
	}
	inputFile.Close()

	// process the file
	envFileOpen := false
	for i := range lines {
		// trim the space off the line
		line := strings.TrimSpace(lines[i])
		// check if we have enetered a config struct
		if startStructRegex.MatchString(line) {
			inConfig = true
			// check if this file already exists
			envFileName := generateEnvFileName(line, *outputDir)
			if fileExists(envFileName) {
				oldCopyFileName := fmt.Sprintf("%s-%s", envFileName, time.Now().Format(time.DateOnly))
				// copy existing file to keep any settings safe
				err = copyFile(envFileName, oldCopyFileName)
				if err != nil {
					slog.Error("failed to save old env file", "env file name", envFileName, "error", err)
					break
				}
			}
			// open/truncate the .env file
			envFile, err = os.Create(envFileName)
			if err != nil {
				slog.Error("failed to open/truncate env file", "error", err)
				os.Exit(1)
			}
			envFileOpen = true
			// file will be closed once we have reached the end of this config struct
			break
		}
		// check if we have reached the end of a config struct
		if inConfig && endStructRegex.MatchString(line) {
			// close the file for writing, next run we will open a new file to write any further config env files encountered
			envFile.Close()
			// record envfile created
			envFilesCreated = append(envFilesCreated, envFile.Name())
			break
		}
		// if we are processing a config struct
		if inConfig {
			// find the env tag, and if possible the default tag to write to the output file
			tag := getEnvTagFromLine(line)
			value := getEnvValueFromLine(line)
			err := writeEnvFileLine(envFile, tag, value)
			if err != nil {
				slog.Error("failed to write new envfile line", "tag", tag, "value", value)
				break
			}
		}
	}

	// ensure the envfile written is closed
	if envFileOpen {
		err = envFile.Close()
		if err != nil {
			slog.Error("failed to close the env file", "error", err)
		}
	}

	slog.Info("completed processing config struct(s) file!", "env-files created", envFilesCreated)
}
