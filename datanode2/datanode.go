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
var ocupado = false //estado para distribuido
var id_node = 2 //id para distribuido 
var datanode = []string{"10.10.28.82:50052","10.10.28.83:50053","10.10.28.84:50054"}  //arreglo que contiene las direcciones de los datanodes
var directions = []string{"10.10.28.81:50055", "10.10.28.82:50052","10.10.28.83:50053","10.10.28.84:50054"} //arreglo que incluye al namenode (se agrego despues de muchas funciones ya hechas)
var total int64 //guarda el numero total de chunks de una propuesta
var nombrearchivo string //guarda el nombre del archivo a subir
var cont int //indica el total de conectados
var chunks [][]byte //guarda los chunks enviados desde el cliente
var tipo_distribucion string //se guarda el tipo de distribucion del datanode
var this_datanode = datanode[1] //direccion de datanode
var activos []int //guarda cuales son los nodos activos
var datanodestatus = []bool{false,false,false} //guarda el estado de los datanodes en orden
var c_mensajes = 0 // contador de mensajes usado para pruebas que se haran en el informe
var timeout = time.Duration(1)*time.Second //timeout para conexiones 

func SaveChunk (chunk []byte, name string){ //transforma el chunk en un archivo para guardar en almacenamiento acorde al nombre entregado
	fmt.Println("Guardando chunk en memoria\nNombre del Chunk: "+name+"\n")
	fileName := name //sacado del tutorial del enunciado
	_, err := os.Create(fileName)
	if err != nil {
			fmt.Println(err)
			os.Exit(1)
	}
	ioutil.WriteFile(fileName, chunk, os.ModeAppend)
}
func SearchChunk (name string) (chunk []byte){ //busca el chunk solicitado en almacenamiento
	fmt.Println("Buscando chunk en almacenamiento\nNombre: "+name)
	newFileChunk, err := os.Open(name) //sacado del tutorial del enunciado
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
	reader := bufio.NewReader(newFileChunk) //traduce el archivo a un arreglo de bytes para ser enviado
	reader.Read(chunkBufferBytes)
	fmt.Println("Enviando Chunk\n")
	return chunkBufferBytes //retorna el chunk en forma de bytes
}
func GenerarPropuesta (total int)([]int64){ //funcion que genera la propuesta inicial de los 3 datanodes conectados
	fmt.Println("Generando Propuesta Inicial...")
	var propuesta []int64
	for i:=0;i<=total/3;i++{ //genera la propuesta para cada chunk en base al total/3, es decir por cada 3 chunks los distribuye entre los 3 datanodes
		rand.Seed(time.Now().UnixNano())//genera una semilla random basada en el time de la maquina
		lilprop:=rand.Perm(3)//genera un arreglo al azar con los tres datanodes, como por ejemplo (0,2,1),(1,0,2), entre otros
		if i==total/3{ //detecta la ultima serie de chunks
			sobra:=total%3
			for j,num :=range lilprop{
				if j==sobra{
					break//si ya no hay mas chunks termina la generacion de propuesta
				}
				propuesta=append(propuesta,int64(num)) //agrega el nodo al cual correspondera la propuesta
			}
		}else{
			for _,num :=range lilprop{
				propuesta=append(propuesta,int64(num)) //agrega el nodo al cual correspondera la propuesta
			}
		}
	}
	return propuesta //retorna la propuesta inicial
}

func (s *server) Upload(ctx context.Context, msg *pb.UploadRequest) (*pb.UploadResponse, error) {//funcion que recibe el nombre y cantidad de chunks que tiene el archivo del cliente
	fmt.Println("Recibiendo Archivo de cliente...\nNombre de Archivo:"+msg.GetNombre()+"\nCantidad de Partes:"+strconv.Itoa(int(msg.GetTotalchunks()))+"\nProcediendo a recibir partes...")
	nombrearchivo=msg.GetNombre() //guarda el nombre
	total=msg.GetTotalchunks() //guarda el total de chunks
	//c_mensajes = c_mensajes + 2
	return &pb.UploadResponse{Resp : int64(0), }, nil
}

