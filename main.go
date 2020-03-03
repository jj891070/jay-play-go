package main

import (
	"io"
	"log"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
)

func main() {
	gin.DefaultWriter = colorable.NewColorableStderr()
	r := gin.Default()
	r.GET("/stream", func(c *gin.Context) {
		chanStream := make(chan int, 10)
		go func() {
			defer close(chanStream)
			for i := 0; i < 5; i++ {
				chanStream <- i
				time.Sleep(time.Second * 1)
			}
		}()
		c.Stream(func(w io.Writer) bool {
			if msg, ok := <-chanStream; ok {
				c.SSEvent("message", msg)
				return true
			}
			log.Println("cancel")
			return false
		})
	})
	r.Use(static.Serve("/", static.LocalFile("./public", true)))
	r.Run(":8087")
}
