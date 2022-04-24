package work

import (
	"context"
	"os/exec"
	"time"
)

type Worker struct {
	index   int
	calls   int
	qps     int
	timeout time.Duration
	results chan<- Result
	stopCh  <-chan struct{}
}

const (
	second = 1e9
)

// newWorker creates new worker instance
func newWorker(callsCount int, resultCh chan<- Result, opts ...WorkerOption) *Worker {
	w := &Worker{
		calls:   callsCount,
		results: resultCh,
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

// run starts worker job
// it's blocking call
func (w *Worker) run(ctx context.Context, args []string) {
	var throttle <-chan time.Time
	if w.qps > 0 {
		throttle = time.Tick(time.Duration(second / w.qps))
	}

	finished := w.finalizer()
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopCh:
			return
		default:
			if w.qps > 0 {
				<-throttle
			}
			w.results <- unaryCall(args)
		}

		if finished() {
			return
		}

		time.Sleep(w.timeout)
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

// finalizer creates stop check func
func (w *Worker) finalizer() func() bool {
	// create calls count checker
	if w.calls > 0 {
		total := 0
		return func() bool {
			total++
			return w.calls == total
		}
	}

	// create stop signal checker
	return func() bool {
		select {
		case <-w.stopCh:
			return true
		default:
			return false
		}
	}
}
