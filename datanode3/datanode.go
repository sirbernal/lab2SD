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
var tipo_distribucion string

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
	
	if tipo_distribucion == "centralizado"{  // Logica cuando el algoritmo de distrib. es centralizado
		
		chunks=append(chunks,msg.GetChunk()) // Agregamos el chunk a nuestro arreglo

		if int64(len(chunks))==total{ // Cuando llegan todos los chunks del archivo al datanode, realizamos propuesta
			
			/* GENERAR PROPUESTA*/
			/* Primero, generamos la conexion con namenode*/
			conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
			if err != nil {
				log.Fatalln(err)
			}
			defer conn.Close()

			client := pb2.NewNodeServiceClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			propuesta := GenerarPropuesta(int(total)) // Con esta funcion generaremos la propuesta
			msg:= &pb2.PropuestaRequest{Prop: propuesta, Name: nombrearchivo} // generamos el mensaje con la propuesta

			resp, err := client.Propuesta(ctx, msg) // enviamos la propuesta y recibimos la respuesta
			//estado := resp.GetMsg()
			fmt.Println(resp.GetProp())

			if resp.GetMsg() == false{  // cuando se rechaza la propuesta, la actualizamos con la propuesta recibida x namenode
				propuesta = resp.GetProp()
		}

		/* GENERAR DISTRIBUCION*/
		/* Leemos el arrigo de propuesta que tiene las designaciones de cada chunk que ira a cada datanode*/
		for i,j :=range propuesta{
			if datanode[0]==datanode[j]{ //claramente, si un chunk se debe quedar en este datanode, para que enviarlo XD
				SaveChunk(chunks[i],nombrearchivo+"_"+strconv.Itoa(i)) // Simplemente lo guardamos con la funcion savechunk
				continue
			}
			// Procedemos a generar la conexion con el datanode a donde enviaremos el chunk
			conn, err := grpc.Dial(datanode[j], grpc.WithInsecure())
			if err != nil {
				fmt.Println("Proceso abortado, se ha desconectado el nodo durante la distribucion")
				break
			}
			defer conn.Close()
	
			client := pb2.NewNodeServiceClient(conn)
	
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			// Creamos el mensaje con el chunk correspondiente + el nombre del chunk 
			msg:= &pb2.DistribucionRequest{Chunk: chunks[i], Name: nombrearchivo+"_"+strconv.Itoa(i)}
			
			// Enviamos el chunk
			_, err = client.Distribucion(ctx, msg)
			if err != nil {
				continue
			}
		}

		// Una vez realizada la distribucion, reseteamos el arreglo donde guardamos los chunks
		chunks = [][]byte{}
		return &pb.UploadChunksResponse{Resp : "El servidor guardo el archivo", }, nil //Avisamos de que estamos ok con la subida

		}
		// Esto se envia cuando no se tienen todos los chunks del archivo
		return &pb.UploadChunksResponse{Resp : "chunk recibido en el server", }, nil
	}
	
	return &pb.UploadChunksResponse{Resp : "Fallo algo", }, nil
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
	return &pb.DownloadChunksResponse{Chunk : SearchChunk(msg.GetName()) }, nil
}

func (s *server) LocationsofChunks(ctx context.Context, msg *pb.LoCRequest) (*pb.LoCResponse, error) {

	return &pb.LoCResponse{Location: []int64{} }, nil
}



func main() {
	
	lis, err := net.Listen("tcp", ":50054")
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