package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hey_grpc/internal/duration"
	"github.com/spf13/cobra"
)

var (
	callsNumber   int32
	concurrency   int32
	callsDuration string
	callTimeout   string
	jsonIn        bool
	jsonOut       bool

	callCmd = &cobra.Command{
		Use:     "call",
		Short:   "process grpc calls",
		Example: `  hey_grpc call -n 10 -c 1 host:port protoServer.Method 'i: 1, foo: "bar"'`,
		RunE:    execute,
	}
)

func init() {
	callCmd.PersistentFlags().Int32VarP(&callsNumber, "number", "n", 1, "-n 10")
	callCmd.PersistentFlags().Int32VarP(&concurrency, "concurrent", "c", 1, "-c 5")
	callCmd.PersistentFlags().StringVarP(&callsDuration, "during", "d", "", "-d 10s")
	callCmd.PersistentFlags().StringVarP(&callTimeout, "timeout", "t", "0s", "-t 1s")
	callCmd.PersistentFlags().BoolVar(&jsonIn, "json_in", false, "--json_in enables json_input for grpc_cli call")
	callCmd.PersistentFlags().BoolVar(&jsonOut, "json_out", false, "--json_out enables json_output for grpc_cli call")

	rootCmd.AddCommand(callCmd)
}

// execute runs calls
func execute(_ *cobra.Command, args []string) error {
	timeout, err := duration.ParseFlag(callTimeout)
	if err != nil {
		return err
	}

	if len(args) < 3 {
		return errors.New("unexpected arguments number\nsee https://github.com/grpc/grpc/blob/master/doc/command_line_tool.md")
	}

	if len(callsDuration) > 0 {
		return executeDuring(args, timeout)
	} else {
		return executeNumbers(args, timeout)
	}
}

// executeNumbers process specified call count
func executeNumbers(args []string, timeout time.Duration) error {
	var (
		wg    sync.WaitGroup
		calls int32
	)

	for i := 0; i < int(concurrency); i++ {
		wg.Add(1)
		go func(args []string, wg *sync.WaitGroup) {
			defer wg.Done()
			for calls < callsNumber {
				atomic.AddInt32(&calls, 1)
				makeCall(args)
				time.Sleep(timeout)
			}
		}(args, &wg)
	}

	wg.Wait()
	return nil
}

// executeDuring process calls during specified time
func executeDuring(args []string, timeout time.Duration) error {
	var (
		wg  sync.WaitGroup
		dur time.Duration
		err error
	)

	dur, err = duration.ParseFlag(callsDuration)
	if err != nil {
		return err
	}

	doneCh := make(chan struct{})
	go func(done chan<- struct{}) {
		time.Sleep(dur)
		close(done)
	}(doneCh)

	for i := 0; i < int(concurrency); i++ {
		wg.Add(1)
		go func(args []string, done <-chan struct{}, wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case <-time.After(timeout):
					makeCall(args)
				}
			}
		}(args, doneCh, &wg)
	}

	wg.Wait()
	return nil
}

// makeCall makes a grpc_cli call
func makeCall(args []string) {
	args = append([]string{"call"}, args...)

	if jsonIn {
		args = append(args, "--json_input")
	}
	if jsonOut {
		args = append(args, "--json_output")
	}

	cmd := exec.Command("grpc_cli", args...)
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}
