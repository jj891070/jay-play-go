package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
)

var serverClose chan os.Signal

// GracefulDown 優雅結束程式
func GracefulDown() <-chan os.Signal {
	return serverClose
}

func main() {

	tellParentMessage := make(chan int, 1)
	// 設定伺服器路由
	r := SetUpRoute()

	server := SetServerConf(r, ":8098")

	go func() {
		// 绑定端口，然后启动应用
		err := ServerRun(server)
		if err != nil {
			fmt.Println("server error ---> ", err)
		}
		tellParentMessage <- 2
	}()

	go func() {
		serverClose = make(chan os.Signal, 1)
		signal.Notify(serverClose, os.Interrupt, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

		s := <-GracefulDown()
		switch s {
		case os.Interrupt:
			log.Infof("SIGSTOP")
			tellParentMessage <- 1
			return
		case syscall.SIGTERM:
			log.Infof("SIGSTOP")
			tellParentMessage <- 1
			return
		case syscall.SIGINT:
			log.Infof("SIGHUP")
			tellParentMessage <- 1
			return
		case syscall.SIGKILL:
			log.Infof("SIGKILL")
			tellParentMessage <- 1
			return
		default:
			log.Infof("default")
			tellParentMessage <- 1
			return
		}

	}()

	var force bool
	for {
		p := <-tellParentMessage
		switch p {
		case 1:
			fmt.Println("來自使用者")
			force = true
		case 2:
			fmt.Println("來自伺服器")
		default:
			fmt.Println("無名")
		}
		if force {
			break
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	fmt.Println("Server exiting")

}

// SetServerConf 設定伺服器conf
func SetServerConf(router *gin.Engine, port string) (srv *http.Server) {
	srv = &http.Server{
		Addr:    port,
		Handler: router,
	}
	return
}

// ServerRun 跑起伺服器
func ServerRun(serverConf *http.Server) (err error) {
	err = serverConf.ListenAndServe()
	return
}
