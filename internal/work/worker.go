package work

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jhump/protoreflect/dynamic"

	"github.com/jhump/protoreflect/dynamic/grpcdynamic"

	"github.com/mispon/hey_grpc/internal/reflection"
)

type Worker struct {
	index   int
	calls   int
	qps     int
	timeout time.Duration
	results chan<- Result
	stopCh  <-chan struct{}
}

const (
	second = 1e9
)

// newWorker creates new worker instance
func newWorker(callsCount int, resultCh chan<- Result, opts ...WorkerOption) *Worker {
	w := &Worker{
		calls:   callsCount,
		results: resultCh,
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

// run starts worker job
// it's blocking call
func (w *Worker) run(ctx context.Context, args []string) {
	unaryCall, err := w.prepareUnaryCall(ctx, args)
	if err != nil {
		w.results <- Result{Err: err}
		return
	}

	var throttle <-chan time.Time
	if w.qps > 0 {
		throttle = time.Tick(time.Duration(second / w.qps))
	}

	finished := w.finalizer()
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopCh:
			return
		default:
			if w.qps > 0 {
				<-throttle
			}
			w.results <- unaryCall()
		}

		if finished() {
			return
		}

		time.Sleep(w.timeout)
	}
}

// prepareUnaryCall is unary call builder
func (w *Worker) prepareUnaryCall(ctx context.Context, args []string) (func() Result, error) {
	// create reflection client
	refClient, conn, err := reflection.NewClientConn(ctx, args[0])
	if err != nil {
		return nil, err
	}

	// parse <service>/<method> string
	service, method, err := parseService(args[1])
	if err != nil {
		return nil, err
	}

	// resolve service by name
	svcDesc, err := refClient.ResolveService(service)
	if err != nil {
		return nil, err
	}

	// resolve method by name
	methodDesc := svcDesc.FindMethodByName(method)
	if methodDesc == nil {
		return nil, errors.New("method not exists")
	}

	// resolve method's input message
	in := methodDesc.GetInputType().GetFullyQualifiedName()
	messageDesc, err := refClient.ResolveMessage(in)
	if err != nil {
		return nil, err
	}

	// enrich message the data
	message := dynamic.NewMessage(messageDesc)
	err = message.UnmarshalText([]byte(args[2]))
	if err != nil {
		return nil, err
	}

	// create dynamic conn
	dynConn := grpcdynamic.NewStub(conn)

	return func() Result {
		start := time.Now()
		_, rpcErr := dynConn.InvokeRpc(context.TODO(), methodDesc, message)
		rpcDur := time.Since(start)

		return Result{
			RequestDur: rpcDur,
			Err:        rpcErr,
		}
	}, nil
}

func parseService(service string) (string, string, error) {
	p := strings.Split(service, "/")
	if len(p) != 2 {
		return "", "", errors.New("failed to parse <service>/<method> string")
	}

	return p[0], p[1], nil
}

// finalizer creates stop check func
func (w *Worker) finalizer() func() bool {
	// create calls count checker
	if w.calls > 0 {
		total := 0
		return func() bool {
			total++
			return w.calls == total
		}
	}

	// create stop signal checker
	return func() bool {
		select {
		case <-w.stopCh:
			return true
		default:
			return false
		}
	}
}
