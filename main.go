package main

import (
	"flag"
	"os"
	"strings"

	"github.com/NeedMoreVolume/environ/cli"
)

func main() {
	switch os.Args[1] {
	case cli.GenerateEnvFilesCommand:
		// provide the flag arguments to be injected into the cli
		envFileGenerator, err := cli.NewEnvFileGenerator(os.Args[2:])
		if err != nil {
			os.Exit(1)
		}
		if envFileGenerator == nil {
			os.Exit(1)
		}
		_, err = envFileGenerator.CreateEnvFiles()
		if err != nil {
			// determine if writing error to exit(2)
			if strings.Contains(err.Error(), "failed to write") {
				os.Exit(2)
			}
			os.Exit(1)
		}
	default:
		flag.Usage()
	}
}
