package environ_test

import (
	"log/slog"
	"testing"

	"github.com/NeedMoreVolume/environ"
)

func TestErrors(t *testing.T) {
	testCases := map[string]struct {
		envErr         environ.EnvError
		expectedOutput string
	}{
		"no extra": {
			envErr: environ.EnvError{
				Err: environ.ErrLoading,
				Key: "1",
			},
			expectedOutput: "env: 1 encountered error loading value",
		},
		"with extra": {
			envErr: environ.EnvError{
				Err:   environ.ErrInvalidInput,
				Key:   "2",
				Extra: "need valid input for parameters",
			},
			expectedOutput: "env: 2 must be a pointer to a struct | extra: need valid input for parameters",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			output := tc.envErr.Error()
			if output != tc.expectedOutput {
				slog.Error("output does not match expected output", "output", output, "expected output", tc.expectedOutput)
				t.Fail()
			}
		})
	}
}
