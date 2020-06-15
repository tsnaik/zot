package main

import (
	"os"

	"github.com/anuvu/zot/pkg/cli"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configPath := home + "/.zot"

	if err := cli.NewRootCmd(configPath).Execute(); err != nil {
		os.Exit(1)
	}
}
