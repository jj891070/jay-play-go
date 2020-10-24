package main

import (
	"flag"
	"log"
	"sync"
	"time"

	redis "github.com/go-redis/redis"
)

var agent = flag.Bool("b", false, "bool类型参数")

func main() {
	flag.Parse()
	conn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if *agent {
		log.Println("A")
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			log.Println("B")
			sub := conn.Subscribe("out1")
			conn.LPush("in", "1")
			log.Println(<-sub.Channel())
			wg.Done()
		}()
		go func() {
			log.Println("C")
			sub := conn.Subscribe("out2")
			conn.LPush("in", "2")
			log.Println(<-sub.Channel())
			wg.Done()
		}()

		wg.Wait()
		return
	}

	for {
		in, err := conn.BRPop(time.Second*5, "in").Result()
		if err != nil {
			if err == redis.Nil {
				log.Println("ya!!")
				continue
			}
			log.Println(err)
			continue
		}

		log.Println("進來: %+v", in)
		conn.Publish("out"+in[1], "ok -> "+in[1])
	}

}
