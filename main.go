package main

import (
	"github.com/blokje5/iam-server/pkg/log"
	"os"

	"github.com/blokje5/iam-server/cmd"
)

func main() {
	logger := log.GetLogger()
	if err := cmd.RootCommand.Execute(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
