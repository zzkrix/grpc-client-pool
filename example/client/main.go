package main

import (
	"flag"
	"log"

	"protoc-gen-client-pool/example/utils/trace"
	pb "protoc-gen-client-pool/gen/demo"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()

	clientID := "demo-1"

	if err := pb.AddClient(&pb.ClientOptions{
		ClientID: clientID,
		Addr:     *addr,
	}); err != nil {
		panic(err)
	}

	c, err := pb.GetClient(clientID)
	if err != nil {
		panic(err)
	}

	ctx := trace.FromContext(nil)
	log.Printf("%s||request begin", ctx)

	r, err := c.SayHello(ctx.RpcContext(), &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.GetMessage())
}
