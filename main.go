package environ

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	// loading tags
	defaultTag  = "default"  // used to set a default value, any
	envTag      = "env"      // used to get value from env, string
	ssmTag      = "ssm"      // used to get value from AWS Parameter store, string
	asmTag      = "asm"      // used to get value from AWS Secrets Manager, string
	gsmTag      = "gsm"      // used to get value from GCP Secrets, string
	swiftTag    = "swift"    // used to get value from Swift based storage
	requiredTag = "required" // used to set requirements for env params, bool: causes errors when not loaded

	// formatting tags
	separatorTag   = "separator"    // used to select custom separators for slices and map items
	kvSeparatorTag = "kv_separator" // used to select custom separators for key value pairs in maps

	// defaults
	defaultSeparator   = ","
	defaultKvSeparator = ":"

	// misc helpers
	durationUnits = "smh"
)

// loads values based on tags provided on the struct
func Load(config any) error {
	configStruct, err := validateConfig(config)
	if err != nil {
		return err
	}
	err = handleStruct(configStruct)
	if err != nil {
		return err
	}
	return nil
}

// validates that a config is a pointer to a struct
func validateConfig(config any) (reflect.Value, error) {
	var output reflect.Value
	ptrRef := reflect.ValueOf(config)
	if ptrRef.Kind() != reflect.Ptr {
		return output, newError(ErrInvalidInput, "config", "must be provided a pointer to a struct")
	}
	output = ptrRef.Elem()
	if output.Kind() != reflect.Struct {
		return output, newError(ErrInvalidInput, "config", "must be provided a pointer to a struct")
	}
	return output, nil
}

// wraps handling fields of a struct
func handleStruct(input reflect.Value) error {
	var (
		inputType = input.Type()
		err       error
	)
	for i := 0; i < input.NumField(); i++ {
		var (
			field       = input.Field(i)
			structField = inputType.Field(i)
		)
		if !field.CanSet() {
			return newError(ErrUnsettableParam, structField.Name, "")
		}
		switch field.Kind() {
		case reflect.Struct:
			err = handleStruct(field)
		default:
			err = handleField(field, structField)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// wraps reading and setting a param value
func handleField(input reflect.Value, structField reflect.StructField) error {
	value, err := getValue(structField)
	if err != nil {
		return err
	}
	if value != "" {
		err = setValue(structField, input, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// reads value from env/stores based on field tags
func getValue(structField reflect.StructField) (string, error) {
	var (
		value    = structField.Tag.Get(defaultTag)
		required bool
		loaded   bool
		err      error
	)
	t, found := structField.Tag.Lookup(requiredTag)
	if found {
		required, err = strconv.ParseBool(t)
		if err != nil {
			return value, newError(ErrInvalidFormat, structField.Name, "required tag value is not a valid boolean representation")
		}
	}
	// check env
	t, found = structField.Tag.Lookup(envTag)
	if found {
		v := os.Getenv(t)
		if v != "" {
			loaded = true
			value = v
		}
	}
	// check other sources
	// TODO: get the parameter value from AWS
	// t, found = structField.Tag.Lookup(ssmTag)
	// if found {
	// }
	// TODO: get the secret value from AWS
	// t, found = structField.Tag.Lookup(asmTag)
	// if found {
	// }
	// TODO: get the secret value from GCP
	// t, found = structField.Tag.Lookup(gsmTag)
	// if found {
	// }
	// check if the field is required but not found/loaded
	if required && !loaded {
		return value, newError(ErrRequiredNotFound, structField.Name, "required field not loaded")
	}

	return value, nil
}

// set will set the loaded value to the param, or return an error
func setValue(structField reflect.StructField, param reflect.Value, value string) error {
	switch param.Type().Kind() {
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return newError(ErrInvalidFormat, structField.Name, "value is not a valid boolean representation")
		}
		param.SetBool(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var (
			v   int64
			err error
		)
		// handle parsing a stringified time.Duration env value
		if param.Type() == reflect.TypeOf(time.Duration(0)) && strings.ContainsAny(value, durationUnits) {
			var dur time.Duration
			dur, err = time.ParseDuration(value)
			v = dur.Nanoseconds()
		} else {
			v, err = strconv.ParseInt(value, 0, param.Type().Bits())
		}
		if err != nil {
			return newError(ErrInvalidFormat, structField.Name, "value is not a valid integer representation")
		}
		if v != 0 {
			param.SetInt(v)
		}
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(value, param.Type().Bits())
		if err != nil {
			return newError(ErrInvalidFormat, structField.Name, "value is not a valid float representation")
		}
		param.SetFloat(v)
	case reflect.Map:
		var (
			separator   = getSeparator(structField.Tag)
			values      = strings.Split(value, separator)
			kvSeparator = getKvSeparator(structField.Tag)
			t           = reflect.MakeMapWithSize(param.Type(), len(values))
		)
		for i := range values {
			var (
				kv    = strings.Split(values[i], kvSeparator)
				key   = reflect.New(param.Type().Key()).Elem()
				value = reflect.New(param.Type().Elem()).Elem()
			)
			if len(kv) != 2 {
				return newError(ErrInvalidFormat, structField.Name, "a map item has more than one kv_separator")
			}
			err := setValue(structField, key, kv[0])
			if err != nil {
				return err
			}
			err = setValue(structField, value, kv[1])
			if err != nil {
				return err
			}
			t.SetMapIndex(key, value)
		}
		param.Set(t)
	case reflect.Slice:
		values := strings.Split(value, getSeparator(structField.Tag))
		param.Grow(len(values))
		param.SetCap(len(values))
		param.SetLen(len(values))
		for i := range values {
			err := setValue(structField, param.Index(i), values[i])
			if err != nil {
				return err
			}
		}
	case reflect.String:
		param.SetString(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 0, param.Type().Bits())
		if err != nil {
			return newError(ErrInvalidFormat, structField.Name, "value is not a valid uint representation")
		}
		param.SetUint(val)
	default:
		return newError(ErrUnsupportedType, structField.Name, "provided type is not supported in this version")
	}
	return nil
}

func getSeparator(structTag reflect.StructTag) string {
	separator := defaultSeparator
	// get the separator from the tags
	if s, ok := structTag.Lookup(separatorTag); ok {
		separator = s
	}
	return separator
}

func getKvSeparator(structTag reflect.StructTag) string {
	separator := defaultKvSeparator
	// get the separator from the tags
	if s, ok := structTag.Lookup(kvSeparatorTag); ok {
		separator = s
	}
	return separator
}
