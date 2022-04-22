package report

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type Report struct {
	request         string
	successRequests int
	failedRequests  int
	requestsTotal   int
	startTime       time.Time
	minDur          time.Duration
	maxDur          time.Duration

	stopCh chan struct{}
}

// New creates new report instance
func New(args []string) *Report {
	return &Report{
		request: strings.Join(args, " "),
		minDur:  math.MaxInt,
		stopCh:  make(chan struct{}),
	}
}

// Watch process results and apply it to report
func (r *Report) Watch(resultCh <-chan Result) {
	r.startTime = time.Now()
	go func() {
		for {
			select {
			case result, ok := <-resultCh:
				if ok {
					r.apply(result)
				}
			case <-r.stopCh:
				return
			}
		}
	}()
}

func (r *Report) apply(result Result) {
	if r.minDur > result.RequestDur {
		r.minDur = result.RequestDur
	}
	if r.maxDur < result.RequestDur {
		r.maxDur = result.RequestDur
	}

	if result.Err == nil {
		r.successRequests++
	} else {
		r.failedRequests++
	}
	r.requestsTotal++
}

// Print finish watching goroutine and prints report
func (r Report) Print() {
	close(r.stopCh)
	fmt.Println(r)
}

func (r Report) String() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("Args: %s\n", r.request))

	sb.WriteString("Requests:\n")
	sb.WriteString(fmt.Sprintf("\tok:    %d\n", r.successRequests))
	sb.WriteString(fmt.Sprintf("\tfail:  %d\n", r.failedRequests))
	sb.WriteString(fmt.Sprintf("\ttotal: %d\n", r.requestsTotal))

	sb.WriteString("Durations:\n")
	sb.WriteString(fmt.Sprintf("\tmin:   %v\n", r.minDur))
	sb.WriteString(fmt.Sprintf("\tmax:   %v\n", r.maxDur))
	sb.WriteString(fmt.Sprintf("\ttotal: %v\n", time.Since(r.startTime)))

	return sb.String()
}
