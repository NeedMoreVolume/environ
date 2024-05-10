# Environ

Golang library for loading config values.

This library uses reflection to read attribute tags and set the values of the attributes of a provided struct accordingly.

The library supports the following tags:
-  `env`: used to denote the key for loading an environment variable value 
- `default`: used to set any default value for an attribute
- `required`: used to flag that a value must be readable and not empty (or return error if there is no value read from any source), supports truthy values.
- `separator`: used to override the default `,` separator for slice elements and map items.
- `kv_separator`: used to override the default `:` separator for key value pairs of map items.

To use this library, just make a configuration struct with any of the above tags and then pass a pointer of one in the Load call.
This library also provides a more detailed error structure, providing a Key and Extra with more information about the error but never any raw values to ensure no confidential data is accidentally leaked from logging loading errors.

## Example
```
package main

import github.com/NeedMoreVolume/environ

type Api struct {
  // http timeouts
  IdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT" default"15s"`
  ReadTimeout time.Duration `env:"HTTP_READ_TIMEOUT" default:"15s"`
  WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" default:"60s"`

  // db conn
  DbHost string `env:"DB_HOST" default:"localhost"`
  DbPort int    `env:"DB_PORT" default:"3306"`
  DbName string `env:"DB_NAME" default:"example"`
  DbUser string `env:"DB_USER" default:"root"`
  DbPass string `env:"DB_PASS" required:"true"`
}


func main() {
  var config Api
  err := environ.Load(&config)
  if err != nil {
    var envErr *environ.EnvError
    if errors.As(err, &envErr) {
      ...
    }
  }
  ...
}
```

## Supported value locations
Environment variables, suppported by `env` tags

## Future plans
1. Integrate 3rd party key value stores such as Google Secrets Manager, AWS Secrets Manager, and AWS Parameter Store.
2. Support loading with options in a new func LoadWithOpts(config any, EnvironOpts...) in order to pass options for loading purposes as one might when using the above value stores.
