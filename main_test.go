package environ_test

import (
	"environ"
	"errors"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"
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

	Bool bool `env:"MY_BOOL" default:"false"`

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

func unsetTestEnv() {
	os.Unsetenv("MY_INT")
	os.Unsetenv("MY_INT_8")
	os.Unsetenv("MY_INT_16")
	os.Unsetenv("MY_INT_32")
	os.Unsetenv("MY_INT_64")
	os.Unsetenv("MY_UINT")
	os.Unsetenv("MY_UINT_8")
	os.Unsetenv("MY_UINT_16")
	os.Unsetenv("MY_UINT_32")
	os.Unsetenv("MY_UINT_64")
	os.Unsetenv("MY_FLOAT32")
	os.Unsetenv("MY_FLOAT64")
	os.Unsetenv("MY_DURATION")
	os.Unsetenv("MY_STRINGIFIED_DURATION")
	os.Unsetenv("MY_BOOL")
	os.Unsetenv("MY_STRING")
	os.Unsetenv("MY_MAP")
	os.Unsetenv("MY_CUSTOM_MAP")
	os.Unsetenv("MY_SLICE")
	os.Unsetenv("MY_CUSTOM_SLICE")
	os.Unsetenv("MY_CONFIG.A")
	os.Unsetenv("B")
}

func TestLoad(t *testing.T) {
	testCases := map[string]struct {
		prep           func()
		input          interface{}
		expectedResult interface{}
		expectedError  environ.EnvError
		clean          func()
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
			prep:  unsetTestEnv,
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
		"with good env values": {
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
				os.Setenv("MY_BOOL", "true")
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
				Bool:                true,
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
			clean: unsetTestEnv,
		},
		"with bad int value": {
			prep: func() {
				os.Setenv("MY_INT", "not an int")
			},
			input: &exampleDefaultConfig{},
			expectedError: environ.EnvError{
				Err:   environ.ErrInvalidFormat,
				Key:   "Int",
				Extra: "value is not a valid integer representation",
			},
			clean: func() {
				os.Unsetenv("MY_INT")
			},
		},
		"with bad float value": {
			prep: func() {
				os.Setenv("MY_FLOAT32", "not a float")
			},
			input: &exampleDefaultConfig{},
			expectedError: environ.EnvError{
				Err:   environ.ErrInvalidFormat,
				Key:   "Float32",
				Extra: "value is not a valid float representation",
			},
			clean: func() {
				os.Unsetenv("MY_FLOAT32")
			},
		},
		// TODO: add AWS Parameter Store
		// TODO: adsd AWS Secrets Manager
		// TODO: add GCP Secrets
		// TODO: add Swift Object Store
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// prep the test
			if tc.prep != nil {
				tc.prep()
			}
			// defer the clean up
			if tc.clean != nil {
				defer tc.clean()
			}
			// run the test
			err := environ.Load(tc.input)
			// validate error
			var envErr *environ.EnvError
			if errors.As(err, &envErr) {
				if tc.expectedError != *envErr {
					slog.Error("expected error didn't match error", "expected error", tc.expectedError, "error", envErr)
					t.FailNow()
				}
			}
			// done checking if this should have errored
			if tc.expectedError.Err != nil {
				return
			}
			// validate result
			if tc.input == tc.expectedResult {
				return
			}
			// otherwise use reflection to verify the results
			result, ok := tc.input.(*exampleDefaultConfig)
			if !ok {
				slog.Error("failed to convert input to config")
				t.FailNow()
			}
			expectedResult, ok := tc.expectedResult.(*exampleDefaultConfig)
			if !ok {
				slog.Error("failed to convert expected result to config")
				t.FailNow()
			}
			if !reflect.DeepEqual(result, expectedResult) {
				slog.Error("expected result does not match result", "expected result", expectedResult, "result", result)
				t.FailNow()
			}
		})
	}
}
