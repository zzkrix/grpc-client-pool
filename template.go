package main

var packageTpl = `
package %s

import (
	"context"

	"google.golang.org/grpc"
)
`

var globalTpl = `
package %s

import (
	"errors"
	"sync"

	"github.com/shimingyah/pool"
)

var clientMap = sync.Map{}

type ClientOptions struct {
	ClientID string
	Addr string
}

// AddClient 不能重复注册
func AddClient(opt *ClientOptions) error {
	c, err := newClientFactory(opt.Addr)
	if err != nil {
		return err
	}

	// 保存到全局map
	_, loaded := clientMap.LoadOrStore(opt.ClientID, c)
	if loaded {
		return errors.New("client exist, refused register")
	}

	return nil
}

func GetClient(clientID string) (*ClientFactory, error) {
	v, ok := clientMap.Load(clientID)
	if !ok {
		return nil, errors.New("client not found")
	}

	client, ok := v.(*ClientFactory)
	if !ok {
		return nil, errors.New("client type invalid")
	}

	return client, nil
}

type ClientFactory struct {
	pool pool.Pool
}

func newClientFactory(addr string) (*ClientFactory, error) {
	p, err := pool.New(addr, pool.DefaultOptions)
	if err != nil {
		return nil, err
	}
	
	return &ClientFactory{pool: p}, nil
}
`

var svcProcessTpl = `
func (f *ClientFactory) process%s(fn func(%sClient)) (err error) {
	conn, err := f.pool.Get()
	if err != nil {
		return err
	}
	
	defer func() {
		err = conn.Close()
	}()
	
	var impl = New%sClient(conn.Value())
	
	fn(impl)
	return nil
}
`

var svcMethodTpl = `
func (f *ClientFactory) %s(ctx context.Context, in *%s, opts ...grpc.CallOption) (resp *%s, err error) {
	err = f.process%s(func(impl %sClient) {
		resp, err = impl.%s(ctx, in, opts...)
	})
	return
}`
