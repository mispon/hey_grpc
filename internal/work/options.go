package work

import (
	"time"
)

// WorkerOption worker modificator
type WorkerOption func(worker *Worker)

// WithTimeout specify timeout between each call
func WithTimeout(d time.Duration) WorkerOption {
	return func(w *Worker) {
		w.timeout = d
	}
}

// WithQPS specify number of callsNumbers per second
func WithQPS(q int) WorkerOption {
	return func(w *Worker) {
		w.qps = q
	}
}

// withStopCh specify stop channel
func withStopCh(sc <-chan struct{}) WorkerOption {
	return func(w *Worker) {
		w.stopCh = sc
	}
}
