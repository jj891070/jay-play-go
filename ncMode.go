package main

import (
	"log"
	"strings"
)

// GetNcCommand 取得netcat執行參數
func GetNcCommand() (ncGrammar []string) {

	// 你想要curl的網址
	wantCurlURL := GetURL()
	tmp := strings.Split(wantCurlURL, ":")
	ncGrammar = append(ncGrammar, "-vz")
	ncGrammar = append(ncGrammar, tmp[0])
	if tmp[1] == "" {
		ncGrammar = append(ncGrammar, "80")
		return
	}
	ncGrammar = append(ncGrammar, tmp[1])
	return
}

// ParseNcStdout 解析netcat回傳
func ParseNcStdout(result string, messageOut bool, outTime float64) {
	if result != "" && messageOut {
		log.Println("port is ok ---> ", result)
	} else if result == "" && messageOut {
		log.Println("port is ok. but no response")
	}
	return
}
