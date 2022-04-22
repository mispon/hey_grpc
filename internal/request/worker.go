package request

import (
	"context"
	"os/exec"
	"time"

	"github.com/cheggaaa/pb/v3"
)

type Worker struct {
	// Calls is required number of unary calls
	Calls int
	// QPS is query per second
	QPS int
	// Delay is timeout between each request
	Delay time.Duration
	// ResultCh is chan for request's results
	ResultCh chan<- Result
	// Progress is progress bar
	Progress *pb.ProgressBar
}

const (
	second = 1e9
)

// Run starts worker job
// it's blocking call
func (w *Worker) Run(ctx context.Context, args []string) {
	if w.Calls == 0 {
		return
	}

	var throttle <-chan time.Time
	if w.QPS > 0 {
		throttle = time.Tick(time.Duration(second / w.QPS))
	}

	for w.Calls > 0 {
		select {
		case <-ctx.Done():
			return
		default:
			if w.QPS > 0 {
				<-throttle
			}
			w.ResultCh <- unaryCall(args)
			w.Progress.Increment()
		}

		time.Sleep(w.Delay)
		w.Calls--
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
