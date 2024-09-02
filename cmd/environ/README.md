[![Go Coverage](https://github.com/NeedMoreVolume/environ/wiki/coverage.svg)](https://raw.githack.com/wiki/NeedMoreVolume/environ/coverage.html)

# Environ CLI Tool

The environ CLI tool helps facilirate good developer experience when creating new configurations for applications. Currently, this tool is capable of taking a file containing configuration structs to create a series of env files with default values for each field of the structure.

## Supported Flags

The CLI tool supports the following flags:
- `input`: used to denote the location and filename to load the golang structs to create env files for/ Default: ./config/config.go
- `output`: used to denote the directory to save the env files. Default: .env/
- `target`: used to generate only a single env file for a given config name.

## Hows, whys, limitations

Nobody writes OpenAPI specs by hand... right? Well, I sure don't want to create env files either. This cli tool helps to alleviate the toil in that process. Write your configuration structures, and let this do the dirty work of creating your env files with all the parameters!

## Usage

environ -input <file> -output <file>

### Example

The following is an example of a simple Mysql config that contains the necessary information to open a DB connection.
```
import "github.com/NeedMoreVolume/environ/load"

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

## Supported locations to load values from

Default values, supported by `default` tags

Environment variables, suppported by `env` tags
