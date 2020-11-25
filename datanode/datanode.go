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

//namenode := "10.10.28.81:50051"
//datanodes := ["10.10.28.82:50051","10.10.28.83:50051","10.10.28.84:50051"]
var total int64 
var nombrearchivo string
var cont int64

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
	nombrearchivo=msg.GetNombre()
	total=msg.GetTotalchunks()
	fmt.Println(msg.GetNombre())
	fmt.Println(msg.GetTotalchunks())
	return &pb.UploadResponse{Resp : int64(0), }, nil
}

func (s *server) UploadChunks(ctx context.Context, msg *pb.UploadChunksRequest) (*pb.UploadChunksResponse, error) {
	fmt.Println(len(chunks))


	fmt.Println("Recibido")
	//fmt.Println(msg.GetChunk())
	chunks=append(chunks,msg.GetChunk())
	
	if int64(len(chunks))==total{
		
		/* GGENERAR PROPUESTA*/
		estado := false
		if estado == false{
			
			conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
			if err != nil {
				log.Fatalln(err)
			}
			defer conn.Close()

			client := pb2.NewNodeServiceClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			msg:= &pb2.PropuestaRequest{Msg: "Propuesta de ejemplo"}

			resp, err := client.Propuesta(ctx, msg)
			estado = resp.GetMsg()
		}

		/* Enviar mensaje devuelta a cliente diciendole que ya terminamos la wea */
		chunks = [][]byte{}
		return &pb.UploadChunksResponse{Resp : "El servidor acepto su propuesta pete", }, nil
			
		
		
		/*go func(){
			Unchunker(nombrearchivo)
			GuardarLibro(nombrearchivo,int(total))
			ActualizarLibro()
		}()*/
	}
	return &pb.UploadChunksResponse{Resp : "recibido en el server", }, nil
	
}

func (s *server) Alive(ctx context.Context, msg *pb2.AliveRequest) (*pb2.AliveResponse, error) {
	return &pb2.AliveResponse{Msg : "Im Alive bitch", }, nil
}

func (s *server) Propuesta(ctx context.Context, msg *pb2.PropuestaRequest) (*pb2.PropuestaResponse, error) {
	
	/* Esta interfaz se tiene que definir porque esta en el 2do proto, pero en este lado no lo usaremos asi que
	no tiene sentido lo que esta aca , porciacaso*/

	return &pb2.PropuestaResponse{Msg : true,}, nil
}

func main() {
	
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal("Error conectando: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterClientServiceServer(s, &server{})
	pb2.RegisterNodeServiceServer(s, &server{}) //recibe conexión con el camion
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	
}