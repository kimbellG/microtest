package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "microtest/user/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":12345"
)

type repo interface {
	AddUser(*pb.UserInfo) (*pb.UserInfo, error)
}

type Repository struct {
	mu sync.RWMutex
}

func (r *Repository) AddUser(user *pb.UserInfo) (*pb.UserInfo, error) {
	r.mu.Lock()
	fmt.Printf("new user:\n\temail: %v\n\tSN: %v %v\n", user.Email, user.Surname, user.Name)
	r.mu.Unlock()

	return user, nil
}

type service struct {
	r repo
}

func (s *service) SignUp(ctx context.Context, user *pb.UserInfo) (*pb.Response, error) {
	_, err := s.r.AddUser(user)
	if err != nil {
		return nil, err
	}

	return &pb.Response{
		IsCreated:   true,
		ErrorString: "Успех.",
	}, nil
}

func main() {
	r := &Repository{}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterAuthServiceServer(s, &service{r})
	reflection.Register(s)

	log.Println("Running on port: ", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
