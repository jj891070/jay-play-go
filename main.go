package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET(
		"/*any",
		func(c *gin.Context) {
			// 記錄輸入
			// 1. Path
			// 2. Params
			// 3. Header
			path := c.Request.URL.Path
			log.Println(path)
			reqHeader := c.Request.Header
			log.Println(reqHeader)
			log.Println("==========================")
			temp := map[string]string{}
			for i := range reqHeader {
				temp[i] = reqHeader.Get(i)
			}

			reqHeaderByte, err := json.Marshal(temp)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("reqHeaderByte ---> ", string(reqHeaderByte))

			switch c.Request.Method {
			case "GET":
				tmp := map[string]string{}
				query := c.Request.URL.Query()
				for k := range query {
					tmp[k] = query.Get(k)
				}
				buf, err := json.Marshal(tmp)
				if err == nil {
					fmt.Println("JSON ---> ", string(buf))
				}
			default:
				// 讀取請求的Body
				buf, err := ioutil.ReadAll(c.Request.Body)
				if err == nil {
					fmt.Println("JSON ---> ", string(buf))
				}

				// 重新複製一份Body存回去請求連線
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
			}

			c.Next()

			// 記錄輸出
			// 1. HTTP Status
			// 2. Header
			// 3. ErrorCode
			// 3. ErrorText
			// 3. Response Data
			fmt.Println("HTTP CODE", c.Writer.Status())
			// res := struct {
			// 	ErrorCode int64       `json:"error_code"`
			// 	ErrorText string      `json:"error_text"`
			// 	Data      interface{} `json:""`
			// }{}
			// json.Unmarshal(c.Writer, &res)

			// 塞入DB
		},
		func(c *gin.Context) {
			fmt.Println("=============")
			c.JSON(http.StatusOK, "OK")
		},
	)
	r.POST(
		"/*any",
		func(c *gin.Context) {
			input := map[string]interface{}{}
			err := c.ShouldBindJSON(&input)
			if err != nil {
				c.JSON(http.StatusOK, "Error -> "+err.Error())
				return
			}
			c.JSON(http.StatusOK, input)
			c.JSON(http.StatusOK, "OK")
		},
	)
	r.Run(":8002")
}
