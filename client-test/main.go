package main

import (
	"encoding/json"
	"fmt"
	pb "microtest/user/proto/user"
	"os"
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

}
