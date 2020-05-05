package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/blokje5/iam-server/pkg/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	runCommand.Flags().StringP("connectionstring", "c", "postgresql://admin:password@localhost:5432", "Pass the DB connection string to the server")
	viper.BindPFlag("connectionstring", runCommand.Flags().Lookup("connectionstring"))
	RootCommand.AddCommand(runCommand)
}

func run() {
	ctx := context.Background()


	params := server.NewParams()
	viper.Unmarshal(&params)
	server := server.New(params)
	if err := server.Init(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
