package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	version = "0.2.0"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints current version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("hey_grpc v%s", version)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
