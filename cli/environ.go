package cli

import (
	"errors"
	"flag"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

var (
	// ErrParsingFlags is returned for errors parsing the flags for the commands
	ErrParsingFlags = errors.New("failed to parse flags")
	// ErrMissingRequirements is returned for empty input/output flags.
	ErrMissingRequirements = errors.New("please provide an input file and a directory to write env files")
)

// EnvFileGenerator is a struct containing the info the Environ Cli requires to generate env files
type EnvFileGenerator struct {
	Logger           *slog.Logger
	StructFilename   string
	OutputDir        string
	TargetConfigName string
}

// NewEnvFileGenerator sets parses flag values and returns a generator with all required infomation injected.
func NewEnvFileGenerator(args []string) (*EnvFileGenerator, error) {
	err := envFileGeneratorFlagSet.Parse(args)
	if err != nil {
		// print usage help
		envFileGeneratorFlagSet.Usage()
		return nil, ErrParsingFlags
	}
	return &EnvFileGenerator{
		Logger:           setupLogger(*logLevel),
		StructFilename:   *inputFileName,
		OutputDir:        *outputDir,
		TargetConfigName: *targetConfigName,
	}, nil
}

// CreateEnvFiles reads the go file for config structs to create env files for.
// returns the file names created for testing purposes
func (cli *EnvFileGenerator) CreateEnvFiles() ([]string, error) {
	var (
		inputFile       *os.File
		envFile         *os.File
		hasTarget       bool
		envFileOpen     bool
		envFilesCreated []string
		err             error

		startStructRegex = regexp.MustCompile(startStructPattern)
	)

	// if the required flags are empty, print usage and exit
	if cli.OutputDir == "" || cli.StructFilename == "" {
		flag.Usage()
		return envFilesCreated, ErrMissingRequirements
	}
	// flag if we are generating a specific env file
	hasTarget = cli.TargetConfigName != ""

	// open input file to read structs from
	inputFile, err = loadFile(cli.StructFilename)
	if err != nil {
		cli.Logger.Error("failed to load input file with config struct(s)",
			"error", err,
		)
		return envFilesCreated, err
	}

	// read the lines of the file into memory
	lines, err := readFile(inputFile)
	inputFile.Close()
	if err != nil {
		cli.Logger.Error("failed to read the file into memory",
			"error", err,
		)
		return envFilesCreated, err
	}
	cli.Logger.Debug("read input file",
		"# of lines", len(lines),
	)

	// process the file
	for i := range lines {
		// trim the space off the line
		line := strings.TrimSpace(lines[i])
		cli.Logger.Debug("processing",
			"index", i,
			"line", line,
			"env file open?", envFileOpen,
		)
		switch {
		case startStructRegex.MatchString(line):
			// check if this is the desired target provided
			configStructName := getConfigNameFromLine(line)
			cli.Logger.Debug("start of a struct detected",
				"struct name", configStructName,
			)
			if hasTarget {
				if configStructName != cli.TargetConfigName {
					// skip until we find the target struct
					continue
				}
			}
			// start handling a new env file for the config struct
			envFile, err = openNewEnvFile(configStructName, &cli.OutputDir)
			if err != nil {
				cli.Logger.Error("failed to open the new env file",
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
				cli.Logger.Error("failed to close the env file",
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
				cli.Logger.Error("failed to write new envfile line",
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
			cli.Logger.Error("failed to close the env file",
				"error", closeErr,
			)
		}
		return envFilesCreated, err
	}
	cli.Logger.Info("completed processing",
		slog.String("files created", strings.Join(envFilesCreated, ",")),
	)
	// exit if we errored, instead of logging success
	if err != nil {
		return envFilesCreated, err
	}

	cli.Logger.Info("successfully processed all configs")
	return envFilesCreated, nil
}
