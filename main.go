package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	api "github.com/alanchchen/grpc-lb-istio/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	host   string
	port   int
	repeat int
)

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "The server host")
	flag.IntVar(&port, "port", 7000, "The server port")
	flag.IntVar(&repeat, "repeat", 1, "Times to call server")
	flag.Parse()
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewIdentityClient(conn)

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"Content-Type": "application/grpc",
	}))

	for i := 0; i < repeat; i++ {
		r, err := c.Who(ctx, &api.WhoRequest{})
		if err != nil {
			log.Printf("%v\n", err)
			continue
		}
		log.Printf("%s\n", r.Name)
	}
}
