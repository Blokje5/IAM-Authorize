package main

import (
	"github.com/blokje5/iam-server/pkg/log"
	"os"

	"github.com/blokje5/iam-server/cmd"
)

func main() {
	if err := cmd.RootCommand.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
