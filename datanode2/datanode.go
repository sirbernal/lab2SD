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
	"bufio"
	"reflect"
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
var tipo_distribucion string
var this_datanode = datanode[1]
var activos []int
var datanodestatus = []bool{false,false,false}
//BORRAR VARIABLES TEMPORALES
var registroname []string
var registroprop [][]int64
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
	reader := bufio.NewReader(newFileChunk)
	reader.Read(chunkBufferBytes)
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

		
	}else if tipo_distribucion == "distribuida" {
		chunks=append(chunks,msg.GetChunk()) // Agregamos el chunk a nuestro arreglo

		if int64(len(chunks))==total{// Cuando llegan todos los chunks del archivo al datanode, realizamos propuesta
			
		
			/* Logica de esto, es que si una propuesta se rechaza, pues hay que volver a realizar una
			nueva propuesta, por lo que rompemos el segundo for que esta hecho para enviar la propuesta a cada
			datanode y volvemos a realizarlo con una nueva propuesta, para eso usamos una variable booleana
			proceso que cuando se aceptan todas las propuestas, recien procederiamos a distribuir */
			var propuesta []int64
			proceso := true  
			for proceso{
				for _,dire:= range datanode{
					/* GENERAR PROPUESTA*/
					/* Primero, generamos la conexion con cada datanode*/
					if this_datanode == dire{ // No enviaremos a este mismo nodo la propuesta a generar
						continue
					}
					conn, err := grpc.Dial(dire, grpc.WithInsecure()) // enviamos propuesta a un nodo especifico
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
						break 
					}

				}
			}

			/* GENERAR DISTRIBUCION*/
			/* Leemos el arrigo de propuesta que tiene las designaciones de cada chunk que ira a cada datanode*/
			for i,j :=range propuesta{
				if this_datanode==datanode[j]{ //claramente, si un chunk se debe quedar en este datanode, para que enviarlo XD
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

			chunks = [][]byte{}
			return &pb.UploadChunksResponse{Resp : "El servidor guardo el archivo", }, nil //Avisamos de que estamos ok con la subida
		}

		return &pb.UploadChunksResponse{Resp : "chunk recibido en el server", }, nil

	}
	
	return &pb.UploadChunksResponse{Resp : "Fallo algo", }, nil
}

func (s *server) Alive(ctx context.Context, msg *pb2.AliveRequest) (*pb2.AliveResponse, error) {
	return &pb2.AliveResponse{Msg : "Im Alive, datanode1", }, nil
}
func ActualizarRegistro(){
	file,err:= os.OpenFile("registro.txt",os.O_CREATE|os.O_WRONLY,0777) //actualiza archivo de registro
	defer file.Close()
	if err !=nil{
		os.Exit(1)
	}
	for i,propuesta :=range registroprop{
		word:=registroname[i]+" Cantidad_Partes_"+strconv.Itoa(len(propuesta))
		_, err := file.WriteString(word + "\n")
		if err != nil {
			log.Fatal(err)
		}
		for j,node:= range propuesta{
			word="parte_"+strconv.Itoa(j)+" "+datanode[int(node)]
			_, err := file.WriteString(word + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	file.Close()
}
func GuardarPropuesta(name string, partes []int64){
	registroname=append(registroname,name)
	registroprop=append(registroprop,partes)
	ActualizarRegistro()
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
		GuardarPropuesta(msg.GetName(),msg.GetProp())
		return &pb2.PropuestaResponse{Msg : true, Prop : []int64{}}, nil
	} else {
		nuevaprop:=GenerarPropuestaNueva(len(msg.GetProp()),TotalConectados())
		GuardarPropuesta(msg.GetName(),nuevaprop)
		return &pb2.PropuestaResponse{Msg : false, Prop : nuevaprop}, nil
	}
		
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

func (s *server) TypeDis(ctx context.Context, msg *pb.TypeRequest) (*pb.TypeResponse, error) {
	if msg.GetType()=="inicio"{
		return &pb.TypeResponse{Resp: tipo_distribucion}, nil
	}else{
		tipo_distribucion = msg.GetType()
		fmt.Println("Tipo distribucion: ", tipo_distribucion)
	}
	return &pb.TypeResponse{Resp: "ok" }, nil
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