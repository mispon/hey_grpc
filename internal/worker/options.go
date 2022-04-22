package worker

import (
	"time"

	"github.com/cheggaaa/pb/v3"
)

// CreateOption worker modificator
type CreateOption func(worker *Worker)

// WithCallsNumber specify number of unary calls
func WithCallsNumber(cn int) CreateOption {
	return func(w *Worker) {
		w.calls = cn
	}
}

// WithDelay specify delay between each call
func WithDelay(d time.Duration) CreateOption {
	return func(w *Worker) {
		w.delay = d
	}
}

// WithQPS specify number of calls per second
func WithQPS(q int) CreateOption {
	return func(w *Worker) {
		w.q = q
	}
}

// WithProgressBar specify progress bar
func WithProgressBar(pb *pb.ProgressBar) CreateOption {
	return func(w *Worker) {
		w.progress = pb
	}
}

// WithStopCh specify stop channel
func WithStopCh(sc <-chan struct{}) CreateOption {
	return func(w *Worker) {
		w.stopCh = sc
	}
}
