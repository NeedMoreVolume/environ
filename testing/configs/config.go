package testdata

// Config is an example of a config struct
type Config struct {
	LogLevel string `env:"LOG_LEVEL" default:"info"`
}
