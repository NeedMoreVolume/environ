# Environ [![Go Coverage](https://github.com/NeedMoreVolume/environ/wiki/coverage.svg)](https://raw.githack.com/wiki/NeedMoreVolume/environ/coverage.html) [![Go Report Card](https://goreportcard.com/badge/github.com/NeedMoreVolume/environ)](https://goreportcard.com/report/github.com/NeedMoreVolume/environ) [![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT) [![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/NeedMoreVolume/environ)

## Description

There are two main components of environ. The library which facilitates config management in code, and the CLI tool to facilitate environment management.
The Library documentation can be found in the [Library Section](#library)
The CLI tool documentation can be found in the [CLI Section](#cli).

### Library
Here begins the documentation for how to use the library portion of this repo.

#### Tags

The library supports the following tags:
-  `env`: used to denote the key for loading an environment variable value 
- `default`: used to set any default value for an attribute
- `required`: used to flag that a value must be loaded and not empty (or return error if there is no value read from any source), supports truthy values.
- `separator`: used to override the default `,` separator for slice elements and map items.
- `kv_separator`: used to override the default `:` separator for key value pairs of map items.

#### Hows, whys, limitations

This library uses reflection to read attribute tags and set the values of the attributes of a provided struct accordingly.

This library treats unloaded required variables as an error. The reasoning behind this descision is that if truely required values are not loaded sucessfully it can lead to degraded service health or even total outage. This should help developers capture any configuration issues during the intialization phase, much like when using a Ping after opening a Mysql connection to validate the database is available and accessible. 

This library also provides a more detailed error structure, providing a Key and Extra with more information about the error but never any raw values to ensure no confidential data is accidentally leaked from logging loading errors.

Currently, the noteworthy limitations of this library are that config files are not supported, and maps of slices are not supported (IE: `map[string][]string`).

#### Usage

To use this library, just make a configuration struct with any of the above tags and then pass a pointer of one in the Load call.

##### Example

The following is an example of a simple Mysql config that contains the necessary information to open a DB connection.
```
import "github.com/NeedMoreVolume/environ"

type MysqlConfig struct {
	Username string `env:"MYSQL_USERNAME" required:"true"`
	Password string `env:"MYSQL_PASSWORD" required:"true"`
	Host     string `env:"MYSQL_HOST" required:"true"`
	Port     int    `env:"MYSQL_PORT" default:"3306"`
	Database string `env:"MYSQL_DATABASE" default:"example"`
}

func main() {
	var cfg MysqlConfig
	// load config
	err := environ.Load(&cfg)
	if err != nil {
   		// handle error
	}

    // can also take the errors.As approach
	var envErr *environ.EnvError
	if errors.As(err, &envErr) {
		// handle EnvError
	}
}
```
This config would fail to load if any of the username, password, or host values are not loaded successfully from a given environment.

#### Supported locations to load values from

Default values, supported by `default` tags

Environment variables, suppported by `env` tags

More to come!

### CLI

The environ CLI tool helps facilirate good developer experience when creating new configurations for applications. Currently, this tool is capable of taking a file containing configuration structs to create a series of env files with default values for each field of the structure.

#### Supported flags

The CLI tool supports the following flags:
- `input`: used to denote the location and filename to load the golang structs to create env files for/ Default: ./config/config.go
- `output`: used to denote the directory to save the env files. Default: .env/
- `target`: used to generate only a single env file for a given config name.

#### Hows, whys, limitations

Nobody writes OpenAPI specs by hand... right? Well, I sure don't want to write env files either. This cli tool helps to alleviate the toil in that process. Write your configuration structures, and let this do the dirty work of creating your env files with all the parameters!

#### Usage

environ -input=/path/to/file -output=/path/to/env-file-dir
