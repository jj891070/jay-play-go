package main

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron"
)

func main() {
	log.Println("Starting...")

	c := cron.New() // 新建一個定時任務物件
	c.AddFunc("* * * * * *", func() {
		log.Println("hello world!!!")
	}) // 給物件增加定時任務

	c.AddFunc("30 * * * *", func() { fmt.Println("Every hour on the half hour") })

	c.Start()
	c.Entries()
	time.Sleep(10 * time.Second)
	log.Println("Stop start")
	c.Stop()
	log.Println("Stop end")

	time.Sleep(10 * time.Second)

}
