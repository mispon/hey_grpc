package worker

import (
	"context"
	"os/exec"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/mispon/hey_grpc/internal/report"
)

type Worker struct {
	// required number of unary calls
	calls int
	// query per sec
	q int
	// timeout between each call
	delay time.Duration
	// chan for request's results
	results chan<- report.Result
	// progress bar
	progress *pb.ProgressBar
	// stop channel
	stopCh <-chan struct{}
}

const (
	second = 1e9
)

// New creates new worker instance
func New(resultCh chan<- report.Result, opts ...CreateOption) *Worker {
	w := &Worker{
		results: resultCh,
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

// Run starts worker job
// it's blocking call
func (w *Worker) Run(ctx context.Context, args []string) {
	var throttle <-chan time.Time
	if w.q > 0 {
		throttle = time.Tick(time.Duration(second / w.q))
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopCh:
			return
		default:
			if w.q > 0 {
				<-throttle
			}
			w.results <- unaryCall(args)
		}

		if w.mustFinish() {
			return
		}

		time.Sleep(w.delay)
	}
}

// unaryCall makes a grpc_cli call
func unaryCall(args []string) report.Result {
	args = append([]string{"call"}, args...)
	cmd := exec.Command("grpc_cli", args...)

	start := time.Now()
	err := cmd.Run()
	dur := time.Since(start)

	return report.Result{
		RequestDur: dur,
		Err:        err,
	}
}

// mustFinish checks should worker stops
func (w *Worker) mustFinish() bool {
	if w.calls > 0 {
		w.progress.Increment()
		w.calls--
		return w.calls == 0
	}

	select {
	case <-w.stopCh:
		return true
	default:
		return false
	}
}
