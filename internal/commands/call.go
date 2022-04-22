package commands

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/mispon/hey_grpc/internal/flags"
	"github.com/mispon/hey_grpc/internal/report"
	"github.com/mispon/hey_grpc/internal/worker"
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

	resultBuf = 100
)

type (
	workerFactory func(i int) *worker.Worker
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

	ctx, cancelFn := context.WithCancel(cmd.Context())
	defer cancelFn()

	resultCh := make(chan report.Result, resultBuf)

	rt := report.New(args)
	rt.Watch(resultCh)

	var (
		createWorker workerFactory
		progress     *pb.ProgressBar

		workersNum = flags.Clamp(int(workersNumber), minWorkers, maxWorkers)
		baseOpts   = []worker.CreateOption{
			worker.WithQPS(int(queryPerSec)),
			worker.WithDelay(timeout),
		}
	)

	if duration > 0 {
		stopCh := make(chan struct{})
		progress = pb.StartNew(int(duration.Seconds()))

		runTimeline(duration, stopCh, progress)
		createWorker = timelineWorkerFactory(resultCh, stopCh, baseOpts...)
	} else {
		cn := flags.Clamp(int(callsNumber), minCalls, maxCalls)
		progress = pb.StartNew(cn)

		callsBatches := flags.Batches(cn, workersNum)
		createWorker = quantityWorkerFactory(resultCh, callsBatches, progress, baseOpts...)
	}

	wg := sync.WaitGroup{}
	wg.Add(workersNum)
	for i := 0; i < workersNum; i++ {
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()
			w := createWorker(i)
			w.Run(ctx, args)
		}(&wg, i)
	}
	wg.Wait()

	close(resultCh)
	progress.Finish()
	rt.Print()

	return nil
}

func timelineWorkerFactory(
	resultCh chan<- report.Result,
	stopCh <-chan struct{},
	baseOpts ...worker.CreateOption,
) workerFactory {
	return func(i int) *worker.Worker {
		baseOpts = append(
			baseOpts,
			worker.WithStopCh(stopCh),
		)
		return worker.New(resultCh, baseOpts...)
	}
}

func quantityWorkerFactory(
	resultCh chan<- report.Result,
	callsBatches []int,
	pb *pb.ProgressBar,
	baseOpts ...worker.CreateOption,
) workerFactory {
	return func(i int) *worker.Worker {
		baseOpts = append(
			baseOpts,
			worker.WithCallsNumber(callsBatches[i]),
			worker.WithProgressBar(pb),
		)
		return worker.New(resultCh, baseOpts...)
	}
}

func runTimeline(dur time.Duration, stopCh chan struct{}, pb *pb.ProgressBar) {
	go func() {
		for {
			select {
			case <-stopCh:
				return
			default:
				pb.Increment()
				time.Sleep(time.Second)
			}
		}
	}()

	go func() {
		time.Sleep(dur)
		close(stopCh)
	}()
}
