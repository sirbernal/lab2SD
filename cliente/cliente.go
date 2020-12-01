package main

import (
	"context"
	"fmt"
	"log"
	"time"
	//"bufio"
	"io"
	"math"
	"os"
	"strconv"
	pb "github.com/sirbernal/lab2SD/proto/client_service"
	"google.golang.org/grpc"
)
var datanode = []string{"localhost:50052","localhost:50053","localhost:50054"} //arreglo que contiene las direcciones de los datanodes
var directions = []string{"localhost:50055", "localhost:50052","localhost:50053","localhost:50054"} //arreglo que incluye al namenode (se agrego despues de muchas funciones ya hechas)
var nodemode = []string{"","","",""} //arreglo que guarda el tipo de distribucion de cada nodo
var nodestatus = []bool{false,false,false,false} //arreglo que guarda el estado de un nodo= true:conectado false:desconectado
var chunks [][]byte //donde se guarda los chunks para subir
var rechunks [][]byte //donde se guarda los chunks para bajar
var tipo_distribucion string  //se guarda el tipo de distribucion del cliente
type ChunkAndN struct{  //estructura usada para hacer chunks
	Chunk []byte //chunk
	N int //numero necesario para envio de chunks visto en este tutorial https://ops.tips/blog/sending-files-via-grpc/
}


func Chunker(archivo string) []ChunkAndN{ //funcion que descomprime el archivo en chunks
	file, err := os.Open(archivo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} 
	defer file.Close()
	var chunks []ChunkAndN //arreglo que guardara los chunks y N
	fileInfo, _ := file.Stat() //sacado practicamente todo del tutorial del enunciado
	var fileSize int64 = fileInfo.Size()
	const fileChunk = 256000//tamaño solicitado en enunciado
	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
	for i := uint64(0); i < totalPartsNum; i++ {
		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)
		n,err :=file.Read(partBuffer)
		if err != nil {
			if err == io.EOF {
				err = nil
				continue
			}
			return nil
		}
		chunks = append(chunks,ChunkAndN{partBuffer,n})//aqui se almacena el chunk en el arreglo
	}
	file.Close()
	return chunks//retorna el arreglo
}
func Unchunker(name string){ //funcion que vuelve el conjunto de chunks guardados en "rechunk" en el archivo original
	_, err := os.Create(name)//funcion tambien sacada del tutorial del enunciado
	if err != nil {
			fmt.Println(err)
			os.Exit(1)
	}
	file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _,chunk:= range rechunks{
		_, err := file.Write(chunk)	
		if err != nil {
				fmt.Println(err)
				os.Exit(1)
		}
		file.Sync()
	}
	file.Close()
	rechunks = [][]byte{} //limpia memoria de los chunks descargados
}

