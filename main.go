package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	fmt.Println("💖 Hello World")

	f, err := os.OpenFile("/ocean/log/jay.json",
	os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	var data string
	now := time.Now() 


	data = `{"name": "💖 Hello Jay `+now.Format(time.RFC3339) +` "}`+"\n"
	if _, err := f.WriteString(data); err != nil {
		log.Println(err)
		f.Close()
		return
	}

	f.Close()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello E.F.K")
	})

	r.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello API")
	})

	r.Run()
}
