package main

import (
	"fmt"
	"log"
	"net"
	"context"
	"time"
	"os"
	"strconv"
	"math/rand"
	"reflect"
	pb "github.com/sirbernal/lab2SD/proto/client_service"
	pb2 "github.com/sirbernal/lab2SD/proto/node_service"
	"google.golang.org/grpc"
)
//datanodes := ["10.10.28.82:50051","10.10.28.83:50051","10.10.28.84:50051"]




type server struct {
}

type reg struct{
	nombre string
	direccion string //o cantidad 
}

var chunks [][]byte
var total int64 
var nombrearchivo string 
var datanode = []string{"localhost:50052","localhost:50053","localhost:50054"}
var datanodestatus = []bool{false,false,false}
var activos []int

func Unchunker(name string){
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
}
var registro [][]reg
func GuardarLibro(name string, partes int){
	var libro []reg
	libro = append(libro,reg{name,strconv.Itoa(partes)})
	for i:=0; i < partes; i++ {
		libro = append(libro,reg{name+"_parte_"+strconv.Itoa(i),"aqui va una IP"})
	}
	registro=append(registro,libro)
	fmt.Println(registro)
}
func ActualizarLibro(){
	file,err:= os.OpenFile("registro.txt",os.O_CREATE|os.O_WRONLY,0777) //actualiza archivo de registro
	defer file.Close()
	if err !=nil{
		os.Exit(1)
	}
	for _,libro :=range registro{
		for i,detalle :=range libro{
			var word string
			if i==0{
				word=detalle.nombre+" Cantidad_Partes_"+detalle.direccion
			}else{
				word=detalle.nombre+" "+detalle.direccion
			}
			_, err := file.WriteString(word + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	file.Close()
}
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

func (s *server) UploadChunks(ctx context.Context, msg *pb.UploadChunksRequest) (*pb.UploadChunksResponse, error) {
	fmt.Println(len(chunks))


	fmt.Println("Recibido")
	//fmt.Println(msg.GetChunk())
	chunks=append(chunks,msg.GetChunk())
	//ItsAlive()
	
	if int64(len(chunks))==total{
		go func(){
			Unchunker(nombrearchivo)
			GuardarLibro(nombrearchivo,int(total))
			ActualizarLibro()
		}()
	}
	return &pb.UploadChunksResponse{Resp : "recibido en el server", }, nil
	
}
func (s *server) Upload(ctx context.Context, msg *pb.UploadRequest) (*pb.UploadResponse, error) {
	nombrearchivo=msg.GetNombre()
	total=msg.GetTotalchunks()
	fmt.Println(msg.GetNombre())
	fmt.Println(msg.GetTotalchunks())
	return &pb.UploadResponse{Resp : int64(0), }, nil
}
func (s *server) Alive(ctx context.Context, msg *pb2.AliveRequest) (*pb2.AliveResponse, error) {
	return &pb2.AliveResponse{Msg : "Im Alive, namenode", }, nil
}

func (s *server) Distribucion(ctx context.Context, msg *pb2.DistribucionRequest) (*pb2.DistribucionResponse, error) {
	/* Esta interfaz se tiene que definir porque esta en el 2do proto, pero en este lado no lo usaremos asi que
	no tiene sentido lo que esta aca , porciacaso*/
	return &pb2.DistribucionResponse{Resp : "",}, nil
}

func TotalConectados()int{
	activos= []int{}
	cont:=0
	for i,data:= range datanodestatus{
		if data{
			cont++
			activos=append(activos, i)
		}
	}
	fmt.Println(activos)
	return cont
}
func GenerarPropuestaNueva (total int, conectados int)([]int64){
	var propuesta []int64
	for i:=0;i<=total/conectados;i++{
		rand.Seed(time.Now().UnixNano())
		lilprop:=rand.Perm(conectados)
		if i==total/conectados{
			sobra:=total%conectados
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
	switch conectados{
	case 1:
		if activos[0]==0{
			return propuesta
		}else{
			for i,_:= range propuesta{
				propuesta[i]=int64(activos[0])
			}
			return propuesta
		}
	case 2:
		if reflect.DeepEqual(activos,[]int{0,2}){
			for i,j:=range propuesta{
				if j==1{
					propuesta[i]=2
				}
			}
			return propuesta
		}else if reflect.DeepEqual(activos,[]int{1,2}){
			for i,j:=range propuesta{
				if j==0{
					propuesta[i]=2
				}
			}
			return propuesta
		}else{
			return propuesta
		}
	default:
		return propuesta
	}
}

func AllAlive () (bool){
	for j,dire :=range datanode{
		conn, err := grpc.Dial(dire, grpc.WithInsecure())
		if err != nil {
			fmt.Println("No esta el datanode", j+1)
			datanodestatus[j]=false
			continue
		}
		defer conn.Close()

		client := pb2.NewNodeServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		msg:= &pb2.AliveRequest{Msg: "Are u alive?"}

		resp, err := client.Alive(ctx, msg)
		if err != nil {
			fmt.Println("No esta el datanode", j+1)
			datanodestatus[j]=false
			continue
		}
		datanodestatus[j]=true
		fmt.Println(resp.GetMsg())
	}
	
	if reflect.DeepEqual(datanodestatus,[]bool{true,true,true}){
		fmt.Println("Todos los datanodes estan vivos")
		return true
	}else{
		return false
	}
}



func (s *server) Propuesta(ctx context.Context, msg *pb2.PropuestaRequest) (*pb2.PropuestaResponse, error) {
	
	/* RECEPCION DE PROPUESTA DE DATANODE */

	fmt.Println("Recibida propuesta!")
	fmt.Println(msg.GetProp())

	/* VERIFICAR QUE LOS NODOS DE LA PROPUESTA ESTEN ALIVE*/
	//allalive := AllAlive()
	//fmt.Println(allalive)
	if AllAlive() {
		return &pb2.PropuestaResponse{Msg : true, Prop : []int64{}}, nil
	} else {
		return &pb2.PropuestaResponse{Msg : false, Prop : GenerarPropuestaNueva(len(msg.GetProp()),TotalConectados())}, nil
	}
		
}

func (s *server) DownloadNames(ctx context.Context, msg *pb.DownloadNamesRequest) (*pb.DownloadNamesResponse, error) {

	return &pb.DownloadNamesResponse{Names : []string{} }, nil
}

func (s *server) DownloadChunks(ctx context.Context, msg *pb.DownloadChunksRequest) (*pb.DownloadChunksResponse, error) {

	return &pb.DownloadChunksResponse{Prop : []int64{} }, nil
}


func main()  {
	lis, err := net.Listen("tcp", ":50055")
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