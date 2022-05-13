package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var (
	NotEnoughArgsErr = errors.New("insufficient number of arguments")
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