func SolicitarLibros()[]string{//solicita la lista de libros al namenode
	conn, err := grpc.Dial(directions[0], grpc.WithInsecure()) //establece conexion con el namenode
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	client := pb.NewClientServiceClient(conn)    
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	msg:= &pb.DownloadNamesRequest{Req: "Libros",} //envia el request con la cadena "Libros" la cual el namenode detecta para esta petición
	resp, err := client.DownloadNames(ctx, msg)	
	fmt.Println("Archivos(s) disponible(s): "+strconv.Itoa(len(resp.GetNames())))//imprime el total de archivos disponibles para descargar
	for i,name:=range resp.GetNames(){
		fmt.Println(strconv.Itoa(i+1)+".-"+name) //imprime la lista de archivos con su numero asociado
	}
	return resp.GetNames()//retorna el arreglo de los titulos por si se necesita descargar
}
func SolicitarUbicacion(libro string)([]int64){ //solicita la ubicacion de los chunks de un archivo especifico
	conn, err := grpc.Dial(directions[0], grpc.WithInsecure())//establece conexion con el namenode
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	client := pb.NewClientServiceClient(conn)    
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	msg:= &pb.LoCRequest{Req: libro,}//solicita el archivo que requiere por medio de su nombre
	resp, err := client.LocationsofChunks(ctx, msg)
	return resp.GetLocation()//retorna la ubicacion de cada chunk disponible en los datanodes
		
}
func DescargarChunks(name string, prop []int64){//funcion que descarga y almacena los chunks acorde al orden que se solicito en la funcion anterior
	for i, node:= range prop{
		datan:=datanode[int(node)] //selecciona el datanode a la cual deberia poseer el chunk nro i del archivo a rearmar
		chunkadescargar:=name+"_"+strconv.Itoa(i) //genera el nombre del chunk que deberia solicitar
		conn, err := grpc.Dial(datan, grpc.WithInsecure()) //hace la conexion con el datanode respectivo
		if err != nil {
			log.Fatalln(err)
		}
		defer conn.Close()
		client := pb.NewClientServiceClient(conn)		
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		msg:= &pb.DownloadChunksRequest{Name: chunkadescargar} //envia el nombre del chunk a descargar
		resp, err := client.DownloadChunks(ctx, msg)
		rechunks=append(rechunks,resp.GetChunk())//lo guarda en el arreglo de chunks para su rearmado
	}
	Unchunker(name)//llama a la funcion que rearma el archivo en base a los chunks
}
func SubirArchivo(node int, archivo string)(){ //funcion que envia el archivo al datanode seleccionado en forma de chunks
	conn, err := grpc.Dial(datanode[node], grpc.WithInsecure()) //genera la conexion con el datanode de destino
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	client := pb.NewClientServiceClient(conn)    
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	chunks := Chunker("./"+archivo) //divide en chunks el archivo seleccionado
	msg:= &pb.UploadRequest{Tipo: 0, Nombre: archivo, Totalchunks: int64(len(chunks))} //envia el nombre del archivo con el total de chunks para generar la prop en datanode
	resp, err := client.Upload(ctx, msg)
	if resp.GetResp()==int64(0){//si la propuesta es aceptada, manda los chunks al datanode para su distribución
		for _,chunk :=range chunks{
			msg:= &pb.UploadChunksRequest{Chunk: chunk.Chunk[:chunk.N]}
			_, err := client.UploadChunks(ctx, msg)
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
		}	
	}
}
func VerifNodos(){ //funcion que revisa y actualiza el estado de los datanodes y el namenode
	for i, dire := range directions{ 
		conn, err := grpc.Dial(dire, grpc.WithInsecure()) //genera la conexion con el node
		if err != nil {
			continue
		}
		defer conn.Close()

		client := pb.NewClientServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		msg:= &pb.TypeRequest{Type: "status" } //envia la consulta por medio de la palabra "status"

		resp, err := client.TypeDis(ctx, msg)
		if err != nil {
			continue
			}
		if resp.GetResp()=="online"{ //si el node le responde "online" actualiza el estado del nodo
			nodestatus[i]=true
		}	
	}
}

