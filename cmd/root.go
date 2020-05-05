package cmd

import (
	"github.com/spf13/cobra"
)

// RootCommand is the base CLI command that all subcommands are added to.
var RootCommand = &cobra.Command{
	Use:   "iam",
	Short: "IAM authorization server",
	Long:  "CLI to manage the IAM-Authorize server",
}