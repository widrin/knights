# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
protoc --go_out=.     --go-grpc_out=.     ./center.proto
