package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

const socketPath = "/tmp/my_unix_socket"

func main() {
	// 删除之前可能残留的套接字文件
	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		fmt.Println("无法删除套接字文件:", err)
		return
	}

	serverMux := http.NewServeMux()
	serverMux.Handle("/", http.FileServer(http.Dir("html")))
	serverMux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello golang network developer!"))
	})

	server := http.Server{Handler: serverMux}
	// 创建Unix域套接字监听器
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println("无法创建监听器:", err)
		return
	}
	defer listener.Close()
	log.Println("開始監聽")

	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}
