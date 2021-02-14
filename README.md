## 好用的參考網址
1. bloomRPC(好用的client工具，類似於postman) : https://github.com/uw-labs/bloomrpc
2. proto變數型態 : https://developers.google.com/protocol-buffers/docs/proto3
## go 
### 安裝proto工具
```shell
go get google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

### 產生grpc和rpc檔案，兩者都要
```shell
protoc --go_out=golang/server protos/hello.proto
protoc --go-grpc_out=golang/server protos/hello.proto
```

## python
### 安裝proto工具
```shell
python -m pip install grpcio
python -m pip install grpcio-tools
```

### 產生grpc和rpc檔案
```shell
# server的產生grpc文檔方式
python -m grpc_tools.protoc --proto_path=protos/ --python_out=python/server/proto --grpc_python_out=python/server/proto hello.proto

# client的產生grpc文檔方式
python -m grpc_tools.protoc --proto_path=protos/ --python_out=python/client/proto --grpc_python_out=python/client/proto hello.proto


# 記得要去hello_pb2_grpc.py的檔案裡，將
# import hello_pb2 as hello__pb2
# 改成
# from . import hello_pb2 as hello__pb2
```
