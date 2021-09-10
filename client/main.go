/**
 * Auth :   liubo
 * Date :   2021/9/10 13:44
 * Comment: 支持自动重连，放心用
 */
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	pb "google.golang.org/grpc/examples/features/proto/echo"
)

var addr = flag.String("addr", "localhost:50052", "the address to connect to")

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	for i:=0; i<100; i++ {
		c := pb.NewEchoClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		fmt.Println(i, " Performing unary request")
		res, err := c.UnaryEcho(ctx, &pb.EchoRequest{Message: "keepalive demo"})
		if err != nil {
			fmt.Println(i, " unexpected error from UnaryEcho: %v", err)
		}
		fmt.Println(i, " RPC response:", res)

		time.Sleep(time.Second * 1)
	}

	select {} // Block forever; run with GODEBUG=http2debug=2 to observe ping frames and GOAWAYs due to idleness.
}