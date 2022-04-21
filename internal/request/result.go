package request

import "time"

type Result struct {
	RequestDur time.Duration
	Err        error
}
