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
var datanode = []string{"localhost:50052","localhost:50053","localhost:50054"}
var directions = []string{"localhost:50055", "localhost:50052","localhost:50053","localhost:50054"}
var nodemode = []string{"","","",""}


var chunks [][]byte //donde guardo los chunks para subir
var rechunks [][]byte //donde guardo los chunks para bajar
var tipo_distribucion string 

type ChunkAndN struct{
	Chunk []byte
	N int 
}


func Chunker(archivo string) []ChunkAndN{
	file, err := os.Open(archivo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} 
	defer file.Close()
	var chunks []ChunkAndN
	fileInfo, _ := file.Stat()
	var fileSize int64 = fileInfo.Size()
	const fileChunk = 256000
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
		chunks = append(chunks,ChunkAndN{partBuffer,n})
		fmt.Println("Se creo y guardo un chunk")
	}
	file.Close()
	return chunks
}
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
	for _,chunk:= range rechunks{
		_, err := file.Write(chunk)	
		if err != nil {
				fmt.Println(err)
				os.Exit(1)
		}
		file.Sync()
	}
	file.Close()
	rechunks = [][]byte{}
}

func SolicitarLibros()[]string{
	conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := pb.NewClientServiceClient(conn)
    
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	msg:= &pb.DownloadNamesRequest{Req: "Libros",}
	resp, err := client.DownloadNames(ctx, msg)
	
	fmt.Println("Archivos(s) disponible(s)")
	for i,name:=range resp.GetNames(){
		fmt.Println(strconv.Itoa(i+1)+".-"+name)
	}
	return resp.GetNames()
}
func SolicitarUbicacion(libro string)([]int64){
	conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := pb.NewClientServiceClient(conn)
    
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	msg:= &pb.LoCRequest{Req: libro,}
	resp, err := client.LocationsofChunks(ctx, msg)
	
	return resp.GetLocation()
		
}
func DescargarChunks(name string, prop []int64){
	for i, node:= range prop{
		fmt.Println(i)
		//pide el chunk a la maquina que le corresponde+
		datan:=datanode[int(node)]
		// i es la parte, datan es la direccion y name, el nombre del archivo
		chunkadescargar:=name+"_"+strconv.Itoa(i)
		fmt.Println(datan,chunkadescargar)

		//////////// Pedir el chunk al datanode correspondiente ///////////
		conn, err := grpc.Dial(datan, grpc.WithInsecure())
		if err != nil {
			log.Fatalln(err)
		}
		defer conn.Close()

		client := pb.NewClientServiceClient(conn)
		
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		//fmt.Println(chunkadescargar)
		msg:= &pb.DownloadChunksRequest{Name: chunkadescargar}
		resp, err := client.DownloadChunks(ctx, msg)
		//////////////////////////////////////////////////////////////////

		//aca se pide el chunk al datanode q corresponda
		//guardarlo en rechunk
		//chunky:=[]byte{}

		rechunks=append(rechunks,resp.GetChunk())
	}
	//fmt.Println(rechunks)
	Unchunker(name)
}
func SubirArchivo(node int, archivo string)(){
	conn, err := grpc.Dial(datanode[node], grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := pb.NewClientServiceClient(conn)
    
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	chunks := Chunker("./"+archivo)
	msg:= &pb.UploadRequest{Tipo: 0, Nombre: archivo, Totalchunks: int64(len(chunks))}
	resp, err := client.Upload(ctx, msg)
	fmt.Println(len(chunks))
	if resp.GetResp()==int64(0){
		for i,chunk :=range chunks{
			msg:= &pb.UploadChunksRequest{Chunk: chunk.Chunk[:chunk.N]}
			resp, err := client.UploadChunks(ctx, msg)
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			if i==(len(chunks)-1){
				fmt.Println("aqui avisa q mando todo al namenode y manda el len")
			}
			fmt.Println(resp.GetResp()) 
		}	
	}
}
func menu_centralizado(){
	var menu2 string
	Menu2:
		for {
			fmt.Print("\n\n---Menu Algoritmo de Exclusión Mutua Centralizada--- \nIngrese opción\n1.-Subir Archivo\n2.-Ver archivos en sistema\n3.-Cerrar Sistema\nIngrese opción:")
			_,err:=fmt.Scanln(&menu2)
				if err!=nil{
					fmt.Print("\nFormato de ingreso no válido, pruebe nuevamente...")
					continue
				}
			switch menu2{
			case "1": 
				var datasubida string
				var dataint int
				fmt.Println("--Lista de Datanodes--")
				for i,j:= range datanode{
					fmt.Println(strconv.Itoa(i+1)+".- Datanode "+strconv.Itoa(i+1)+" Ip:"+j)
				}
				for {
					fmt.Print("Seleccione un datanode:")
					_,err:=fmt.Scanln(&datasubida)
					if err!=nil{
						fmt.Println("Formato no válido, pruebe nuevamente...")
						continue
					}
					dataint,err=strconv.Atoi(datasubida)
					if err!=nil{
						fmt.Println("Formato no válido, pruebe nuevamente...")
						continue
					}
					if dataint>0 && dataint<len(datanode)+1{
						break
					}else{
						fmt.Println("Opción no válida, pruebe nuevamente...")
					}
				}
				var archivo string
				fmt.Print("Ingrese nombre de archivo con su extension(ejemplo: archivo.pdf):")
				fmt.Scanln(&archivo)
				SubirArchivo(dataint-1,archivo)
			case "2":
				libros:=SolicitarLibros()
				if len(libros)==0{
					fmt.Println("\n\nSistema sin archivos disponibles")
					continue Menu2
				}
				var opt string
				var optint int
				for{
					fmt.Print("Seleccione un archivo para descargar o escriba 'n' para volver: ")
					_,err:=fmt.Scanln(&opt)
					if err!=nil{
						fmt.Println("Formato no válido, pruebe nuevamente...")
						continue
					}
					if opt=="n"{
						continue Menu2
					}
					optint,err=strconv.Atoi(opt)
					if err!=nil{
						fmt.Println("Formato no válido, pruebe nuevamente...")
						continue
					}
					if optint>0 && optint<len(libros)+1{
						break
					}else{
						fmt.Println("Opción no válida, pruebe nuevamente...")
					}
				}
				fmt.Print("Descargando archivo: "+libros[optint-1])
				ubicacion := SolicitarUbicacion(libros[optint-1])
				DescargarChunks(libros[optint-1], ubicacion)
			case "3":
				fmt.Println("Terminando ejecución cliente...")
				break Menu2
			default:
				fmt.Println("\nFormato u opción no válida, pruebe nuevamente:\n\n")
				continue Menu2
			}
		}
}
func menu_distribuido(){
	fmt.Println("Aqui deberia aplicarle el ricky wala")
	menu_centralizado()
}
func DistribStatus()int{ //0 estan sin modo, 1 esta centralizado, 2 esta distribuido, 3 fallo todo
	cent:=0
	dist:=0
	for _,stat:=range nodemode{
		if stat=="centralizado"{
			cent++
		}else if stat=="distribuido"{
			dist++
		}
	}
	if dist==0 && cent==0{
		return 0
	}
	if cent>0 && dist==0{
		return 1
	}
	if dist>0 && cent==0{
		return 2
	}
	if dist>0 && cent>0{
		return 3
	}
	return -1
}
func VerifInicial()int{
	for i, dire := range directions{
		conn, err := grpc.Dial(dire, grpc.WithInsecure())
		if err != nil {
			//fmt.Println("No esta el nodo")
			continue
		}
		defer conn.Close()

		client := pb.NewClientServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		msg:= &pb.TypeRequest{Type: "inicio" }

		resp, err := client.TypeDis(ctx, msg)
		if err != nil {
			//fmt.Println("No esta el nodo")
			continue
			}
		nodemode[i]=resp.GetResp()
	}
	return DistribStatus()
}
func menu(){
	init:=VerifInicial()
	switch init{
	case 0: 
		var menu1 string
		for {
			fmt.Print("\n\n---Menu Principal--- \nIngrese que algoritmo de exclusión mutua desea experimentar\n1.-Algoritmo de Exclusión Mutua Centralizada\n2.-Algoritmo de Exclusión Mutua Distribuida\nIngrese opción:")
			_,err:=fmt.Scanln(&menu1)
				if err!=nil{
					fmt.Print("\nFormato de ingreso no válido, pruebe nuevamente:")
					continue
				}
			if menu1=="1"{
				tipo_distribucion = "centralizado"
				
				for _, dire := range directions{
					conn, err := grpc.Dial(dire, grpc.WithInsecure())
					if err != nil {
						fmt.Println("No esta el nodo")
						continue
					}
					defer conn.Close()

					client := pb.NewClientServiceClient(conn)

					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()
					msg:= &pb.TypeRequest{Type: tipo_distribucion }

					_, err = client.TypeDis(ctx, msg)
					if err != nil {
						fmt.Println("No esta el nodo")
						continue
						}
				}
				menu_centralizado()
				break
			}else if menu1=="2"{
				tipo_distribucion = "distribuido"
				
				for _, dire := range directions{
					conn, err := grpc.Dial(dire, grpc.WithInsecure())
					if err != nil {
						fmt.Println("No esta el nodo")
						continue
					}
					defer conn.Close()

					client := pb.NewClientServiceClient(conn)

					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()
					msg:= &pb.TypeRequest{Type: tipo_distribucion }

					_, err = client.TypeDis(ctx, msg)
					if err != nil {
						fmt.Println("No esta el nodo")
						continue
						}
				}
				menu_distribuido()
				break
			}else{
				fmt.Println("\nFormato u opción no válida, pruebe nuevamente:\n\n")
				continue
			}
		}
	case 1:
		fmt.Println("Algoritmo de Exclusión Mutua Centralizada Detectado\nRedireccionando a Menu del algoritmo...")
		menu_centralizado()
	case 2:
		fmt.Println("Algoritmo de Exclusión Mutua Distribuida Detectado\nRedireccionando a Menu del algoritmo...")
		menu_distribuido()
	case 3:
		fmt.Println("Algoritmo corrompido detectado! (Presencia de distintos algoritmos en distintos nodos)\nPor favor, reiniciar el sistema completo\n\n---Cerrando Sistema--")
		return
	}
	
}
func main() {
	menu()
	return
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := pb.NewClientServiceClient(conn)
    
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	chunks := Chunker("./ejemplo.pdf")
	msg:= &pb.UploadRequest{Tipo: 0, Nombre: "lel.pdf", Totalchunks: int64(len(chunks))}
	resp, err := client.Upload(ctx, msg)
	fmt.Println(len(chunks))
	if resp.GetResp()==int64(0){
		for i,chunk :=range chunks{
			msg:= &pb.UploadChunksRequest{Chunk: chunk.Chunk[:chunk.N]}
			resp, err := client.UploadChunks(ctx, msg)
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			if i==(len(chunks)-1){
				fmt.Println("aqui avisa q mando todo al namenode y manda el len")
			}
			fmt.Println(resp.GetResp()) 
		}	
	}
	SolicitarLibros()
	ubicacion := SolicitarUbicacion("lel.pdf")
	DescargarChunks("lel.pdf", ubicacion)
	fmt.Println(ubicacion)
	

	
	/*for i,chunk :=range chunks{
		msg:= &pb.UploadChunksRequest{Chunk: chunk.Chunk[:chunk.N]}
		resp, err := client.Upload(ctx, msg)
		if err != nil {
			log.Fatalf("can not receive %v", err)
		}
		if i==(len(chunks)-1){
			fmt.Println("aqui avisa q mando todo al namenode y manda el len")
		}
		fmt.Println(resp.GetIdLibro()) 
	}*/

			// write to disk
			//fileName := "part_" + strconv.FormatUint(i, 10)
			//_, err := os.Create(fileName)

			//msg:= &pb.UploadRequest{Chunk: partBuffer[:n]}

			//stream.Send(msg)
			/*if err != nil {
				fmt.Println("cago el chunk")
			}*/
			//resp, err := stream.Recv()

			/*if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			fmt.Println(resp.IdLibro) 
			// write/save buffer to disk
	}*/



	
	
}
