/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"fmt"
	"github.com/leejamin/conn-pool-go/pool"
	"log"
	"os"
	"time"

	pb "github.com/leejamin/conn-pool-go/examples/grpc/helloworld/helloworld"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	connPool := pool.NewConnPool(&pool.Options{
		//Addr:         "www.baidu.com:80",
		MinIdleConns: 2,
		Dialer: func(ctx context.Context) (pool.Conn, error) {
			return grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		},
		OnClose: func(conn pool.Conn) error {
			fmt.Println("close")
			return nil
		},
	})

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()


	for i := 0; i < 5; i++ {
		connPool.WithConn(ctx, func(ctx context.Context, conn pool.Conn) error {

			if grpcConn, ok := conn.(*grpc.ClientConn); ok {
				c := pb.NewGreeterClient(grpcConn)
				r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
				if err != nil {
					log.Fatalf("could not greet: %v", err)
				}
				log.Printf("Greeting: %s", r.GetMessage())
				return err
			}
			return nil
		})
		time.Sleep(100 * time.Millisecond)
	}
	connPool.Close()

}
