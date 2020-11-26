package main

import (
	"fmt"
	"log"
	"net"
	"context"
	"time"
	"io/ioutil"
	"strconv"
	"os"
	"math/rand"
	pb "github.com/sirbernal/lab2SD/proto/client_service"
	pb2 "github.com/sirbernal/lab2SD/proto/node_service"
	"google.golang.org/grpc"
)

type server struct {
}

//namenode := "10.10.28.81:50051"
var datanode = []string{"localhost:50052","localhost:50053","localhost:50054"}
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
func SaveChunk (chunk []byte, name string){
	fileName := name
	_, err := os.Create(fileName)
	if err != nil {
			fmt.Println(err)
			os.Exit(1)
	}
	// write/save buffer to disk
	ioutil.WriteFile(fileName, chunk, os.ModeAppend)
}
func SearchChunk (name string) (chunk []byte){
	newFileChunk, err := os.Open(name)
	if err != nil {
			fmt.Println(err)
			os.Exit(1)
	}
	defer newFileChunk.Close()

	chunkInfo, err := newFileChunk.Stat()
	if err != nil {
			fmt.Println(err)
			os.Exit(1)
	}
	var chunkSize int64 = chunkInfo.Size()
	chunkBufferBytes := make([]byte, chunkSize)
	return chunkBufferBytes
}
func GenerarPropuesta (total int)([]int64){
	var propuesta []int64
	for i:=0;i<=total/3;i++{
		rand.Seed(time.Now().UnixNano())
		lilprop:=rand.Perm(3)
		if i==total/3{
			sobra:=total%3
			for j,num :=range lilprop{
				if j==sobra{
					break
				}
				propuesta=append(propuesta,int64(num))
			}
		}else{
			for _,num :=range lilprop{
				propuesta=append(propuesta,int64(num))
			}
		}
	}
	return propuesta
}

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
	
	if int64(len(chunks))==total{  // Cuando llegan todos los chunks del archivo al datanode, realizamos propuesta
		
		/* GGENERAR PROPUESTA*/
	
		conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
		if err != nil {
			log.Fatalln(err)
		}
		defer conn.Close()

		client := pb2.NewNodeServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		propuesta := GenerarPropuesta(int(total))
		msg:= &pb2.PropuestaRequest{Prop: propuesta, Name: nombrearchivo}

		resp, err := client.Propuesta(ctx, msg)
		//estado := resp.GetMsg()
		fmt.Println(resp.GetProp())

		if resp.GetMsg() == false{
			propuesta = resp.GetProp()
		}

		/* GENERAR DISTRIBUCION*/
		for i,j :=range propuesta{
			if datanode[1]==datanode[j]{
				SaveChunk(chunks[i],nombrearchivo+"_"+strconv.Itoa(i))
				continue
			}
			conn, err := grpc.Dial(datanode[j], grpc.WithInsecure())
			if err != nil {
				fmt.Println("Proceso abortado, se ha desconectado el nodo durante la distribucion")
				break
			}
			defer conn.Close()
	
			client := pb2.NewNodeServiceClient(conn)
	
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			msg:= &pb2.DistribucionRequest{Chunk: chunks[i], Name: nombrearchivo+"_"+strconv.Itoa(i)}
	
			_, err = client.Distribucion(ctx, msg)
			if err != nil {
				continue
			}
		}

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
	return &pb2.AliveResponse{Msg : "Im Alive, datanode1", }, nil
}

func (s *server) Propuesta(ctx context.Context, msg *pb2.PropuestaRequest) (*pb2.PropuestaResponse, error) {
	
	/* Esta interfaz se tiene que definir porque esta en el 2do proto, pero en este lado no lo usaremos asi que
	no tiene sentido lo que esta aca , porciacaso*/

	return &pb2.PropuestaResponse{Msg : true,}, nil
}
func (s *server) Distribucion(ctx context.Context, msg *pb2.DistribucionRequest) (*pb2.DistribucionResponse, error) {
	
	SaveChunk(msg.GetChunk(),msg.GetName())
	return &pb2.DistribucionResponse{Resp : "",}, nil
}

func (s *server) DownloadNames(ctx context.Context, msg *pb.DownloadNamesRequest) (*pb.DownloadNamesResponse, error) {

	return &pb.DownloadNamesResponse{Names : []string{} }, nil
}

func (s *server) DownloadChunks(ctx context.Context, msg *pb.DownloadChunksRequest) (*pb.DownloadChunksResponse, error) {

	return &pb.DownloadChunksResponse{Prop : []int64{} }, nil
}

func main() {
	
	lis, err := net.Listen("tcp", ":50053")
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