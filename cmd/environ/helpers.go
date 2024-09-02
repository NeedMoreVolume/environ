package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

// loads the file containing the config struct(s)
// returns the file, and any errors - does not defer closing the file
func loadFile(fileName string) (*os.File, error) {
	configFile, err := os.Open(fileName)
	if err != nil {
		slog.Error("failed to read the input file containing the config struct(s)", "error", err)
		return configFile, err
	}
	return configFile, err
}

// reads the opened file by line
// returns lines as []string and any error
func readFile(file *os.File) ([]string, error) {
	var (
		scanner = bufio.NewScanner(file)
		err     error
	)
	scanner.Split(bufio.ScanLines)

	var configLines []string
	for scanner.Scan() {
		configLines = append(configLines, scanner.Text())
	}
	err = file.Close()
	if err != nil {
		slog.Error("failed to close file after reading", "error", err)
	}

	return configLines, err
}

// returns true if a file exists already
func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil
}

// copies source to destination file
func copyFile(sourceFileName, destFileName string) error {
	// open sourceFile
	sourceFile, err := os.Open(sourceFileName)
	if err != nil {
		slog.Error("failed to open the source file to copy", "source file", sourceFileName, "error", err)
		return err
	}
	defer sourceFile.Close()

	// create destFile
	destFile, err := os.Create(destFileName)
	if err != nil {
		slog.Error("failed to create the dest file", "dest file", destFileName, "error", err)
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		slog.Error("failed to save copy existing file", "error", err)
	}

	return err
}

// returns the config struct name from a file line
func getConfigNameFromLine(line string) string {
	return strings.Split(line, " ")[1]
}

// output path should contain the trailing slash.
func generateEnvFileName(fileName string, outputPath string) string {
	var envFileName strings.Builder
	envFileName.WriteString(outputPath)
	envFileName.WriteString(fileName)
	return envFileName.String()
}

func getEnvTagFromLine(line string) string {
	var output string

	if start := strings.Index(line, `env:"`); start != -1 {
		end := strings.Index(line[start:], `"`)
		if end != -1 {
			output = line[start:(start + end)]
		}
	}

	return output
}

func getEnvValueFromLine(line string) string {
	var output string

	if start := strings.Index(line, `default:"`); start != -1 {
		end := strings.Index(line[start:], `"`)
		if end != -1 {
			output = line[start:(start + end)]
		}
	}

	return output
}

// function wraps creation process of an env file
//
//	copys existing env file before truncating/creating a new env file from structs.
//	files returned must be closed by the caller, this function makes no effort to ensure files opened are closed.
func openNewEnvFile(configStructName string, outputDir *string) (*os.File, error) {
	var (
		envFile *os.File
		err     error
	)

	// check if this file already exists, if exists, copy it
	envFileName := generateEnvFileName(configStructName, *outputDir)
	if fileExists(envFileName) {
		oldCopyFileName := fmt.Sprintf("%s-%s", envFileName, time.Now().Format(time.DateOnly))
		err = copyFile(envFileName, oldCopyFileName)
		if err != nil {
			slog.Error("failed to save old env file", "env file name", envFileName, "error", err)
		}
		return envFile, err
	}

	envFile, err = os.Create(envFileName)
	if err != nil {
		slog.Error("failed to open/truncate env file", "error", err)
	}
	return envFile, err
}

// writes an env file line given a envFile, tag and value to set
func writeEnvFileLine(envFile *os.File, envTag, envValue string) error {
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
		_, err := envFile.WriteString(envLine.String())
		if err != nil {
			slog.Error("failed to write env file line", "line", envLine, "error", err)
			return err
		}
	}
	return nil
}
