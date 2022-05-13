package report

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/mispon/hey_grpc/internal/work"
)

type report struct {
	request         string
	successRequests int
	failedRequests  int
	requestsTotal   int
	startTime       time.Time
	minDur          time.Duration
	maxDur          time.Duration
}

// Print prints results
func Print(args []string, startAt time.Time, results []work.Result) {
	rep := &report{
		request:   strings.Join(args, " "),
		startTime: startAt,
		minDur:    math.MaxInt,
	}

	for _, res := range results {
		rep.apply(res)
	}

	fmt.Println(rep)
}

func (r *report) apply(result work.Result) {
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

func (r *report) String() string {
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
