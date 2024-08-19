package environ

import (
	"errors"
	"strings"
)

var (
	ErrRequiredNotFound = errors.New("is required but failed to load value") // returned for required variables that are not successfully loaded in the env
	ErrLoading          = errors.New("encountered error loading value")      // returned for sources that are not loading properly
	ErrInvalidFormat    = errors.New("has invalid format")                   // returned for tags/params that are not in the correct format
	ErrInvalidInput     = errors.New("must be a pointer to a struct")        // returned for not providing a pointer to a struct
	ErrUnsupportedType  = errors.New("has unsupported type")                 // returned for types that are not supported
	ErrUnsettableParam  = errors.New("must be a settable parameter")         // returned when unsettable params are encountered in a struct
)

type EnvError struct {
	Err   error
	Key   string
	Extra string
}

// output: env: <key> <err message> | extra: <extra>
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
