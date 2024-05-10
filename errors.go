package environ

import (
	"errors"
	"strings"
)

var (
	ErrRequiredNotFound = errors.New("is required but not found")     // returned for required variables that are not found in the env
	ErrInvalidFormat    = errors.New("has invalid format")            // returned for tags/params that are not in the correct format
	ErrInvalidInput     = errors.New("must be a pointer to a struct") // returned for not providing a pointer to a struct
	ErrUnsupportedType  = errors.New("has unsupported type")          // returned for types that are not supported
)

type EnvError struct {
	error
	Key   string
	Extra string
}

// output: env: <key> <err message> | extra: <extra>
func (e *EnvError) Error() string {
	var sb strings.Builder
	sb.WriteString("env: ")
	sb.WriteString(e.Key)
	sb.WriteString(e.error.Error())
	sb.WriteString(" | extra: ")
	sb.WriteString(e.Extra)
	return sb.String()
}

func newError(err error, key, extra string) *EnvError {
	return &EnvError{
		error: err,
		Key:   key,
		Extra: extra,
	}
}
