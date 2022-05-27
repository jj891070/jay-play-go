package main

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	// You can generate a Token from the "Tokens Tab" in the UI
	const token = "qwe123"
	const bucket = "jay-test"
	const org = "agocean"

	client := influxdb2.NewClient("http://127.0.0.1:8086", token)
	// Get non-blocking write client
	writeAPI := client.WriteAPI(org, bucket)

	poolNum := 10
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = poolNum
	transport.MaxConnsPerHost = poolNum
	transport.MaxIdleConnsPerHost = poolNum
	apiclient := http.Client{
		// Timeout:   30 * time.Second,
		Transport: transport,
	}

	// write some points
	for i := 0; i < 100; i++ {
		code, _, _ := HTTPGet(
			apiclient,
			"https://agapi.bravocasino.net/healthz", nil,
		)

		// create point
		p := influxdb2.NewPoint(
			"sc-api",
			map[string]string{
				"hostname": "health-api",
			},
			map[string]interface{}{
				"status":  code,
				"resTime": rand.Float64() * 300.0,
			},
			time.Now())
		// write asynchronously
		writeAPI.WritePoint(p)
	}
	// Force all unwritten data to be sent
	writeAPI.Flush()
	// always close client at the end
	defer client.Close()
}

func HTTPGet(
	client http.Client,
	url string, wg *sync.WaitGroup,
) (statusCode int, body string, err error) {
	//发送请求获取响应
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	//结束网络释放资源
	if resp != nil {
		defer func() {
			resp.Body.Close()
			if wg != nil {
				wg.Done()
			}
		}()
	}
	//判断响应状态码
	statusCode = resp.StatusCode

	//读取响应实体
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	body = string(bs)
	return
}
