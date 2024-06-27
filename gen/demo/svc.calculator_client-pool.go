package demo

import (
	"context"

	"google.golang.org/grpc"
)

func (f *ClientFactory) processCalculator(fn func(CalculatorClient)) (err error) {
	conn, err := f.pool.Get()
	if err != nil {
		return err
	}

	defer func() {
		err = conn.Close()
	}()

	var impl = NewCalculatorClient(conn.Value())

	fn(impl)
	return nil
}

func (f *ClientFactory) Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (resp *AddResponse, err error) {
	err = f.processCalculator(func(impl CalculatorClient) {
		resp, err = impl.Add(ctx, in, opts...)
	})
	return
}
func (f *ClientFactory) Sub(ctx context.Context, in *SubRequest, opts ...grpc.CallOption) (resp *SubResponse, err error) {
	err = f.processCalculator(func(impl CalculatorClient) {
		resp, err = impl.Sub(ctx, in, opts...)
	})
	return
}
