package main

import (
	"fmt"
	"log"
	"net"
	"context"
	//"os"
	pb "github.com/sirbernal/lab2SD/proto/client_service"
	"google.golang.org/grpc"
)

type server struct {
}

var chunks [][]byte
/* func Unchunker(name string){
	_, err := os.Create(name)
	if err != nil {
			fmt.Println(err)
			os.Exit(1)
	}
	file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _,chunk:= range chunks{
		_, err := file.Write(chunk)	
		if err != nil {
				fmt.Println(err)
				os.Exit(1)
		}
		file.Sync()
	}
	file.Close()
	fmt.Println(len(chunks))
	chunks = [][]byte{}
	fmt.Println(len(chunks))
} */

func (s *server) Upload(ctx context.Context, msg *pb.UploadRequest) (*pb.UploadResponse, error) {
	fmt.Println(len(chunks))
	fmt.Println("Recibido")
	//fmt.Println(msg.GetChunk())
	chunks=append(chunks,msg.GetChunk())
	// generarpropuesta(chunks)

	/*if len(chunks)==4{
		go func(){
			Unchunker("ejemploarmado.pdf")
		}()
	} */
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