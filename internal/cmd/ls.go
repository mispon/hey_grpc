package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/mispon/hey_grpc/internal/reflection"
	"github.com/spf13/cobra"
)

var (
	lsCmd = &cobra.Command{
		Use:   "ls",
		Short: "List services",
		Example: `	hey_grpc ls host:port`,
		RunE: ls,
	}
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

func ls(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return NotEnoughArgsErr
	}

	ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
	defer cancel()

	refClient, err := reflection.NewClient(ctx, args[0])
	if err != nil {
		return err
	}

	services, err := refClient.ListServices()
	if err != nil {
		return err
	}

	for _, s := range services {
		fmt.Println(s)
	}
	return nil
}
