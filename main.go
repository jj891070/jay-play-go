package main

import (
	"bytes"
	"log"
	"os/exec"
	"time"
)

func main() {

	var (
		cmd            *exec.Cmd
		commandGrammar []string
	)

	// server 回傳的內容是否要輸出
	messageOut := GetOutPutMessage()
	// 當時間大於幾秒（second）的時候輸出
	outTime := GetOutRestrictTime()
	// 延遲多久curl一次
	durationTime := GetDurationTime()
	// 取得要執行的模式
	mode := GetMode()
	// 取得要執行的指令
	shellCommand := GetShellCommand(mode)

	switch shellCommand {
	case "nc":
		commandGrammar = GetNcCommand()
	case "curl":
		commandGrammar = GetCurlCommand(messageOut)
	default:
		log.Println(" 😱 ", shellCommand)
		return
	}

	for {
		cmd = exec.Command(shellCommand, commandGrammar...)

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})
		cmd.Stderr = stderr
		cmd.Stdout = stdout

		err := cmd.Run()
		if err != nil {
			log.Println(" ☠  Stderr Err ==> ", stderr.String())
			log.Println(" ☠  Shell Command Excute Error => ", err)
			log.Println("=========================================")
		}
		// log.Println("Out ==> ", stdout.String())
		// log.Println(stdout.String())
		if shellCommand == "nc" {
			ParseNcStdout(stdout.String(), messageOut, outTime)
		} else {
			ParseCurlStdout(stdout.String(), messageOut, outTime)
		}

		time.Sleep(time.Duration(durationTime) * time.Second)
	}

}
