package main

import (
	"fmt"
	"os"

	"github.com/blokje5/iam-server/cmd"
)

func main() {
	if err := cmd.RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}