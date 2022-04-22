package commands

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/mispon/hey_grpc/internal/flags"
	"github.com/mispon/hey_grpc/internal/request"
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
	maxCalls   = 1_000_000
	maxWorkers = 500
)

func init() {
	callCmd.PersistentFlags().Int32VarP(&callsNumber, "number", "n", 0, "-n 10")
	// todo
	// callCmd.PersistentFlags().StringVarP(&callsDuration, "during", "d", "0s", "-d 10s")
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

	// todo
	/*duration, err := flags.ParseDuration(callsDuration)
	if err != nil {
		return err
	}*/

	ctx, cancelFn := context.WithCancel(cmd.Context())
	defer cancelFn()

	var (
		cn  = min(int(callsNumber), maxCalls)
		wn  = min(int(workersNumber), maxWorkers)
		qps = int(queryPerSec)

		resultCh = make(chan request.Result, cn*wn)
		progress = pb.StartNew(cn * wn)
	)

	startTime := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(wn)

	for i := 0; i < wn; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			w := request.Worker{
				Calls:    cn,
				QPS:      qps,
				Delay:    timeout,
				ResultCh: resultCh,
				Progress: progress,
			}
			w.Run(ctx, args)
		}(&wg)
	}

	wg.Wait()
	progress.Finish()
	close(resultCh)

	request.PrintReport(resultCh, args, time.Since(startTime))
	return nil
}

func min(value, limit int) int {
	if value > limit {
		return limit
	}
	return value
}
