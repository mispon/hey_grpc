package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	version = "0.0.1"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "prints hey_grpc version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("hey_grpc v%s", version)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
