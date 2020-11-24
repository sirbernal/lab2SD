package main

import (
	"fmt"
	"log"
	"net"
	"context"
	"time"
	//"os"
	pb "github.com/sirbernal/lab2SD/proto/client_service"
	pb2 "github.com/sirbernal/lab2SD/proto/node_service"
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

func ItsAlive(){
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := pb2.NewNodeServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	msg:= &pb2.AliveRequest{Msg: "Are u alive?"}

	resp, err := client.Alive(ctx, msg)
	fmt.Println(resp.GetMsg())
}

func (s *server) Upload(ctx context.Context, msg *pb.UploadRequest) (*pb.UploadResponse, error) {
	fmt.Println(len(chunks))


	fmt.Println("Recibido")
	//fmt.Println(msg.GetChunk())
	chunks=append(chunks,msg.GetChunk())
	ItsAlive()
	
	/*if len(chunks)==4{
		go func(){
			Unchunker("ejemploarmado.pdf")
		}()
	} */
	return &pb.UploadResponse{IdLibro : "recibido en el server", }, nil
	
}

func (s *server) Alive(ctx context.Context, msg *pb2.AliveRequest) (*pb2.AliveResponse, error) {
	return &pb2.AliveResponse{Msg : "Im Alive bitch", }, nil
}



func main()  {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Error conectando: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterClientServiceServer(s, &server{}) //recibe conexión con el camion
	pb2.RegisterNodeServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}



/*
	chblock := make(chan struct{}) 
	go func(){
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatal("Error conectando: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterClientServiceServer(s, &server{}) //recibe conexión con el camion
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	<-chblock
	*/
}