package request

import (
	"context"
	"os/exec"
	"time"
)

type Worker struct {
	// Calls is number of requests to make
	Calls int
	// QPS is query per second
	QPS int
	// Delay is timeout between each request
	Delay time.Duration
	// ResultCh is chan for request's results
	ResultCh chan<- Result
}

// Run starts worker job
// it's blocking call
func (w Worker) Run(ctx context.Context, args []string) {
	var throttle <-chan time.Time
	if w.QPS > 0 {
		throttle = time.Tick(time.Duration(1e6/(w.QPS)) * time.Microsecond)
	}

	for i := 0; i < w.Calls; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			if w.QPS > 0 {
				<-throttle
			}
			w.ResultCh <- unaryCall(args)
		}
	}
}

// unaryCall makes a grpc_cli call
func unaryCall(args []string) Result {
	args = append([]string{"call"}, args...)
	cmd := exec.Command("grpc_cli", args...)

	start := time.Now()
	err := cmd.Run()
	dur := time.Since(start)

	return Result{
		RequestDur: dur,
		Err:        err,
	}
}
