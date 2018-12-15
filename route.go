package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SetUpRoute 設定路由
func SetUpRoute() *gin.Engine {
	var r *gin.Engine
	r = gin.Default()
	// // 注册一个路由和处理函数
	v2 := r.Group("/v2", middleware1)
	v2.GET("/login", loginEndpoint)
	v2.GET("/submit", submitEndpoint)
	v2.GET("/read", readEndpoint)
	return r
}

func middleware1(c *gin.Context) {
	log.Println("exec middleware1")

	//你可以写一些逻辑代码

	// 执行该中间件之后的逻辑
	c.Next()
	log.Println("exec last middleware1")
}

func loginEndpoint(c *gin.Context) {
	log.Println("i'm login")

	c.String(http.StatusOK, "hello, login")
	for {
		log.Println("i'm go to sleep")

		time.Sleep(20 * time.Second)
	}
}

func submitEndpoint(c *gin.Context) {
	log.Println("i'm submitEndpoint")
	c.String(http.StatusOK, "hello, submitEndpoint")
}

func readEndpoint(c *gin.Context) {
	log.Println("i'm readEndpoint")
	c.String(http.StatusOK, "hello, readEndpoint")
}
