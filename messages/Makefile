message.pb.go:
	protoc --go_out=plugins=grpc:. message.proto

clean:
	rm -f message.pb.go

deps:
	go get -u google.golang.org/grpc
	go get -u github.com/golang/protobuf/protoc-gen-go
