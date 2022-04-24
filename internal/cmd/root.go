package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:                    "hey_grpc",
		Short:                  "grpc_cli wrapper for simple load tests",
		BashCompletionFunction: "todo",
	}
)

func Execute() error {
	return rootCmd.Execute()
}
