package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("💖 Hello World")

	f,err := os.Create("jay.json")


	if err !=nil {

		fmt.Println( err.Error() )
		f.Close()

	} else {

		_,err=f.Write([]byte(`{"name": "💖 Hello Jay"}`))
		if err!= nil{
			fmt.Println( err.Error() )
		}

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
