package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
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

// GetOutPutMessage 取得輸出訊息
func GetOutPutMessage() bool {
	messageOut := os.Getenv("MESSAGE_OUT")
	if messageOut == "" {
		return false
	}

	v, err := strconv.ParseBool(messageOut)
	if err != nil {
		log.Println(" ☠  Parse error ----> ", err)
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
		log.Println(" ☠  Parse error ----> ", err)
	}
	return s
}

func main() {

	var (
		cmd *exec.Cmd
	)
	// 你想要curl的網址
	wantCurlURL := GetURL()
	// server 回傳的內容是否要輸出
	messageOut := GetOutPutMessage()
	// 當時間大於幾秒（second）的時候輸出
	outTime := GetOutRestrictTime()

	if messageOut {
		cmd = exec.Command("curl", wantCurlURL, "-w", `
			{
				"size_download": %{size_download},
				"speed_download": %{speed_download},
				"time_namelookup": %{time_namelookup},
				"time_connect": %{time_connect},
				"time_appconnect": %{time_appconnect},
				"time_pretransfer": %{time_pretransfer},
				"time_redirect": %{time_redirect},
				"time_starttransfer": %{time_starttransfer},
				"time_total": %{time_total}
			}`,
		)
	} else {
		cmd = exec.Command("curl", wantCurlURL, "-w", `
			{
				"size_download": %{size_download},
				"speed_download": %{speed_download},
				"time_namelookup": %{time_namelookup},
				"time_connect": %{time_connect},
				"time_appconnect": %{time_appconnect},
				"time_pretransfer": %{time_pretransfer},
				"time_redirect": %{time_redirect},
				"time_starttransfer": %{time_starttransfer},
				"time_total": %{time_total}
			}
			`,
			"-s", "-f", "-o", "/dev/null",
		)
	}

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	err := cmd.Run()
	if err != nil {
		log.Println(" ☠  Err ==> ", stderr.String())
		log.Println(" ☠  Run Error => ", err)
	}
	// log.Println("Out ==> ", stdout.String())

	data := map[string]float64{}
	err = json.Unmarshal(stdout.Bytes(), &data)
	if err != nil {
		log.Println("JSON Error => ", err)
	}

	total := data["time_total"]
	nameLookup := data["time_namelookup"]
	tcpConnect := data["time_pretransfer"] - data["time_namelookup"]
	sslConnect := data["time_appconnect"]
	preTransfer := data["time_pretransfer"] - data["time_appconnect"]
	if data["time_appconnect"] == 0 {
		sslConnect = 0
		preTransfer = data["time_pretransfer"] - data["time_connect"]
	}
	redirect := data["time_redirect"]
	serverHandle := data["time_starttransfer"] - data["time_pretransfer"]
	returnTime := data["time_total"] - data["time_starttransfer"]
	if outTime != 0 && total > outTime {
		log.Printf(`
		🍾 總時間 %f
			-> 解析網址 %f (%.3f％)
			-> TCP握手 %f (%.3f％)
			-> SSL檢查 %f (%.3f％)
			-> 傳入資料 %f (%.3f％)
			-> 轉導 %f (%.3f％)
			-> Server處理時間 %f (%.3f％)
			-> 內容傳輸時間 %f (%.3f％)
	`,
			total,
			nameLookup, nameLookup/total*100,
			tcpConnect, tcpConnect/total*100,
			sslConnect, sslConnect/total*100,
			preTransfer, preTransfer/total*100,
			redirect, redirect/total*100,
			serverHandle, serverHandle/total*100,
			returnTime, returnTime/total*100,
		)
		return
	}

	if outTime == 0 {
		log.Printf(`
		🌞 總時間 %f
			-> 解析網址 %f (%.3f％)
			-> TCP握手 %f (%.3f％)
			-> SSL檢查 %f (%.3f％)
			-> 傳入資料 %f (%.3f％)
			-> 轉導 %f (%.3f％)
			-> Server處理時間 %f (%.3f％)
			-> 內容傳輸時間 %f (%.3f％)
	`,
			total,
			nameLookup, nameLookup/total*100,
			tcpConnect, tcpConnect/total*100,
			sslConnect, sslConnect/total*100,
			preTransfer, preTransfer/total*100,
			redirect, redirect/total*100,
			serverHandle, serverHandle/total*100,
			returnTime, returnTime/total*100,
		)
		return
	}

}
