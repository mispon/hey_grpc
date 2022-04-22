package request

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// PrintReport prints requests statistic
func PrintReport(resultCh <-chan Result, args []string, totalDur time.Duration) {
	rep := report{
		request:       strings.Join(args, " "),
		requestsTotal: len(resultCh),
		minDur:        math.MaxInt,
		totalDur:      totalDur,
	}

	for res := range resultCh {
		rep.apply(res)
	}

	fmt.Println(rep)
}

type report struct {
	request         string
	successRequests int
	failedRequests  int
	requestsTotal   int
	totalDur        time.Duration
	minDur          time.Duration
	maxDur          time.Duration
}

func (r *report) apply(result Result) {
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
}

func (r report) String() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("Args: %s\n", r.request))

	sb.WriteString("Requests:\n")
	sb.WriteString(fmt.Sprintf("\tok:    %d\n", r.successRequests))
	sb.WriteString(fmt.Sprintf("\tfail:  %d\n", r.failedRequests))
	sb.WriteString(fmt.Sprintf("\ttotal: %d\n", r.requestsTotal))

	sb.WriteString("Durations:\n")
	sb.WriteString(fmt.Sprintf("\tmin:   %v\n", r.minDur))
	sb.WriteString(fmt.Sprintf("\tmax:   %v\n", r.maxDur))
	sb.WriteString(fmt.Sprintf("\ttotal: %v\n", r.totalDur))

	return sb.String()
}
