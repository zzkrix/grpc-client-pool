package demo

import (
	"context"

	"google.golang.org/grpc"
)

func (f *ClientFactory) processGreeter(fn func(GreeterClient)) (err error) {
	conn, err := f.pool.Get()
	if err != nil {
		return err
	}

	defer func() {
		err = conn.Close()
	}()

	var impl = NewGreeterClient(conn.Value())

	fn(impl)
	return nil
}

func (f *ClientFactory) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (resp *HelloResponse, err error) {
	err = f.processGreeter(func(impl GreeterClient) {
		resp, err = impl.SayHello(ctx, in, opts...)
	})
	return
}
