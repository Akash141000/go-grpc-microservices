package main

import (
	"context"
	"flag"
	"fmt"
	"grpc-go/client"
	"log"
)

func runClient() {
	client := client.New("http://localhost:3000")
	price, err := client.FetchPrice(context.Background(), "ET")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", price)
}

func runServer() {
	listenAddr := flag.String("listenAddr", ":3000", "Port on which service is running on")
	flag.Parse()

	svc := NewLogginService(NewMetricService(&priceFetcher{}))

	server := NewJSONAPIServer(*listenAddr, svc)
	server.Run()

	log.Println("Go server started...")
}

func main() {
	runClient()
	runServer()
}
