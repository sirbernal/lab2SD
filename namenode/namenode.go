package main

import (
	"fmt"
	"log"
	"net"
	"context"

	pb "github.com/sirbernal/lab2SD/proto/client_service"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) Upload(ctx context.Context, msg *pb.UploadRequest) (*pb.UploadResponse, error) {


	fmt.Println("Recibido")
	//fmt.Println(msg.GetChunk())
	return &pb.UploadResponse{IdLibro : "recibido en el server", }, nil
	
}



func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Error conectando: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterClientServiceServer(s, &server{}) //recibe conexión con el camion
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	
}