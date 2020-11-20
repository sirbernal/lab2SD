package main

import (
	"context"
	"fmt"
	"log"




	pb "github.com/sirbernal/lab2SD/proto/client_service"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	
	client := pb.NewClientServiceClient(conn)
 

	msg := &pb.UploadRequest{IdLibro: "enviado", }

	stream, err := client.Upload(context.Background())

	stream.Send(msg)
	resp, err := stream.Recv()

	if err != nil {
		log.Fatalf("can not receive %v", err)
	}
	
	fmt.Println(resp.IdLibro)
}

