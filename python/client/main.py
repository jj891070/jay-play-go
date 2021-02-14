from __future__ import print_function
import logging

import grpc

from proto.hello_pb2 import * 
from proto.hello_pb2_grpc import * 


def run():
    # NOTE(gRPC Python Team): .close() is possible on a channel and should be
    # used in circumstances in which the with statement does not fit the needs
    # of the code.
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = GreeterStub(channel)
        response = stub.SayHello(HelloRequest(name='jay'))
    print("Greeter client received: " + response.message)


if __name__ == '__main__':
    logging.basicConfig()
    run()