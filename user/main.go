package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	pb "microtest/user/proto/user"

	"github.com/micro/go-micro/v2"
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

type userService struct {
	r repo
}

func (s *userService) SignUp(ctx context.Context, user *pb.UserInfo, resp *pb.Response) error {
	_, err := s.r.AddUser(user)
	if err != nil {
		return err
	}

	resp.IsCreated = true
	resp.ErrorString = "Успех"
	return nil
}

func main() {
	r := &Repository{}

	service := micro.NewService(
		micro.Name("user.service"),
	)

	service.Init()

	if err := pb.RegisterAuthServiceHandler(service.Server(), &userService{r}); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
