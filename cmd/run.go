package cmd

import (
	"github.com/blokje5/iam-server/pkg/server"
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	runCommand := &cobra.Command{
		Use:   "run",
		Short: "Run the IAM-Authorize Server",
		Long:  "Run the IAM-Authorize Server",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}

	RootCommand.AddCommand(runCommand)
}

func run() {
	ctx := context.Background()

	server := server.New()
	if err := server.Init(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
