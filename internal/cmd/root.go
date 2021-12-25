package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "hey_grpc",
		Short: "very simple grpc_cli wrapper",
	}
)

func Execute() error {
	return rootCmd.Execute()
}
