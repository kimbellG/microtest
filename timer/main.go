package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/kimbellG/microtest/timer/proto/timer"

	micro "github.com/micro/go-micro/v2"
)

const (
	rewriteChar = '\b'
)

func getBackspaceInterval(t time.Duration) string {
	result := ""
	for len(result) <= len(fmt.Sprintf("%.1f", float32(t/time.Second))) {
		result += "\b"
	}

	return result
}

func timer(req *pb.TimerRequest) pb.TimerResponse {

	result := pb.TimerResponse{}

	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1))
		abort <- struct{}{}
	}()

	fmt.Printf("Timer start. End of timer after %.1fs\n", float32(time.Duration(req.T)/time.Second))
	fmt.Println("For stop press any key.")
	fmt.Print("Current time: ")

	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	resp := ticker.C
	for cur, backspace := time.Duration(0), "\b\b\b\b"; cur <= time.Duration(req.T); {
		select {
		case <-resp:
			backspace = getBackspaceInterval(cur)
			fmt.Printf("%.1fs%v", float32(cur)/float32(time.Second), backspace)
			cur += time.Millisecond * 100
		case <-abort:
			fmt.Println("\nStop timer. Goodbye")
			result.IsOK = false
			return result
		}
	}

	fmt.Println("\a\nTimer end. Goodbye")
	result.IsOK = true

	return result
}

type timerService struct{}

func (t *timerService) Wait(ctx context.Context, r *pb.TimerRequest, resp *pb.TimerResponse) error {
	result := timer(r)
	resp.IsOK = result.IsOK

	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("Timer"),
	)

	service.Init()
	if err := pb.RegisterTimerHandler(service.Server(), &timerService{}); err != nil {
		log.Fatalf("register service handler: %v", err)
	}

	if err := service.Run(); err != nil {
		log.Fatalf("run service: %v", err)
	}
}
