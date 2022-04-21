package flags

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ParseFlag(t *testing.T) {
	type testCase struct {
		flag     string
		result   time.Duration
		hasError bool
	}

	testCases := []testCase{
		{
			flag:     "1h",
			result:   time.Hour,
			hasError: false,
		},
		{
			flag:     "10m",
			result:   time.Minute * 10,
			hasError: false,
		},
		{
			flag:     "200s",
			result:   time.Second * 200,
			hasError: false,
		},
		{
			flag:     "600ms",
			result:   time.Millisecond * 600,
			hasError: false,
		},
		{
			flag:     "foo",
			result:   0,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("check %s parsed to %d", tc.flag, tc.result), func(t *testing.T) {
			result, err := ParseDuration(tc.flag)
			assert.Equal(t, result, tc.result)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
