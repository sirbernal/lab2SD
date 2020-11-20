package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/sirbernal/lab2SD/proto/client_service"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) Upload(stream pb.ClientService_UploadServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	fmt.Println(req.Chunk)
	rsp := &pb.UploadResponse{IdLibro : "recibido", }
	stream.Send(rsp)
	return nil
	
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