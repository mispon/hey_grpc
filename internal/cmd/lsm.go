package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/mispon/hey_grpc/internal/reflection"
	"github.com/spf13/cobra"
)

var (
	lsmCmd = &cobra.Command{
		Use:   "lsm",
		Short: "List service methods",
		Example: `	hey_grpc lsm host:port <service>`,
		RunE: lsm,
	}
)

func init() {
	rootCmd.AddCommand(lsmCmd)
}

func lsm(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return NotEnoughArgsErr
	}

	ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
	defer cancel()

	refClient, err := reflection.NewClient(ctx, args[0])
	if err != nil {
		return err
	}

	svcDesc, err := refClient.ResolveService(args[1])
	if err != nil {
		return err
	}

	for _, md := range svcDesc.GetMethods() {
		fmt.Println(md.GetName())
	}
	return nil
}