func RicartyAgrawala()bool{
	
	for _,dire:= range datanode{
		if this_datanode == dire{ // Solo contactaremos con nodos distintos al nuestro
			continue
		}
		conn, err := grpc.Dial(dire, grpc.WithInsecure()) // enviamos propuesta a un nodo especifico
		if err != nil {
			continue
		}
		defer conn.Close()

		client := pb2.NewNodeServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		msg:= &pb2.RicandAgraRequest{Id: int64(id_node)} // generamos el mensaje con la propuesta

		resp, err := client.RicandAgra(ctx, msg) // enviamos la propuesta y recibimos la respuesta
		//c_mensajes = c_mensajes + 2
		if resp.GetResp() != "mensaje"{
			if resp.GetId() > int64(id_node){
				return false
			}
		}

	}
	return true
}


func (s *server) UploadChunks(ctx context.Context, msg *pb.UploadChunksRequest) (*pb.UploadChunksResponse, error) { //funcion que guarda los chunks y cuando llegan todos genera la distribucion
	//mt.Println(len(chunks))
	//fmt.Println("Recibido")
	
	if tipo_distribucion == "centralizado"{  // Logica cuando el algoritmo de distrib. es centralizado
		
		chunks=append(chunks,msg.GetChunk()) // Agregamos el chunk a nuestro arreglo

		if int64(len(chunks))==total{ // Cuando llegan todos los chunks del archivo al datanode, realizamos propuesta
			fmt.Println("Recepción Completa de partes...")
			for{
				fmt.Println("Verificando estado de Namenode...")
				conn, err := grpc.Dial(directions[0], grpc.WithInsecure())//nos conectamos al datanode a ver su estado
				if err != nil {
					log.Fatalln(err)
				}
				defer conn.Close()
				client := pb2.NewNodeServiceClient(conn)

				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel() // Con esta funcion generaremos la propuesta
				msg:= &pb2.StatusRequest{Id: 1} // generamos el mensaje con la consulta de estado al namenode
				resp, err := client.Status(ctx, msg)
				//c_mensajes = c_mensajes + 2
				if resp.GetResp()==1{ //si esta libre, procedemos a seguir el algoritmo
					fmt.Println("Namenode disponible...")
					break
				}
				fmt.Println("Namenode ocupado... Reiniciando solicitud...\n")
				//si no, volvemos a consultar si se encuentra disponible el namenode
			}	
				/* GENERAR PROPUESTA*/
			/* Primero, generamos la conexion con namenode*/
			conn, err := grpc.Dial(directions[0], grpc.WithInsecure())
			if err != nil {
				log.Fatalln(err)
			}
			defer conn.Close()
			client := pb2.NewNodeServiceClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(),timeout)
			defer cancel()
			propuesta := GenerarPropuesta(int(total)) // Con esta funcion generaremos la propuesta
			fmt.Println("Enviando propuesta a Namenode")		
			msg:= &pb2.PropuestaRequest{Prop: propuesta, Name: nombrearchivo} // generamos el mensaje con la propuesta
			resp, err := client.Propuesta(ctx, msg) // enviamos la propuesta y recibimos la respuesta
			//c_mensajes = c_mensajes + 2
			if resp.GetMsg() == false{  // cuando se rechaza la propuesta, la actualizamos con la propuesta recibida x namenode
				fmt.Println("Propuesta Rechazada... Actualizando Propuesta")
				propuesta = resp.GetProp()
			}else{
				fmt.Println("Propuesta Aceptada")
			}
			msg2:= &pb2.StatusRequest{Id: 2} // liberamos al datanode, ya que se reviso la propuesta y se procedera a distriburi
			fmt.Println("Liberando Namenode...")
			_, err=client.Status(ctx, msg2)
			//c_mensajes = c_mensajes + 2
			if err != nil {
				log.Fatalln(err)
			}
			/* GENERAR DISTRIBUCION*/
			/* Leemos el arreglo de propuesta que tiene las designaciones de cada chunk que ira a cada datanode*/
			fmt.Println("Iniciando Distribucion de Chunks...")
			for i,j :=range propuesta{
				if this_datanode==datanode[j]{ //si el chunk corresponde al datanode local, simplemente se guarda
					fmt.Println("Guardando Parte "+strconv.Itoa(i)+" en datanode actual")
					SaveChunk(chunks[i],nombrearchivo+"_"+strconv.Itoa(i)) // se guarda en almacenamiento el chunk
					continue
				}
				// Procedemos a generar la conexion con el datanode a donde enviaremos el chunk
				conn, err := grpc.Dial(datanode[j], grpc.WithInsecure())
				if err != nil { //desconectar un nodo destino durante la distribucion producira un error y abortará todo lo anterior
					fmt.Println("Proceso abortado, se ha desconectado el nodo durante la distribucion")
					break
				}
				defer conn.Close()
				client := pb2.NewNodeServiceClient(conn)
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				fmt.Println("Enviando Parte "+strconv.Itoa(i)+" a datanode "+strconv.Itoa(int(j)+1)+"\n")
				msg:= &pb2.DistribucionRequest{Chunk: chunks[i], Name: nombrearchivo+"_"+strconv.Itoa(i)}// Creamos el mensaje con el chunk correspondiente + el nombre del chunk 
				_, err = client.Distribucion(ctx, msg)// Enviamos el chunk
				//c_mensajes = c_mensajes + 2
				if err != nil {
					continue
				}
			}
			fmt.Println("Distribucion Completada\n\n")
			chunks = [][]byte{}// Una vez realizada la distribucion, reiniciamos el arreglo donde guardamos los chunks
			//c_mensajes = c_mensajes + 2
			//mensaj := "El servidor guardo el archivo, num. de mensajes totales: " + strconv.Itoa(c_mensajes)
			//c_mensajes = 0
			return &pb.UploadChunksResponse{Resp : "subida completa" }, nil //Avisamos al cliente de que esta completa la subida y distribucion del archivo
			}
		//c_mensajes = c_mensajes + 2
		return &pb.UploadChunksResponse{Resp : "chunk recibido en el server", }, nil// Esto se envia cuando no se tienen todos los chunks del archivo

		
	}else if tipo_distribucion == "distribuido" {   // Logica cuando el algoritmo de distrib. es distribuido
		chunks=append(chunks,msg.GetChunk()) // Agregamos el chunk a nuestro arreglo
		if int64(len(chunks))==total{// Cuando llegan todos los chunks del archivo al datanode, realizamos propuesta
			fmt.Println("Recepción Completa de partes...")
			/* Logica de esto, es que si una propuesta se rechaza, pues hay que volver a realizar una
			nueva propuesta, por lo que rompemos el segundo for que esta hecho para enviar la propuesta a cada
			datanode y volvemos a realizarlo con una nueva propuesta, para eso usamos una variable booleana
			proceso que cuando se aceptan todas las propuestas, recien procederiamos a distribuir */
			var propuesta []int64   // Aqui almacenamos nuestra propuesta
			propuesta = GenerarPropuesta(int(total))// Con esta funcion generaremos la propuesta
			var contador int  // Usamos este contador para comparar la cantidad de envios con respecto a la cantidad de nodos conectados
			/* Lo primero que haremos es un loop que nos servira para la generacion (y posible regeneracion) de propuesta
			de distribucion de chunks, puesto que si se rechaza una propuesta hay que generar una nueva y volver a enviar a los
			demas datanodes esta nueva propuesta, por eso es que necesitabamos agregar un nuevo for mas (el "Proceso")
			*/
			Proceso:
			for{
				/* Esta funcion AllAlive la llamamos  por que ademas de poder comprobar si los nodos que implican la 
				propuesta estan disponibles, actualiza el arreglo datanodestatus que nos indica los datanodes que estan vivos,
				basicamente para esto la llamamos realmente aca
				*/
				AllAlive([]int64{})
				fmt.Println("Enviando Propuesta a datanodes")  
				/* Primero, generamos la conexion con cada datanode*/
				for i,dire:= range datanode{
					
					if this_datanode == dire{ // No enviaremos a este mismo nodo la propuesta a generar
						contador++
						continue
					}
					if !datanodestatus[i]{ // SI el nodo esta desconectado, no haremos envio a aquel nodo
						continue
					}
					
					conn, err := grpc.Dial(dire, grpc.WithInsecure()) // generamos conexion a un nodo especifico
					if err != nil {
						continue
					}
					defer conn.Close()

					client := pb2.NewNodeServiceClient(conn)

					ctx, cancel := context.WithTimeout(context.Background(), timeout)
					defer cancel()

					msg:= &pb2.PropuestaRequest{Prop: propuesta, Name: nombrearchivo} // generamos el mensaje con la propuesta 

					resp, err := client.Propuesta(ctx, msg) // enviamos la propuesta y recibimos la respuesta
					//c_mensajes = c_mensajes + 2
					//estado := resp.GetMsg()
					//fmt.Println(resp.GetProp(),resp.GetMsg(),contador,dire)
					if resp.GetMsg() == false{  // cuando se rechaza la propuesta, la actualizamos con la propuesta recibida x namenode	
						fmt.Println("Propuesta rechazada...")
						propuesta=GenerarPropuestaNueva(len(propuesta),TotalConectados()) // Generamos una nueva propuesta
						contador = 0 // reiniciamos el contador
						break
					}else{
						contador++ //Sumamos de que se logro un envio
					}
				}
				// Si notamos que enviamos todos los mensajes a todos los nodos que estan conectados, pues estamos listos para el siguiente
				// paso, por lo cual rompemos este Proceso
				if int(contador)==TotalConectados(){ 
					break Proceso
				}
				contador=0
			}
			//Enviar mensaje a namenode con la propuesta aceptada por nodos 
			fmt.Println("Propuesta Aceptada")
		
			
			//actualizamos esta variable que indica que estaremos contactando con namenode
			ocupado = true
			fmt.Println("Actualizando estado: OCUPADO")
			// Se usa el algoritmo de Ricart y Agrawala para pedir permisos de acceso a namenode
			
			for !RicartyAgrawala(){} 
			/* Si se tiene aprobacion de los demas nodos para contactar namenode, procedera a contactarlo
			en caso contrario, estara consultando constantemente a la autorizacion de los demas nodos*/
			fmt.Println("Enviando Propuesta a Namenode...")
			// Contactamos con namenode 
			conn, err := grpc.Dial(directions[0], grpc.WithInsecure())
			if err != nil {
				log.Fatalln(err)
			}
			defer conn.Close()

			client := pb2.NewNodeServiceClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			msg:= &pb2.PropuestaRequest{Prop: propuesta, Name: nombrearchivo} // enviamos propuesta a namenode para que la escriba

			_, err = client.Propuesta(ctx, msg) // enviamos la propuesta y recibimos la respuesta
			//c_mensajes = c_mensajes + 2
			ocupado = false // como ya terminamos de contactarnos con namenode, desactivamos esta variable
			fmt.Println("Actualizando estado: LIBRE")
			/* GENERAR DISTRIBUCION*/
			/* Leemos el arrigo de propuesta que tiene las designaciones de cada chunk que ira a cada datanode*/
			fmt.Println("Iniciando Distribucion de Chunks")
			for i,j :=range propuesta{
				if this_datanode==datanode[j]{ //claramente, si un chunk se debe quedar en este datanode, para que enviarlo XD
					fmt.Println("Guardando Parte "+strconv.Itoa(i)+" en datanode actual")
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
		
				ctx, cancel := context.WithTimeout(context.Background(),timeout)
				defer cancel()
				// Creamos el mensaje con el chunk correspondiente + el nombre del chunk 
				fmt.Println("Enviando Parte "+strconv.Itoa(i)+" a datanode "+strconv.Itoa(int(j)+1)+"\n")
				msg:= &pb2.DistribucionRequest{Chunk: chunks[i], Name: nombrearchivo+"_"+strconv.Itoa(i)}
				
				// Enviamos el chunk
				_, err = client.Distribucion(ctx, msg)
				//c_mensajes = c_mensajes + 2
				if err != nil {
					continue
				}
			}
			fmt.Println("Distribucion Completada\n\n")
			chunks = [][]byte{}

			//c_mensajes = c_mensajes + 2
			//mensaj := "El servidor guardo el archivo, num. de mensajes totales: " + strconv.Itoa(c_mensajes)
			//c_mensajes = 0
			return &pb.UploadChunksResponse{Resp : "ok", }, nil //Avisamos de que estamos ok con la subida
		}
		// Este mensaje se envia cuando no hemos recibido todos los chunks de cliente
		//c_mensajes = c_mensajes + 2
		return &pb.UploadChunksResponse{Resp : "chunk recibido en el server", }, nil 


	} 
	// En el caso de que el sistema no tenga claro el tipo de distribucion (la variable tipo_distribucion no tiene valor)
	// entonces significa que algo esta fallando, avisamos de eso al cliente
	//c_mensajes = c_mensajes + 2
	return &pb.UploadChunksResponse{Resp : "Fallo algo", }, nil
}

