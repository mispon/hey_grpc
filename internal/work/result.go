package work

import "time"

// Result is unary call data
type Result struct {
	RequestDur time.Duration
	Err        error
}
