## 好用的client工具，類似於postman
1. bloomRPC : https://github.com/uw-labs/bloomrpc
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
python -m grpc_tools.protoc --proto_path=protos/ --python_out=python/server/ --grpc_python_out=python/server/ hello.proto

# 記得要去hello_pb2_grpc.py的檔案裡，將
# import hello_pb2 as hello__pb2
# 改成
# from . import hello_pb2 as hello__pb2
```
