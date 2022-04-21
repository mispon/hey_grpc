package request

import (
	"fmt"
	"strings"
	"time"
)

// PrintReport prints requests statistic
func PrintReport(resultCh <-chan Result, totalDur time.Duration) {
	rep := report{
		requestCount: len(resultCh),
		totalDur:     totalDur,
	}

	for res := range resultCh {
		rep.apply(res)
	}

	fmt.Println(rep)
}

type report struct {
	requestCount int
	totalDur     time.Duration
	minDur       time.Duration
	maxDur       time.Duration
}

func (r *report) apply(result Result) {
	if r.minDur > result.RequestDur {
		r.minDur = result.RequestDur
	}
	if r.maxDur < result.RequestDur {
		r.maxDur = result.RequestDur
	}
}

func (r report) String() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("Requests: %v\n", r.requestCount))
	sb.WriteString(fmt.Sprintf("Duration:"))
	sb.WriteString(fmt.Sprintf("\tmin: %v\n", r.minDur))
	sb.WriteString(fmt.Sprintf("\tmax: %v\n", r.maxDur))
	sb.WriteString(fmt.Sprintf("\ttotal: %v\n", r.totalDur))

	return sb.String()
}
