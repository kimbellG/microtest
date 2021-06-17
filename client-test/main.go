package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	timer "github.com/kimbellG/microtest/timer/proto/timer"
	pb "github.com/kimbellG/microtest/user/proto/user"
	micro "github.com/micro/go-micro/v2"
)

const (
	address         = "localhost:12345"
	defaultFilename = "user.json"
)

func parseFile(file string) (*pb.UserInfo, error) {

	result := &pb.UserInfo{}
	fileStream, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("open: %v: %v", file, err)
	}

	if err := json.NewDecoder(fileStream).Decode(result); err != nil {
		return nil, fmt.Errorf("decode: %v: %v", file, err)
	}

	return result, nil
}

func main() {
	service := micro.NewService(micro.Name("user.cli"))
	service.Init()

	client := pb.NewAuthService("AuthService", service.Client())
	t := timer.NewTimerService("Timer", service.Client())

	file := defaultFilename

	test, err := parseFile(file)
	if err != nil {
		log.Fatalf("could not parse file: %v", err)
	}

	_, err = client.SignUp(context.Background(), test)
	if err != nil {
		log.Fatalf("sign up: %v", err)
	}

	fmt.Println("Start timer")
	resp, err := t.Wait(context.Background(), &timer.Request{T: int64(time.Second * 10)})

	if err != nil {
		log.Fatalln("wait timer:", err)
	}

	fmt.Println("Timer is end")
	fmt.Printf("resp: %t", resp.IsOK)

}
