FROM golang:1.11-alpine AS builder

WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o app .
# RUN ls -alt
# RUN pwd


FROM alpine
RUN apk --no-cache add curl
COPY --from=builder /go/src/app/app /go/src/app
WORKDIR /go/src
# RUN pwd
# RUN ls -alt
ENTRYPOINT ["./app"]
