package flags

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Batches(t *testing.T) {
	testCases := []struct {
		n   int
		c   int
		exp []int
	}{
		{n: 10, c: 3, exp: []int{3, 3, 4}},
		{n: 10, c: 2, exp: []int{5, 5}},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			b := Batches(tc.n, tc.c)
			require.Equal(t, tc.exp, b)
		})
	}
}
