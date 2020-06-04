package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCommand is the base CLI command that all subcommands are added to.
var RootCommand = &cobra.Command{
	Use:   "iam",
	Short: "IAM authorization server",
	Long:  "CLI to manage the IAM-Authorize server",
}

func init() {
	viper.SetEnvPrefix("IAM")
	viper.AutomaticEnv()
}
