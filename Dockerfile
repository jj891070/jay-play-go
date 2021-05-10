FROM golang:1.16-alpine AS builder

WORKDIR /go/src/app
COPY . .
# ARG SVC=server
# RUN go build -o app ./${SVC}/
RUN go build -o app .
RUN rm -rf vendor *.go go.*

FROM alpine:3.12
WORKDIR /root/
COPY --from=builder /go/src/app/app .

CMD ["app"]