package work

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrepareUnaryCall(t *testing.T) {
	t.SkipNow()

	args := []string{
		"localhost:7002",       // host:port
		"Strings/ToUpper",      // <service>/<method>
		"str: \"hello world\"", // payload
	}

	w := Worker{verbose: true}
	uc, err := w.prepareUnaryCall(context.TODO(), args)
	require.NoError(t, err)

	res := uc()
	require.NoError(t, res.Err)
}
