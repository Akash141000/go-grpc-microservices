package main

import (
	"context"
	"flag"
	"fmt"
	"grpc-go/client"
	"grpc-go/proto"
	"log"
	"time"
)

func runClient() {
	client := client.New("http://localhost:3000")
	price, err := client.FetchPrice(context.Background(), "ET")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", price)
}

func runGRPCServer() {
	listenAddr := flag.String("listenAddr", ":4000", "Port on which grpc service is running on")
	flag.Parse()
	svc := NewLogginService(NewMetricService(&priceFetcher{}))

	err := makeGRPCServerAndRun(*listenAddr, svc)
	if err != nil {
		log.Fatal(err)
	}

}

func runGRPCClient() {
	listenAddr := flag.String("listenAddr", ":4000", "Port on which grpc client is running on")
	flag.Parse()

	ctx := context.Background()

	grpcClient, err := client.NewGRPCClient(*listenAddr)

	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	resp, err := grpcClient.FetchPrice(ctx, &proto.PriceRequest{
		Ticker: "ETH",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", resp)

}

func runJSONServer() {
	listenAddr := flag.String("listenAddr", ":3000", "Port on which service is running on")
	flag.Parse()

	svc := NewLogginService(NewMetricService(&priceFetcher{}))

	server := NewJSONAPIServer(*listenAddr, svc)
	server.Run()

	log.Println("Go server started...")
}

func main() {
	runClient()
	runJSONServer()
	runGRPCServer()
	runGRPCClient()
}
