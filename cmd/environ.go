// This is an executable that will read a config struct file and create a default.env file
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

const (
	// flags
	inputFlag  = "input"
	outputFlag = "output"

	// regex pattern for structs
	structPattern = "type .*? struct {"
)

var envFileGeneratorFlags = flag.NewFlagSet("env-file-generator", flag.ExitOnError)

func main() {
	var (
		configFile       *os.File
		envFile          *os.File
		envFileRoot      = envFileGeneratorFlags.Lookup(outputFlag).Value.String()
		inConfig         = false
		structRegex, err = regexp.Compile(structPattern)
	)
	envFileGeneratorFlags.String(inputFlag, "./config.go", "input file containing the config structs")
	envFileGeneratorFlags.String(outputFlag, ".env", "output directory")

	err = envFileGeneratorFlags.Parse(os.Args[0:])
	if err != nil {
		envFileGeneratorFlags.Usage()
		slog.Error("failed to parse flags", "error", err)
		os.Exit(1)
	}

	// load the config struct file
	configFile, err = os.Open(envFileGeneratorFlags.Lookup("input").Value.String())
	if err != nil {
		slog.Error("failed to read the input file containing the config struct(s)", "error", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(configFile)
	scanner.Split(bufio.ScanLines)
	var configLines []string
	for scanner.Scan() {
		configLines = append(configLines, scanner.Text())
	}
	configFile.Close()

	// parse the lines
	for _, line := range configLines {
		line = strings.TrimSpace(line)
		if structRegex.MatchString(line) {
			inConfig = true
			// check if this file already exists
			var envFileName strings.Builder
			envFileName.WriteString(envFileRoot)
			envFileName.WriteString(strings.Split(line, " ")[1])
			if _, err = os.Stat(envFileName.String()); err == nil {
				// copy the existing file for safety
				envFile, err = os.Open(envFileName.String())
				defer envFile.Close()
				destFile, err := os.Create(fmt.Sprintf("old-%s", envFileName.String()))
				defer destFile.Close()
				_, err = io.Copy(destFile, envFile)
				if err != nil {
					slog.Error("failed to save old env file", "error", err)
					os.Exit(1)
				}
			}
			// open/truncate the .env file
			envFile, err = os.Create(envFileName.String())
			if err != nil {
				slog.Error("failed to open/truncate env file", "error", err)
				os.Exit(1)
			}
			defer envFile.Close()
			continue
		}
		if inConfig {
			// find the env tag, and if possible the default tag to write to the output file
			var (
				envTag   string
				envValue string
			)
			if start := strings.Index(line, `env:"`); start != -1 {
				end := strings.Index(line[start:], `"`)
				if end != -1 {
					envTag = line[start:(start + end)]
				}
			}
			if start := strings.Index(line, `default:"`); start != -1 {
				end := strings.Index(line[start:], `"`)
				if end != -1 {
					envValue = line[start:(start + end)]
				}
			}
			var envLine strings.Builder
			if envTag != "" {
				// write new line
				envLine.WriteString(envTag)
				envLine.WriteString("=")
			}
			if envValue != "" {
				envLine.WriteString(envValue)
			}
			envLine.WriteString("\n")
			if envLine.Len() != 0 {
				_, err = envFile.WriteString(envLine.String())
				if err != nil {
					slog.Error("failed to write env file line", "line", envLine, "error", err)
					os.Exit(1)
				}
			}
		}
	}
}
