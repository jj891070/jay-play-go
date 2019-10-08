package main

import (
	"log"
	"os"
	"strconv"
)

// GetURL 取得URL
func GetURL() string {
	url := os.Getenv("WANT_TO_CURL_URL")

	if url == "" {
		url = "https://www.google.com/"
	}

	return url
}

// GetSslSwitch 取得https開關
func GetSslSwitch() bool {
	sslSwitch := os.Getenv("TURN_OFF_SSL_SWITCH")
	if sslSwitch == "" {
		return false
	}

	v, err := strconv.ParseBool(sslSwitch)
	if err != nil {
		log.Println(" ☠ sslSwitch Parse error ----> ", err)
	}
	return v
}

// GetStatusCode 取得status code
func GetStatusCode() bool {
	statusSwitch := os.Getenv("TURN_ON_STATUS_CODE_SWITCH")
	if statusSwitch == "" {
		return false
	}

	v, err := strconv.ParseBool(statusSwitch)
	if err != nil {
		log.Println(" ☠ statusSwitch Parse error ----> ", err)
	}
	return v
}

// GetOutPutMessage 取得輸出訊息
func GetOutPutMessage() bool {
	messageOut := os.Getenv("MESSAGE_OUT")
	if messageOut == "" {
		return false
	}

	v, err := strconv.ParseBool(messageOut)
	if err != nil {
		log.Println(" ☠ messageOut Parse error ----> ", err)
	}

	return v
}

// GetOutRestrictTime 取得限制大於多少秒在輸出訊息
func GetOutRestrictTime() float64 {
	outTime := os.Getenv("RESTRICT_TIME")
	if outTime == "" {
		return 0.0
	}

	s, err := strconv.ParseFloat(outTime, 64)
	if err != nil {
		log.Println(" ☠ outTime Parse error ----> ", err)
		return 0.0
	}
	return s
}

// GetDurationTime 取得延遲多久curl一次
func GetDurationTime() int64 {
	durationTime := os.Getenv("DURATION_TIME")
	if durationTime == "" {
		return 0
	}

	n, err := strconv.ParseInt(durationTime, 10, 64)
	if err != nil {
		log.Println(" ☠ durationTime Parse error ----> ", err)
		return 0
	}
	return n
}

// GetPostData 取得延遲多久curl一次
func GetPostData() string {
	postData := os.Getenv("POST_DATA")
	if postData == "" {
		return postData
	}

	return postData
}

// GetHeaders 取得headers
func GetHeaders() string {
	withHeaders := os.Getenv("WITH_HEADERS")
	if withHeaders == "" {
		return withHeaders
	}

	return withHeaders
}
