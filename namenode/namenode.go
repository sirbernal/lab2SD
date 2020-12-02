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
var total int64 //guarda el numero total de chunks de una propuesta
var nombrearchivo string //guarda el nombre del archivo de la propuesta a guardar
var datanode = []string{"localhost:50052","localhost:50053","localhost:50054"} //arreglo que contiene las direcciones de los datanodes
var directions = []string{"localhost:50055", "localhost:50052","localhost:50053","localhost:50054"} //arreglo que incluye al namenode (se agrego despues de muchas funciones ya hechas)
var this_datanode = directions[0] //direccion de namenode
var datanodestatus = []bool{false,false,false} //arreglo que guarda el estado de un nodo= true:conectado false:desconectado
var activos []int //guarda cuales son los nodos activos
var tipo_distribucion string //se guarda el tipo de distribucion del namenode
var estado =true // true: libre , false: ocupado
var registroname []string //registro en memoria que guarda los nombres de los archivos
var registroprop [][]int64 //registro en memoria que guarda la distribucion de chunks de archivos
func GuardarPropuesta(name string, partes []int64){ //guarda la propuesta en memoria
	registroname=append(registroname,name) //guarda el nombre de la propuesta
	registroprop=append(registroprop,partes) //guarda la distribucion de chunks de la propuesta
	ActualizarRegistro() //actualiza el registro agregando el ultimo archivo
}
func ActualizarRegistro(){//actualiza el registro en base a lo que tiene en memoria
	file,err:= os.OpenFile("registro.txt",os.O_CREATE|os.O_WRONLY,0777) //abre o genera el archivo de registro
	defer file.Close()
	if err !=nil{
		os.Exit(1)
	}
	for i,propuesta :=range registroprop{ //revisa cada nombre de registro guardado en memoria
		word:=registroname[i]+" Cantidad_Partes_"+strconv.Itoa(len(propuesta))//primero escribe el encabezado solicitado en pauta
		_, err := file.WriteString(word + "\n")
		if err != nil {
			log.Fatal(err)
		}
		for j,node:= range propuesta{ //revisa donde esta cada parte del archivo a guardar en memoria
			word="parte_"+strconv.Itoa(j)+" "+datanode[int(node)] //escribe la parte y la direccion del nodo que almacena 
			_, err := file.WriteString(word + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	file.Close()
}

func (s *server) UploadChunks(ctx context.Context, msg *pb.UploadChunksRequest) (*pb.UploadChunksResponse, error) {//funcion de datanode
	return &pb.UploadChunksResponse{Resp : "", }, nil	
}

func (s *server) Upload(ctx context.Context, msg *pb.UploadRequest) (*pb.UploadResponse, error) {//funcion de datanode
	return &pb.UploadResponse{Resp : int64(0), }, nil
}
func (s *server) Alive(ctx context.Context, msg *pb2.AliveRequest) (*pb2.AliveResponse, error) { //retorna que esta activo el nodo si le consultan
	return &pb2.AliveResponse{Msg : "Im Alive, namenode", }, nil
}

func (s *server) Distribucion(ctx context.Context, msg *pb2.DistribucionRequest) (*pb2.DistribucionResponse, error) {//funcion de datanode
	return &pb2.DistribucionResponse{Resp : "",}, nil
}

func TotalConectados()int{//funcion que calcula el numero total de conectados y actualiza cuales nodos son los activos
	activos= []int{}//reinicia la lista de activos
	cont:=0//contador de activos
	for i,data:= range datanodestatus{
		if data{
			cont++
			activos=append(activos, i)//guarda el nodo activo
		}
	}
	return cont//retorna la cantidad
}
func GenerarPropuestaNueva (total int, conectados int)([]int64){//en modo centralizado genera una nueva propuesta basada en los conectados actualmente (si hay menos de 3)
	var propuesta []int64 //nueva propuesta
	for i:=0;i<=total/conectados;i++{ //genera la propuesta para cada chunk en base al total/total de conectados, es decir si hay dos conectados, lo hara por cada dos chunks
		rand.Seed(time.Now().UnixNano())//genera una semilla random basada en el time de la maquina
		lilprop:=rand.Perm(conectados)//genera un arreglo al azar de total/total de conectados para distribuir, por ejemplo, si hay dos podria salir (0,1) o (1,0)
		if i==total/conectados{//detecta la ultima serie de chunks
			sobra:=total%conectados
			for j,num :=range lilprop{
				if j==sobra{
					break//si ya no hay mas chunks termina la generacion de propuesta
				}
				propuesta=append(propuesta,int64(num))//agrega el nodo al cual correspondera la propuesta
			}
		}else{
			for _,num :=range lilprop{
				propuesta=append(propuesta,int64(num))//agrega el nodo al cual correspondera la propuesta
			}
		}
	}
	switch conectados{//en base a la cantidad de conectados vera como enviar la propuesta
	case 1: //si hay un conectado
		if activos[0]==0{ //si es el datanode 1 solo manda la propuesta
			return propuesta
		}else{ //si es cualquiera de los otros dos nodos cambia cada valor de la propuesta para este nodo
			for i,_:= range propuesta{ //ej prop=000 , conectado=1 , nueva prop=111
				propuesta[i]=int64(activos[0])
			}
			return propuesta
		}
	case 2: //si hay dos conectados
		if reflect.DeepEqual(activos,[]int{0,2}){ //si los conectados son el datanode 1 y datanode 3
			for i,j:=range propuesta{ 
				if j==1{
					propuesta[i]=2 //reemplaza los designados al datanode 2 por los del datanode 3
				}
			}
			return propuesta
		}else if reflect.DeepEqual(activos,[]int{1,2}){ //si los conectados son el datanode 2 y datanode 3
			for i,j:=range propuesta{
				if j==0{
					propuesta[i]=2 //reemplaza los designados al datanode 1 por los del datanode 3
				}
			}
			return propuesta
		}else{ //si los conectados son el datanode 1 y datanode 2 solo envia la propuesta ya que no requiere cambios
			return propuesta
		}
	default:
		return propuesta
	}
}

func AllAlive () (bool){ //funcion que verifica si estan todos los datanodes conectados, actualiza en memoria cuales lo estan y retorna falso si hay al menos uno que no lo está
	for j,dire :=range datanode{
		conn, err := grpc.Dial(dire, grpc.WithInsecure())//inicia conexion con cada datanode
		if err != nil {
			datanodestatus[j]=false
			continue
		}
		defer conn.Close()
		client := pb2.NewNodeServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		msg:= &pb2.AliveRequest{Msg: "Are u alive?"} //realiza la consulta al datanode
		_ , err = client.Alive(ctx, msg)
		if err != nil {
			datanodestatus[j]=false
			continue
		}
		datanodestatus[j]=true //actualiza la lista de estados
	}
	if reflect.DeepEqual(datanodestatus,[]bool{true,true,true}){//si estan todos conectados retorna true
		fmt.Println(datanodestatus)
		return true
	}else{
		fmt.Println(datanodestatus)
		return false
	}
}


func (s *server) Propuesta(ctx context.Context, msg *pb2.PropuestaRequest) (*pb2.PropuestaResponse, error) {//funcion que recibe la propuesta del datanode
	
	/* RECEPCION DE PROPUESTA DE DATANODE */
	if tipo_distribucion == "centralizado" { //si esta en modo centralizado
		estado=false //se pone en estado ocupado
		if AllAlive() { //verifica que los datanodes este en linea
			GuardarPropuesta(msg.GetName(),msg.GetProp()) //si es asi acepta y guarda la propuesta
			return &pb2.PropuestaResponse{Msg : true, Prop : []int64{}}, nil
		} else {// si hay al menos un datanode offline
			nuevaprop:=GenerarPropuestaNueva(len(msg.GetProp()),TotalConectados()) //genera una nueva propuesta en base a los datanodes conectados
			GuardarPropuesta(msg.GetName(),nuevaprop) //guarda la nueva propuesta 
			return &pb2.PropuestaResponse{Msg : false, Prop : nuevaprop}, nil //niega la propuesta inicial y envia la nueva
		}
	} else if tipo_distribucion == "distribuido"{//si esta en modo distribuido
		GuardarPropuesta(msg.GetName(),msg.GetProp()) //guarda la propuesta recibida
		return &pb2.PropuestaResponse{Msg : true, Prop : []int64{}}, nil
	}	
	return &pb2.PropuestaResponse{Msg : true, Prop : []int64{}}, nil //return solicitado por golang
}

func (s *server) DownloadNames(ctx context.Context, msg *pb.DownloadNamesRequest) (*pb.DownloadNamesResponse, error) {//funcion que retorna al cliente los nombres de archivos disponibles
	return &pb.DownloadNamesResponse{Names : registroname }, nil
}
func VerifProp(prop []int64)bool{ //funcion que verifica que en la propuesta generada esten todos los nodos activos
	AllAlive() //actualiza la lista de nodos activos
	for _,chunk:=range prop{
		if !datanodestatus[chunk]{
			return false //si dentro de la propuesta encuentra un nodo apagado designado retorna falso
		}
	}
	return true 
}
func BuscarChunks(name string)([]int64){ //funcion que busca la distribucion de un archivo
	var propreq= []int64{} //genera una distribucion vacia
	for i,reg :=range registroname{ //busca la distribucion acorde al nombre
		if reg==name{
			propreq=registroprop[i] //guarda la distribucion encontrada
		}
	}
	if reflect.DeepEqual(propreq,[]int64{}){ //si no encuentra nada retorna vacio
		return []int64{}
	}
	if !VerifProp(propreq){ //si dentro de la distribucion hay un nodo desconectado retorna una distribucion vacia para notificar que es imposible rearmar el archivo
		return []int64{}
	}
	return propreq
}
func (s *server) DownloadChunks(ctx context.Context, msg *pb.DownloadChunksRequest) (*pb.DownloadChunksResponse, error) {//funcion de los datanodes
	return &pb.DownloadChunksResponse{Chunk : []byte{} }, nil
}
func (s *server) LocationsofChunks(ctx context.Context, msg *pb.LoCRequest) (*pb.LoCResponse, error) {//funcion que retorna la ubicacion de las partes del archivo solicitado por el cliente
	return &pb.LoCResponse{Location: BuscarChunks(msg.GetReq())} , nil
}
func (s *server) TypeDis(ctx context.Context, msg *pb.TypeRequest) (*pb.TypeResponse, error) {//funcion para multiples usos
	if msg.GetType()=="inicio"{//si el cliente o un nodo esta iniciando y este esta corriendo retornara el tipo de distribucion que actualmente posee
		return &pb.TypeResponse{Resp: tipo_distribucion}, nil
	}else if msg.GetType()=="status"{ //funcion que notifica que se encuentra online el nodo
		return &pb.TypeResponse{Resp: "online"}, nil
	}else{//solo queda la posibilidad de que reciba por parte del cliente el algoritmo con el cual funcionara la maquina
		tipo_distribucion = msg.GetType()//actualiza el tipo de distribución con la cual funcionará el node
	}
	return &pb.TypeResponse{Resp: "" }, nil
}
func (s *server) RicandAgra(ctx context.Context, msg *pb2.RicandAgraRequest) (*pb2.RicandAgraResponse, error) {//funcion de los datanodes
	return &pb2.RicandAgraResponse{Resp: "mensaje" , Id: int64(1)}, nil
}
func (s *server) Status(ctx context.Context, msg *pb2.StatusRequest) (*pb2.StatusResponse, error) {//retorna multiples consultas desde los datanodes
	//0.- Retorna tipo de distribucion 1.- Informa libertad del namenode/centralizado
	switch msg.GetId(){
	case 0: //retorna que tipo de distribucion esta corriendo
		if tipo_distribucion=="centralizado"{
			return &pb2.StatusResponse{Resp: 0 }, nil //0 representa centralizado
		}else if tipo_distribucion=="distribuido"{
			return &pb2.StatusResponse{Resp: 1 }, nil //1 representa distribuido
		}else{
			return &pb2.StatusResponse{Resp: 2 }, nil //2 representa sin designacion actual
		}
	case 1: //para centralizado el namenode notificara si se encuentra en libertad u ocupado
		if estado{
			return &pb2.StatusResponse{Resp: 1 }, nil //esta libre
		}else{
			return &pb2.StatusResponse{Resp: 0 }, nil //esta ocupado con otro datanode
		}
	case 2: //data node manda esto cuando termina distribucion y libera al namenode/centralizado
		estado=true //actualiza el estado a libre
		return &pb2.StatusResponse{Resp: -1 }, nil
	default: 
		return &pb2.StatusResponse{Resp: -1 }, nil
	}
	return &pb2.StatusResponse{Resp: -1 }, nil
	
}
func InicioTipo()bool{ //funcion de inicializacion para verificar si otras maquinas ya online poseen una distribucion implementada
	initipo:= []bool{false,false}  // centralizado,distribuido
	for _,dire :=range directions{ //revisa cada node
		if dire==this_datanode{ //omite la consulta a si mismo
			continue
		}
		conn, err := grpc.Dial(dire, grpc.WithInsecure()) //genera la conexion con el node respectivo
		if err != nil {
			continue
		}
		defer conn.Close()
		client := pb2.NewNodeServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		msg:= &pb2.StatusRequest{Id: 0} //realiza la consulta de estado
		resp, err := client.Status(ctx, msg)
		if err != nil {
			continue
		}
		if resp.GetResp()==0{// 0 representa centralizado
			initipo[0]=true //actualiza el arreglo al inicio de la funcion que representa que encontro a un node en modo centralizado
		}else if resp.GetResp()==1{ // 1 representa distribuido
			initipo[1]=true //actualiza el arreglo al inicio de la funcion que representa que encontro a un node en modo distribuido
		}//2 no representa nada por lo que lo omite
	}
	if !initipo[0] && !initipo[1]{//si el arreglo queda vacio aun no hay un tipo de distribucion implementado por lo que no hace nada
		return true
	}else if initipo[0] && !initipo[1]{//si detecta solo centralizado pone al node en modo centralizado
		tipo_distribucion="centralizado"
		return true
	}else if !initipo[0] && initipo[1]{//si detecta solo distribuido pone al node en modo centralizado
		tipo_distribucion="distribuido"
		return true
	}else{//si detecta ambas notifica el error y cierra el sistema
		fmt.Println("Inconsistencia de nodos, reinicie el sistema completo por favor")
		return false
		
	}
}

func main()  {
	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatal("Error conectando: %v", err)
	}
	s := grpc.NewServer()
	in:=InicioTipo()//verifica el estado de otras maquinas para actualizar de ser necesario el modo a centralizado o distribuido
	if !in{//si hay un estado sin sentido cierra el sistema y notifica
		return
	}
	pb.RegisterClientServiceServer(s, &server{}) //recibe conexión con el camion
	pb2.RegisterNodeServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}