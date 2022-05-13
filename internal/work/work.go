package work

import (
	"context"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/mispon/hey_grpc/internal/flags"
)

const (
	buffer = 100
)

// Work represent "call" cmd progress
type Work struct {
	pb *pb.ProgressBar

	callsNum   int
	workersNum int
	duration   time.Duration
	workerOpts []WorkerOption
}

// New creates new Work instance
func New(cn, wn int, dur time.Duration, wOpts ...WorkerOption) *Work {
	p := cn
	if dur > 0 {
		p = int(dur.Seconds())
	}

	return &Work{
		pb:         pb.StartNew(p),
		callsNum:   cn,
		workersNum: wn,
		duration:   dur,
		workerOpts: wOpts,
	}
}

// Execute starts workers and keep results
func (w *Work) Execute(ctx context.Context, args []string) []Result {
	var (
		result       = make([]Result, 0, buffer)
		resultCh     = make(chan Result)
		doneCh       = make(chan struct{})
		callsBatches = make([]int, w.workersNum)
	)

	ctx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	if w.duration > 0 {
		stopCh := make(chan struct{})
		w.workerOpts = append(w.workerOpts, withStopCh(stopCh))
		go w.keepAlive(stopCh)
	} else {
		callsBatches = flags.Batches(w.callsNum, w.workersNum)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(doneCh)
				return
			case r, ok := <-resultCh:
				if w.callsNum > 0 {
					w.pb.Increment()
				}
				if !ok {
					close(doneCh)
					return
				}
				result = append(result, r)
			}
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(w.workersNum)

	for i := 0; i < w.workersNum; i++ {
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()
			wk := newWorker(callsBatches[i], resultCh, w.workerOpts...)
			wk.run(ctx, args)
		}(&wg, i)
	}

	wg.Wait()

	time.Sleep(time.Second)
	close(resultCh)

	<-doneCh
	return result
}

// keepAlive keeps the work active for a duration time
func (w *Work) keepAlive(stopCh chan struct{}) {
	go func() {
		for {
			select {
			case <-stopCh:
				return
			default:
				w.pb.Increment()
				time.Sleep(time.Second)
			}
		}
	}()

	go func() {
		time.Sleep(w.duration)
		close(stopCh)
	}()
}
