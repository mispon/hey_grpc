package cmd

import (
	"errors"

	"github.com/mispon/hey_grpc/internal/report"

	"github.com/mispon/hey_grpc/internal/work"

	"github.com/mispon/hey_grpc/internal/flags"
	"github.com/spf13/cobra"
)

var (
	callsNumber   int32
	callsDuration string
	callTimeout   string
	workersNumber int32
	queryPerSec   int32

	callCmd = &cobra.Command{
		Use:     "call",
		Short:   "process grpc calls",
		Example: `  hey_grpc call -n 10 -c 1 host:port protoServer.Method 'i: 1, foo: "bar"'`,
		RunE:    execute,
	}
)

const (
	minCalls = 0
	maxCalls = 1_000_000

	minWorkers = 1
	maxWorkers = 500
)

func init() {
	callCmd.PersistentFlags().Int32VarP(&callsNumber, "number", "n", 0, "-n 10")
	callCmd.PersistentFlags().StringVarP(&callsDuration, "duration", "d", "0s", "-d 10s")
	callCmd.PersistentFlags().StringVarP(&callTimeout, "timeout", "t", "0s", "-t 1s")
	callCmd.PersistentFlags().Int32VarP(&workersNumber, "workers", "w", 1, "-w 5")
	callCmd.PersistentFlags().Int32VarP(&queryPerSec, "qps", "q", 0, "-q 100")

	rootCmd.AddCommand(callCmd)
}

// execute runs calls
func execute(cmd *cobra.Command, args []string) error {
	if len(args) < 3 {
		return errors.New("unexpected arguments number\nsee https://github.com/grpc/grpc/blob/master/doc/command_line_tool.md")
	}

	timeout, err := flags.ParseDuration(callTimeout)
	if err != nil {
		return err
	}

	duration, err := flags.ParseDuration(callsDuration)
	if err != nil {
		return err
	}

	if callsNumber == 0 && duration == 0 {
		return errors.New(`one of "number" or "duration" shouldn't be equal zero at the same time`)
	}

	// --duration has higher priority over --number
	if duration > 0 {
		callsNumber = 0
	}

	var (
		cn = flags.Clamp(int(callsNumber), minCalls, maxCalls)
		wn = flags.Clamp(int(workersNumber), minWorkers, maxWorkers)
	)

	w := work.New(
		cn,
		wn,
		duration,
		work.WithQPS(int(queryPerSec)),
		work.WithTimeout(timeout),
	)
	results := w.Execute(cmd.Context(), args)

	report.Print(args, results)

	return nil
}
