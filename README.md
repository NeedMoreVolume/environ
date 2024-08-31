[![Go Coverage](https://github.com/NeedMoreVolume/environ/wiki/coverage.svg)](https://raw.githack.com/wiki/NeedMoreVolume/environ/coverage.html)

# Environ

Environ facilitates loading values for a struct by tags, where the tags define the rules for a value and where it can be loaded from.

## Tags

The library supports the following tags:
-  `env`: used to denote the key for loading an environment variable value 
- `default`: used to set any default value for an attribute
- `required`: used to flag that a value must be loaded and not empty (or return error if there is no value read from any source), supports truthy values.
- `separator`: used to override the default `,` separator for slice elements and map items.
- `kv_separator`: used to override the default `:` separator for key value pairs of map items.

## Hows, whys, limitations
This library uses reflection to read attribute tags and set the values of the attributes of a provided struct accordingly.
This library treats unloaded required variables as an error. The reasoning behind this descision is that if truely required values are not loaded sucessfully it can lead to degraded service health or even total outage. This should help developers capture any configuration issues during the intialization phase, much like when using a Ping after opening a Mysql connection to validate the database is available and accessible. 
This library also provides a more detailed error structure, providing a Key and Extra with more information about the error but never any raw values to ensure no confidential data is accidentally leaked from logging loading errors.
Currently, the noteworthy limitations of this library are that config files are not supported, and maps of slices are not supported (IE: `map[string][]string`).

## Usage

To use this library, just make a configuration struct with any of the above tags and then pass a pointer of one in the Load call.

### Example

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

    // can also take the errors.As approach
		var envErr *environ.EnvError
		if errors.As(err, &envErr) {
      // handle EnvError
    }
	}
}
```
This config would fail to load if any of the username, password, or host values are not loaded successfully from a given environment.

## Supported locations to load values from

Default values, supported by `default` tags
Environment variables, suppported by `env` tags
