package load

import (
	"errors"
	"strings"
)

var (
	// ErrRequiredNotFound is the error for required variables that are not successfully loaded in the env
	ErrRequiredNotFound = errors.New("is required but failed to load value")
	// ErrLoading is the error for sources that are not loading properly
	ErrLoading = errors.New("encountered error loading value")
	// ErrInvalidFormat is the error for tags/params that are not in the correct format
	ErrInvalidFormat = errors.New("has invalid format")
	// ErrInvalidInput is the error for not providing a pointer to a struct
	ErrInvalidInput = errors.New("must be a pointer to a struct")
	// ErrUnsupportedType is the error for types that are not supported
	ErrUnsupportedType = errors.New("has unsupported type")
	// ErrUnsettableParam is the error for unsettable params, or unexported fields encountered in a struct
	ErrUnsettableParam = errors.New("must be a settable parameter")
)

// EnvError implements the error interface with key infomation and some helpful text for fixing the issues with loading a config
type EnvError struct {
	Err   error
	Key   string
	Extra string
}

// Error returns a user friendly error message in the format below
//
//	env: <key> <err message> | extra: <extra>
func (e *EnvError) Error() string {
	var sb strings.Builder
	sb.WriteString("env: ")
	sb.WriteString(e.Key)
	sb.WriteString(" ")
	sb.WriteString(e.Err.Error())
	if e.Extra != "" {
		sb.WriteString(" | extra: ")
		sb.WriteString(e.Extra)
	}
	return sb.String()
}

func newError(err error, key, extra string) *EnvError {
	return &EnvError{
		Err:   err,
		Key:   key,
		Extra: extra,
	}
}