func (s *server) Alive(ctx context.Context, msg *pb2.AliveRequest) (*pb2.AliveResponse, error) {//responde que se encuentra en linea si le es solicitado
	//c_mensajes = c_mensajes + 2
	//fmt.Println("Cant. mensajes realizados aca: ",  c_mensajes)
	//c_mensajes = 0
	return &pb2.AliveResponse{Msg : "Im Alive, datanode", }, nil
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
	fmt.Println("Generando nueva propuesta...")
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
func VerifProp(prop []int64)bool{ //funcion que verifica que en la propuesta generada esten todos los nodos activos
	for _,chunk:=range prop{
		if !datanodestatus[chunk]{
			return false //si dentro de la propuesta encuentra un nodo apagado designado retorna falso
		}
	}
	return true 
}

func AllAlive (prop []int64) (bool){//funcion que verifica los datanodes conectados, actualiza en memoria cuales lo estan y retorna falso si hay discrepancia con la lista de conectados enviada
	for j,dire :=range datanode{
		conn, err := grpc.Dial(dire, grpc.WithInsecure()) //inicia conexion con cada datanode
		if err != nil {
			datanodestatus[j]=false
			continue
		}
		defer conn.Close()
		client := pb2.NewNodeServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		msg:= &pb2.AliveRequest{Msg: "Are u alive?"} //realiza la consulta al datanode
		_, err = client.Alive(ctx, msg)
		//c_mensajes = c_mensajes + 2
		if err != nil {
			datanodestatus[j]=false
			continue
		}
		datanodestatus[j]=true //actualiza la lista de estados
	}
	if reflect.DeepEqual(prop, []int64{}){//utilizado para actualizar estado de nodes
		return true
	}
	if !VerifProp(prop){//si dentro de la propuesta hay un nodo desconectado retorna falso
		fmt.Println("Error!...Propuesta no valida... Datanode desconectado detectado en propuesta")
		return false
	}
	fmt.Println("Propuesta valida")
	return true //retorna true si todo es valido
}

func (s *server) Propuesta(ctx context.Context, msg *pb2.PropuestaRequest) (*pb2.PropuestaResponse, error) {//recibe la propuesta de otro datanode en distribuido
	fmt.Println("Propuesta Recibida...Inicializando Proceso de verificación de Propuesta...")
	if AllAlive(msg.GetProp()) { //revisa que la propuesta sea valida en base a los conectados
		fmt.Println("Aprobando Propuesta")
		return &pb2.PropuestaResponse{Msg : true, Prop : []int64{}}, nil //envia true si lo es
	} else {		
		fmt.Println("Rechazando Propuesta")
		return &pb2.PropuestaResponse{Msg : false, Prop : []int64{}}, nil //envia falso si la propuesta esta mala
	}
}


func (s *server) Distribucion(ctx context.Context, msg *pb2.DistribucionRequest) (*pb2.DistribucionResponse, error) {//recibe el chunk a guardar en almacenamiento
	fmt.Println("Chunk Recibido!")
	SaveChunk(msg.GetChunk(),msg.GetName()) //guarda el chunk en almacenamiento con el nombre respectivo enviado desde otro datanode
	return &pb2.DistribucionResponse{Resp : "",}, nil
}

func (s *server) DownloadNames(ctx context.Context, msg *pb.DownloadNamesRequest) (*pb.DownloadNamesResponse, error) {//funcion del namenode
	return &pb.DownloadNamesResponse{Names : []string{} }, nil
}

func (s *server) DownloadChunks(ctx context.Context, msg *pb.DownloadChunksRequest) (*pb.DownloadChunksResponse, error) {//funcion que busca y envia el chunk solicitado
	return &pb.DownloadChunksResponse{Chunk : SearchChunk(msg.GetName()) }, nil //envia el chunk en forma de []byte encontrado en almacenamuiento
}

func (s *server) LocationsofChunks(ctx context.Context, msg *pb.LoCRequest) (*pb.LoCResponse, error) {//funcion del namenode
	return &pb.LoCResponse{Location: []int64{} }, nil
}

func (s *server) TypeDis(ctx context.Context, msg *pb.TypeRequest) (*pb.TypeResponse, error) {//funcion para multiples usos
	if msg.GetType()=="inicio"{//si el cliente o un nodo esta iniciando y este esta corriendo retornara el tipo de distribucion que actualmente posee
		return &pb.TypeResponse{Resp: tipo_distribucion}, nil
	}else if msg.GetType()=="status"{ //funcion que notifica que se encuentra online el nodo
		return &pb.TypeResponse{Resp: "online"}, nil
	}else{//solo queda la posibilidad de que reciba por parte del cliente el algoritmo con el cual funcionara la maquina
		fmt.Println("Recepcion de algoritmo de distribucion "+msg.GetType()+"\nAplicando Configuración...\nListo\n")
		tipo_distribucion = msg.GetType()//actualiza el tipo de distribución con la cual funcionará el node
	}
	return &pb.TypeResponse{Resp: "" }, nil
}

func (s *server) RicandAgra(ctx context.Context, msg *pb2.RicandAgraRequest) (*pb2.RicandAgraResponse, error) {
	if ocupado == false {
		return &pb2.RicandAgraResponse{Resp: "mensaje" , Id: int64(id_node)}, nil
	} else {
		return &pb2.RicandAgraResponse{Resp: "ocupado" , Id: int64(id_node)}, nil
	}
	
}

func (s *server) Status(ctx context.Context, msg *pb2.StatusRequest) (*pb2.StatusResponse, error) {//retorna multiples consultas desde los datanodes
	switch msg.GetId(){
		case 0: //retorna que tipo de distribucion esta corriendo
		if tipo_distribucion=="centralizado"{
			return &pb2.StatusResponse{Resp: 0 }, nil //0 representa centralizado
		}else if tipo_distribucion=="distribuido"{
			return &pb2.StatusResponse{Resp: 1 }, nil //1 representa distribuido
		}else{
			return &pb2.StatusResponse{Resp: 2 }, nil //2 representa sin designacion actual
		}
	default: 
		return &pb2.StatusResponse{Resp: -1 }, nil
	}
	return &pb2.StatusResponse{Resp: -1 }, nil
	
}
func InicioTipo()bool{ //funcion de inicializacion para verificar si otras maquinas ya online poseen una distribucion implementada
	fmt.Println("Revisando datanodes y namenode...")
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
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
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
		fmt.Println("Algoritmo de Exclusión Mutua no detectado...\nEsperando configuracion de cliente...")
		return true
	}else if initipo[0] && !initipo[1]{//si detecta solo centralizado pone al node en modo centralizado
		fmt.Println("Algoritmo de Exclusión Mutua Centralizada detectado...\nAplicando Configuración...\nListo\n")
		tipo_distribucion="centralizado"
		return true
	}else if !initipo[0] && initipo[1]{//si detecta solo distribuido pone al node en modo centralizado
		fmt.Println("Algoritmo de Exclusión Mutua Distribuida detectado...\nAplicando Configuración...\nListo\n")
		tipo_distribucion="distribuido"
		return true
	}else{//si detecta ambas notifica el error y cierra el sistema
		fmt.Println("Inconsistencia de nodos, reinicie el sistema completo por favor")
		return false
		
	}
}
func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal("Error conectando: %v", err)
	}
	s := grpc.NewServer()
	fmt.Println("INICIALIZANDO DATANODE 2...")
	in:=InicioTipo()//verifica el estado de otras maquinas para actualizar de ser necesario el modo a centralizado o distribuido
	if !in{//si hay un estado sin sentido cierra el sistema y notifica
		return
	}
	pb.RegisterClientServiceServer(s, &server{})
	pb2.RegisterNodeServiceServer(s, &server{}) //recibe conexión con el camion
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	
}