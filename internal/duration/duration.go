package duration

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type (
	durationType struct {
		suffix string
		factor time.Duration
	}
)

func ParseFlag(flag string) (time.Duration, error) {
	types := []durationType{
		{"ms", time.Millisecond},
		{"s", time.Second},
		{"m", time.Minute},
		{"h", time.Hour},
	}

	for _, dt := range types {
		if strings.HasSuffix(flag, dt.suffix) {
			value := strings.TrimRight(flag, dt.suffix)
			return createNew(value, dt.factor)
		}
	}

	return 0, errors.New("unexpected duration type")
}

func createNew(value string, factor time.Duration) (time.Duration, error) {
	dur, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return time.Duration(dur) * factor, nil
}
