package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/IbrohimIsroilov/shippy/shippy-service-consignment/proto/consignment"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
	defaultFileName = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {
	// set up a connection to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewShippingService(conn)
    
	// contact the server and print out the response
	file := defaultFileName
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)

	if err != nil {
		log.Fatalf("Coult not parse file: %v", err)
	}

	r, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatal("Could not greet: %v", err)
	}
	log.Printf("Created: %t", r.Created)
}