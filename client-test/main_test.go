package main

import (
	"context"
	"testing"

	pb "github.com/kimbellG/microtest/user/proto/user"

	"google.golang.org/grpc"
)

func TestUser(t *testing.T) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	file := defaultFilename

	user, err := parseFile(file)
	if err != nil {
		t.Fatalf("Could not parse file: %v", err)
	}

	r, err := client.SignUp(context.Background(), user)
	if err != nil {
		t.Fatalf("Could not sign up: %v", err)
	}
	t.Log(r.ErrorString)

}
