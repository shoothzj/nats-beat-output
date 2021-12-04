package main

import (
	"github.com/elastic/beats/v7/heartbeat/cmd"
	"os"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
