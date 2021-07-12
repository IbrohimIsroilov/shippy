package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "gihub.com/IbrohimIsroilov/shippy/shippy-service-consignment/proto/consignment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
}

type Repository struct {
	mu  sync.RWMutex
	consignments []*pb.Consignment
}

// Create a new consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	repo.mu.Unlock()
	return consignment, nil
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition
type service struct {
	repo repository
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	// Save our consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	// Return matching the "Response" message we created in our protobuf definition
	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func main() {
	repo := &Repository{}

	// set-up our grpc server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// register our service with the grpc server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definiton

	pb.RegisterShippingServiceServer(s, &service{repo})

	// register reflection service on grpc server
	reflection.Register(s)

	log.Println("Running on port: ", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: %v", err)
	}
}