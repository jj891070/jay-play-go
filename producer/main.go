package main

import (
	"log"
	"strconv"

	"github.com/nsqio/go-nsq"
)

func main() {
	config := nsq.NewConfig()

	w, err := nsq.NewProducer("127.0.0.1:4150", config)

	if err != nil {
		log.Panic(err)
	}

	// chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	for i := 0; i < 100000; i++ {
		// buf := make([]int, 4)
		tmp := strconv.Itoa(i)
		log.Printf("Pub: %s", tmp)
		err = w.Publish("test", []byte(tmp))
		if err != nil {
			log.Fatalln("ya!!", err)
		}
		// time.Sleep(time.Second * 1)
	}

	w.Stop()
}
