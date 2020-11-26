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
	//"strconv"



	pb "github.com/sirbernal/lab2SD/proto/client_service"
	"google.golang.org/grpc"
)

var chunks [][]byte

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

func main() {
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
