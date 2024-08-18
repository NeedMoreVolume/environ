package environ_test

import (
	"environ"
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

type exampleDefaultConfig struct {
	// INTS
	Int   int   `env:"MY_INT" default:"1"`
	Int8  int8  `env:"MY_INT_8" default:"1"`
	Int16 int16 `env:"MY_INT_16" default:"1"`
	Int32 int32 `env:"MY_INT_32" default:"1"`
	Int64 int64 `env:"MY_INT_64" default:"1"`

	Uint   uint   `env:"MY_UINT" default:"1"`
	Uint8  uint8  `env:"MY_UINT_8" default:"1"`
	Uint16 uint16 `env:"MY_UINT_16" default:"1"`
	Uint32 uint32 `env:"MY_UINT_32" default:"1"`
	Uint64 uint64 `env:"MY_UINT_64" default:"1"`

	Float32 float32 `env:"MY_FLOAT32" default:"1"`
	Float64 float64 `env:"MY_FLOAT64" default:"1"`

	Duration            time.Duration `env:"MY_DURATION" default:"10"`
	StringifiedDuration time.Duration `env:"MY_STRINGIFIED_DURATION" default:"1s"`

	String string `env:"MY_STRING" default:"1"`

	Map               map[string]string `env:"MY_MAP" default:"1:2,3:4"`
	MapWithCustomSeps map[int]int       `env:"MY_CUSTOM_MAP" kv_separator:"-" separator:":" default:"1-2:3-4"`

	Slice              []string `env:"MY_SLICE" default:"1,2,3,4"`
	SliceWithCustomSep []int    `env:"MY_CUSTOM_SLICE" separator:"|" default:"1|2|3|4"`

	NestedConfig exampleNestedConfig
}

type exampleNestedConfig struct {
	A string `env:"MY_CONFIG.A" default:"nest_1"`
	B int    `env:"B"`
}

func TestLoad(t *testing.T) {
	var (
		logger = zerolog.New(os.Stderr)
	)
	testCases := map[string]struct {
		prep           func()
		input          interface{}
		expectedResult interface{}
		expectedError  environ.EnvError
	}{
		"not a pointer to a struct": {
			input: "this is a pointer to a struct",
			expectedError: environ.EnvError{
				Err:   environ.ErrInvalidInput,
				Key:   "config",
				Extra: "must be provided a pointer to a struct",
			},
			expectedResult: "this is a pointer to a struct",
		},
		"default values, empty env": {
			input: &exampleDefaultConfig{},
			expectedResult: &exampleDefaultConfig{
				Int:                 1,
				Int8:                1,
				Int16:               1,
				Int32:               1,
				Int64:               1,
				Uint:                1,
				Uint8:               1,
				Uint16:              1,
				Uint32:              1,
				Uint64:              1,
				Float32:             1,
				Float64:             1,
				Duration:            10,
				StringifiedDuration: time.Second,
				String:              "1",
				Map:                 map[string]string{"1": "2", "3": "4"},
				MapWithCustomSeps:   map[int]int{1: 2, 3: 4},
				Slice:               []string{"1", "2", "3", "4"},
				SliceWithCustomSep:  []int{1, 2, 3, 4},
				NestedConfig: exampleNestedConfig{
					A: "nest_1",
				},
			},
		},
		"with env values": {
			prep: func() {
				// setup env values
				os.Setenv("MY_INT", "0")
				os.Setenv("MY_INT_8", "2")
				os.Setenv("MY_INT_16", "4")
				os.Setenv("MY_INT_32", "6")
				os.Setenv("MY_INT_64", "8")
				os.Setenv("MY_UINT", "0")
				os.Setenv("MY_UINT_8", "2")
				os.Setenv("MY_UINT_16", "4")
				os.Setenv("MY_UINT_32", "6")
				os.Setenv("MY_UINT_64", "8")
				os.Setenv("MY_FLOAT32", "0.2")
				os.Setenv("MY_FLOAT64", "1.4")
				os.Setenv("MY_DURATION", "1000")
				os.Setenv("MY_STRINGIFIED_DURATION", "1m")
				os.Setenv("MY_STRING", "a longer string")
				os.Setenv("MY_MAP", "10:20,30:40")
				os.Setenv("MY_CUSTOM_MAP", "100-200:300-400")
				os.Setenv("MY_SLICE", "9,8,7,6")
				os.Setenv("MY_CUSTOM_SLICE", "5|4|3|2")
				os.Setenv("MY_CONFIG.A", "nested config value")
				os.Setenv("B", "4")
			},
			input: &exampleDefaultConfig{},
			expectedResult: &exampleDefaultConfig{
				Int:                 0,
				Int8:                2,
				Int16:               4,
				Int32:               6,
				Int64:               8,
				Uint:                0,
				Uint8:               2,
				Uint16:              4,
				Uint32:              6,
				Uint64:              8,
				Float32:             0.2,
				Float64:             1.4,
				Duration:            1000,
				StringifiedDuration: time.Minute,
				String:              "a longer string",
				Map:                 map[string]string{"10": "20", "30": "40"},
				MapWithCustomSeps:   map[int]int{100: 200, 300: 400},
				Slice:               []string{"9", "8", "7", "6"},
				SliceWithCustomSep:  []int{5, 4, 3, 2},
				NestedConfig: exampleNestedConfig{
					A: "nested config value",
					B: 4,
				},
			},
		},
		// TODO: add AWS Parameter Store
		// TODO: add AWS Secrets Manager
		// TODO: add GCP Secrets
		// TODO: add Swift Object Store
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.prep != nil {
				tc.prep()
			}
			err := environ.Load(tc.input)
			var envErr *environ.EnvError
			if errors.As(err, &envErr) {
				if tc.expectedError != *envErr {
					logger.Error().Err(err).Msg("expected error didn't match error")
					t.FailNow()
				}
			}

			if tc.input == tc.expectedResult {
				return
			}

			result, ok := tc.input.(*exampleDefaultConfig)
			if !ok {
				logger.Error().Msg("failed to convert input to config")
				t.FailNow()
			}

			expectedResult, ok := tc.expectedResult.(*exampleDefaultConfig)
			if !ok {
				logger.Error().Msg("failed to convert expected result to config")
				t.FailNow()
			}

			if !reflect.DeepEqual(result, expectedResult) {
				logger.Error().Any("result", result).Any("expected result", expectedResult).Msg("expected result does not match result")
				t.FailNow()
			}
		})
	}
}
