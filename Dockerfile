FROM golang:1.14-alpine3.12 
WORKDIR /helloworld
ADD . /helloworld
RUN cd /helloworld && go build -o helloworld
RUN ls -alt
EXPOSE 80
ENTRYPOINT ./helloworld