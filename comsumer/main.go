package main

import (
	"log"
	"strconv"
	"sync"

	"github.com/nsqio/go-nsq"
)

func main() {

	wg := &sync.WaitGroup{}
	wg.Add(1000)

	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer("test", "ch2", config)
	var last int = -1
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("Got a message: %s", message.Body)
		tmp, err := strconv.Atoi(string(message.Body))
		if err != nil {
			log.Fatal(err)
		}
		if tmp < last {
			log.Fatalln("wait ")
		} else {
			last = tmp
		}

		wg.Done()
		return nil
	}))
	err := q.ConnectToNSQLookupd("127.0.0.1:4161")
	if err != nil {
		log.Panic(err)
	}
	wg.Wait()

}
