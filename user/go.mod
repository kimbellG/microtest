module github.com/kimbellG/microtest/user

go 1.16

replace github.com/kimbellG/microtest/user/ => ../user

require (
	github.com/golang/protobuf v1.5.2
	github.com/micro/go-micro/v2 v2.9.1
	google.golang.org/protobuf v1.26.0
)