func menu2(tipo string){//Menu que se despliega luego de seleccionar el algoritmo de exclusion
	var menu2 string
	Menu2: //menu principal del algoritmo
		for {
			fmt.Print("\n\n---Menu Algoritmo de Exclusión Mutua "+tipo+"--- \nIngrese opción\n1.-Subir Archivo\n2.-Ver archivos en sistema\n3.-Cerrar Sistema\nIngrese opción:")
			VerifNodos()//actualiza el estado de conexion de los nodes
			_,err:=fmt.Scanln(&menu2)
				if err!=nil{
					fmt.Print("\nFormato de ingreso no válido, pruebe nuevamente...")
					continue
				}
			switch menu2{//dependiendo de la opcion pasaremos a los submenus respectivos
			case "1": //menu para subir un archivo
				Menusubida:
					for{
						VerifNodos()//actualiza el estado de los nodos
						if !nodestatus[0]{//si el namenode esta apagado evita la poisibilidad de subir un archivo ya que no quedara registro
							fmt.Println("Namenode apagado, por favor inicializar... regresando al menu principal")
							continue Menu2//retorna a menu principal
							
						}
						var datasubida string
						var dataint int
						fmt.Println("--Lista de Datanodes--")
						for i,j:= range datanode{//imprime los datanodes con su ip y ultimo estado de conexion
							var k string
							if nodestatus[i+1]{
								k="  online"
							}else{
								k="  offline"
							}
							fmt.Println(strconv.Itoa(i+1)+".- Datanode "+strconv.Itoa(i+1)+" Ip:"+j+k)
						}
						for {
							fmt.Print("n.- Volver\nSeleccione un datanode:") // enviar n permite volver al menu principal
							_,err:=fmt.Scanln(&datasubida)
							if err!=nil{
								fmt.Println("Formato no válido, pruebe nuevamente...")
								continue
							}
							if datasubida=="n"{
								continue Menu2
							}
							dataint,err=strconv.Atoi(datasubida) //transforma el input en int
							if err!=nil{
								fmt.Println("Formato no válido, pruebe nuevamente...")
								continue
							}
							if !nodestatus[dataint]{ //evita envio de archivos si el nodo esta desconectado y reinicia el menu de subida
								fmt.Println("Nodo fuera de linea, actualizando información de nodos...")
								continue Menusubida
							}
							if dataint>0 && dataint<len(datanode)+1{ //si la opcion es valida permite seguir con el sistema de subida
								break
							}else{
								fmt.Println("Opción no válida, pruebe nuevamente...")
							}
					}
					var archivo string
					fmt.Print("Ingrese nombre de archivo con su extension(ejemplo: archivo.pdf):") 
					fmt.Scanln(&archivo) //guarda el nombre de su archivo incluida la extension
					SubirArchivo(dataint-1,archivo) //llama a la funcion de subir archivo con el nodo seleccionado
					continue Menu2 //vuelve al menu principal
					} 
			case "2":
				VerifNodos() //actualiza el estado de los nodes
				if !nodestatus[0]{ //si el namenode esta apagado aborta el intento de solicitar los archivos
					fmt.Println("Namenode apagado, por favor inicializar... regresando al menu principal")
					continue Menu2
					
				}
				libros:=SolicitarLibros()//solicita la lista de archivos guardados por el namenode
				if len(libros)==0{//si no hay archivos retorna al menu principal
					fmt.Println("\n\nSistema sin archivos disponibles\n")
					continue Menu2
				}
				//si hay archivos estos se imprimiran con un indice respecivo para ser seleccionados
				var opt string
				var optint int
				for{
					fmt.Print("Seleccione un archivo para descargar o escriba 'n' para volver: ")
					_,err:=fmt.Scanln(&opt)
					if err!=nil{
						fmt.Println("Formato no válido, pruebe nuevamente...")
						continue
					}
					if opt=="n"{//permite volver al menu principal por medio del ingreso de "n"
						continue Menu2
					}
					optint,err=strconv.Atoi(opt)
					if err!=nil{
						fmt.Println("Formato no válido, pruebe nuevamente...")
						continue
					}
					if optint>0 && optint<len(libros)+1{//verifica que el numero se encuentre dentro de las opciones disponibles
						break
					}else{
						fmt.Println("Opción no válida, pruebe nuevamente...")
					}
				}
				fmt.Print("Descargando archivo: "+libros[optint-1])//notifica que archivo se selecciono
				ubicacion := SolicitarUbicacion(libros[optint-1]) //solicita al namenode donde se encuentran los chunks del archivo seleccionado
				if len(ubicacion)==0{ //puede recibir una ubicacion nula (array vacio) lo cual detecta que al menos un nodo que contenia el archivo fue apagado y evita la descarga
					fmt.Println("\n\nArchivo no recuperable\nMotivo: nodos que contenían parte de este estan fuera de línea\nRedireccionando a Menú Principal...\n")
					continue Menu2 //retorna al menu principal
				}
				DescargarChunks(libros[optint-1], ubicacion) //descarga los chunks desde los datanodes respectivos y rearma el archivo respectivo
			case "3": //Permite cerrar el cliente terminando su ejecución
				fmt.Println("Terminando ejecución cliente...")
				break Menu2 
			default:
				fmt.Println("\nFormato u opción no válida, pruebe nuevamente:\n\n")
				continue Menu2
			}
		}
}
func DistribStatus()int{ //Verifica el estado que tienen los nodes si este inicio despues de ellos u otro cliente los designo antes
	cent:=0 //cuenta los nodes que estan en modo centralizado
	dist:=0 //cuenta los nodes que estan en modo distribuido
	for _,stat:=range nodemode{
		if stat=="centralizado"{
			cent++
		}else if stat=="distribuido"{
			dist++
		}
	}
	if dist==0 && cent==0{//si no hay ninguno significa que no hay aun un algoritmo seleccionado
		return 0
	}
	if cent>0 && dist==0{//detecta que al menos un node esta en formato centralizado
		return 1
	}
	if dist>0 && cent==0{//detecta que al menos un node esta en formato distribuido
		return 2
	}
	if dist>0 && cent>0{//detecta que hay al menos un node de cada tipo y notifica de la incongruencia del sistema
		return 3
	}
	return -1 //return que golang exige
}
func VerifInicial()int{ //funcion que verifica inicialmente que algortmo esta corriendo cada node
	for i, dire := range directions{
		conn, err := grpc.Dial(dire, grpc.WithInsecure())//genera la conexion con el node
		if err != nil {
			continue
		}
		defer conn.Close()
		client := pb.NewClientServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		msg:= &pb.TypeRequest{Type: "inicio" } //envia la palabra "inicio" al nodo
		resp, err := client.TypeDis(ctx, msg)
		if err != nil {
			continue
			}
		nodemode[i]=resp.GetResp() //actualiza el tipo de distribucion de cada node
	}
	return DistribStatus() //calcula que estado necesita tener el cliente al iniciar
}
func menu(){//menu de inicializacion del cliente
	init:=VerifInicial() //verificara el estado de los nodos para ver si es posible omitir el menu de inicio
	switch init{
	case 0: //si no hay ningun algoritmo definido en ningun nodo procede a consultar con el menu de inicio
		var menu1 string
		for {
			fmt.Print("\n\n---Menu Principal Cliente--- \nIngrese que algoritmo de exclusión mutua desea experimentar\n1.-Algoritmo de Exclusión Mutua Centralizada\n2.-Algoritmo de Exclusión Mutua Distribuida\nIngrese opción:")
			_,err:=fmt.Scanln(&menu1)
				if err!=nil{
					fmt.Print("\nFormato de ingreso no válido, pruebe nuevamente:")
					continue
				}
			if menu1=="1"{//se selecciona el modo centralizado
				tipo_distribucion = "centralizado" //actualiza el estado de cliente a centralizado
				
				for _, dire := range directions{ //envia el estado a todos los nodes deben estar en modo centralizado
					conn, err := grpc.Dial(dire, grpc.WithInsecure())
					if err != nil {
						continue
					}
					defer conn.Close()
					client := pb.NewClientServiceClient(conn)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()
					msg:= &pb.TypeRequest{Type: tipo_distribucion } //envia "centralizado" a cada nodo para actualizar su modo
					_, err = client.TypeDis(ctx, msg)
					if err != nil {
						continue
						}
				}
				menu2("Centralizada")//redirecciona al menu principal en modo centralizado
				break
			}else if menu1=="2"{//se selecciona el modo distribuido
				tipo_distribucion = "distribuido"
				for _, dire := range directions{ //envia el estado a todos los nodes deben estar en modo distribuido
					conn, err := grpc.Dial(dire, grpc.WithInsecure())
					if err != nil {
						continue
					}
					defer conn.Close()
					client := pb.NewClientServiceClient(conn)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()
					msg:= &pb.TypeRequest{Type: tipo_distribucion } //envia "distribuido" a cada nodo para actualizar su modo
					_, err = client.TypeDis(ctx, msg)
					if err != nil {
						continue
						}
				}
				menu2("Distribuida")//redirecciona al menu principal en modo distribuido
				break
			}else{
				fmt.Println("\nFormato u opción no válida, pruebe nuevamente:\n\n")
				continue
			}
		}
	case 1: //detecta que todos los nodes online estan en modo centralizado
		fmt.Println("Algoritmo de Exclusión Mutua Centralizada Detectado\nRedireccionando a Menu del algoritmo...")
		menu2("Centralizada")
	case 2: //detecta que todos los nodes online estan en modo distribuido
		fmt.Println("Algoritmo de Exclusión Mutua Distribuida Detectado\nRedireccionando a Menu del algoritmo...")
		menu2("Distribuida")
	case 3: //detecta que los nodos online tienen modos distintos entre ellos por lo cual notificara que el sistema no funcionara
		fmt.Println("Algoritmo corrompido detectado! (Presencia de distintos algoritmos en distintos nodos)\nPor favor, reiniciar el sistema completo\n\n---Cerrando Sistema--")
		return
	}
	
}
func main() {
	menu()//inicializa el menu
}
