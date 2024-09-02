package main

import (
	"bufio"
	"io"
	"log/slog"
	"os"
	"strings"
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

// generates the full env file name from the config struct and the output path
func generateEnvFileName(line string, outputPath string) string {
	var envFileName strings.Builder
	envFileName.WriteString(outputPath)
	envFileName.WriteString(strings.Split(line, " ")[1])
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
