package main

import (
	"context"
	"fmt"
	"log"
	//"bufio"
	//"io/ioutil"
	"math"
	"os"
	//"strconv"



	pb "github.com/sirbernal/lab2SD/proto/client_service"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	fmt.Println("im fine")
	client := pb.NewClientServiceClient(conn)

	fileToBeChunked := "./ejemplo.pdf" // change here!

	file, err := os.Open(fileToBeChunked)

	if err != nil {
			fmt.Println(err)
			os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()

	const fileChunk = 256000 // 1 MB, change this to your requirement

	// calculate total number of parts the file will be chunked into

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)
	stream, err := client.Upload(context.Background())
	for i := uint64(0); i < totalPartsNum; i++ {

			partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
			partBuffer := make([]byte, partSize)
			fmt.Println("im fine")
			n, _  := file.Read(partBuffer)
			// write to disk
			//fileName := "part_" + strconv.FormatUint(i, 10)
			//_, err := os.Create(fileName)
			fmt.Println("im fine")
			msg:= &pb.UploadRequest{Chunk: partBuffer[:n]}
			fmt.Println("im fine")
			err:= stream.Send(msg)
			if err != nil {
				fmt.Println("cago el chunk")
			}
			resp, err := stream.Recv()

			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			fmt.Println(resp.IdLibro)
			// write/save buffer to disk
	}



	
	
}
