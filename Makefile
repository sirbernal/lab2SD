grpc:
	export GO111MODULE=on  
	go get github.com/golang/protobuf/protoc-gen-go
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	export PATH="$PATH:$(go env GOPATH)/bin"
    go get -u github.com/golang/protobuf/protoc-gen-go

protos: grpc
	protoc --proto_path=. --go_out=plugins=grpc:proto proto/ClientService.proto
	protoc --proto_path=. --go_out=plugins=grpc:proto proto/NodeService.proto

runcliente:
	cd cliente && \
    go run cliente.go
rundatanode1:
	cd datanode1 && \
    go run datanode.go
rundatanode2:
	cd datanode2 && \
    go run datanode.go
rundatanode3:
	cd datanode3 && \
    go run datanode.go
runnamenode:
	cd namenode && \
    go run namenode.go            
