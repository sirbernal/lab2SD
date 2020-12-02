grpc:
	export GO111MODULE=on  
	go get github.com/golang/protobuf/protoc-gen-go
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	export PATH="$PATH:$(go env GOPATH)/bin"
	go get -u github.com/golang/protobuf/protoc-gen-go

protos: grpc
	protoc --proto_path=. --go_out=plugins=grpc:proto proto/ClientService.proto
	protoc --proto_path=. --go_out=plugins=grpc:proto proto/NodeService.proto

runc:
	cd cliente && \
	go run cliente.go

rund1:
	cd datanode1 && \
	go run datanode.go

rund2:
	cd datanode2 && \
	go run datanode.go

rund3:
	cd datanode3 && \
	go run datanode.go

runn:
	cd namenode && \
	go run namenode.go            
