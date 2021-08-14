module github.com/leejamin/conn-pool-go/examples/grpc

go 1.14

require (
	github.com/golang/protobuf v1.4.2
	github.com/leejamin/conn-pool-go v0.0.0-20210814165124-bbd0f5278890
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/leejamin/conn-pool-go => ../../
