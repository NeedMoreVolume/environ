package cli_test

import (
	"errors"
	"log/slog"
	"os"
	"reflect"
	"testing"

	"github.com/NeedMoreVolume/environ/cli"
)

func TestNewEnvFileGenerator(t *testing.T) {
	testCases := map[string]struct {
		flags          []string
		expectedStruct *cli.EnvFileGenerator
		expectedError  error
	}{
		"with no flags": {
			expectedStruct: &cli.EnvFileGenerator{
				Logger:         slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})),
				StructFilename: "./config/config.go",
				OutputDir:      ".env/",
			},
		},
		"with all flags": {
			flags: []string{
				"-input=goodfile",
				"-output=goodoutput/",
				"-target=mysql",
				"-v=debug",
			},
			expectedStruct: &cli.EnvFileGenerator{
				Logger:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
				StructFilename:   "goodfile",
				OutputDir:        "goodoutput/",
				TargetConfigName: "mysql",
			},
		},
		"with an unknown flag": {
			flags: []string{
				"-bad",
			},
			expectedError: cli.ErrParsingFlags,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			envFileGenerator, err := cli.NewEnvFileGenerator(tc.flags)
			if !errors.Is(err, tc.expectedError) {
				slog.Error("error recieved and expected error do not match",
					"recieved", err,
					"expected", tc.expectedError,
				)
				t.Fail()
			}
			if !reflect.DeepEqual(tc.expectedStruct, envFileGenerator) {
				slog.Error("expected results do not match the results",
					"recieved", envFileGenerator,
					"expected", tc.expectedStruct,
				)
				t.Fail()
			}
		})
	}
}
