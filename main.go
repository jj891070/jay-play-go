package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"strings"
	"time"
)

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
	// 延遲多久curl一次
	durationTime := GetDurationTime()
	// 取得要發curl的資料
	postData := GetPostData()
	// 取得要帶的headers
	withHeaders := GetHeaders()
	// 是否要忽略https
	turnOffSsl := GetSslSwitch()
	// 取得status code
	turnOnStatusCode := GetStatusCode()

	var curlGrammar []string
	curlGrammar = append(curlGrammar, wantCurlURL)

	if turnOffSsl {
		curlGrammar = append(curlGrammar, "-k")
	}

	if turnOnStatusCode {
		curlGrammar = append(curlGrammar, "-i")
	}

	if !messageOut {
		curlGrammar = append(curlGrammar, "-s")
		curlGrammar = append(curlGrammar, "-o")
		curlGrammar = append(curlGrammar, `/dev/null`)
	}

	if withHeaders != "" {
		tmp := strings.Split(withHeaders, ",")
		for _, tmpValue := range tmp {
			curlGrammar = append(curlGrammar, "-H")
			curlGrammar = append(curlGrammar, tmpValue)
		}
	}

	if postData != "" {
		curlGrammar = append(curlGrammar, "-X")
		curlGrammar = append(curlGrammar, "POST")
		curlGrammar = append(curlGrammar, "-d")
		curlGrammar = append(curlGrammar, postData)
	}

	curlGrammar = append(curlGrammar, "-w")
	curlGrammar = append(curlGrammar, `####{
				"size_download": %{size_download},
				"speed_download": %{speed_download},
				"time_namelookup": %{time_namelookup},
				"time_connect": %{time_connect},
				"time_appconnect": %{time_appconnect},
				"time_pretransfer": %{time_pretransfer},
				"time_redirect": %{time_redirect},
				"time_starttransfer": %{time_starttransfer},
				"time_total": %{time_total}
			}`)
	data := map[string]float64{}
	for {
		cmd = exec.Command("curl", curlGrammar...)

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})
		cmd.Stderr = stderr
		cmd.Stdout = stdout

		err := cmd.Run()
		if err != nil {
			log.Println(" ☠  Stderr Err ==> ", stderr.String())
			log.Println(" ☠  Curl Command Excute Error => ", err)
			return
		}
		// log.Println("Out ==> ", stdout.String())
		// log.Println(stdout.String())
		tmpStdout := strings.Split(stdout.String(), "####")
		if messageOut {
			log.Println(tmpStdout[0])
		}

		err = json.Unmarshal([]byte(tmpStdout[1]), &data)
		if err != nil {
			log.Println("JSON Error => ", err)
			return
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
		}
		time.Sleep(time.Duration(durationTime) * time.Second)
	}

}
