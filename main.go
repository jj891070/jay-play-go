package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var (
		wg  sync.WaitGroup
		jay int64
		l   *sync.Mutex
	)
	// 同步鎖
	l = new(sync.Mutex)

	// 每１秒鐘做一次
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	// 通知結束的channel
	c := make(chan int)

	// 共同資料，一個讀一個寫
	jay = 56

	// 工人一
	wg.Add(1)
	go func(wg *sync.WaitGroup, c *chan int, jay *int64) {
	ForEnd:
		for {

			select {

			case <-t.C:
				fmt.Println(time.Now())
			case <-*c:
				fmt.Println("break")
				break ForEnd
			default:
				fmt.Println("I'm start!!!!!!!!!!!!!!!!! ---> ", *jay)
				time.Sleep(10 * time.Second)
				fmt.Println("I'm default!!!!!!!!!!!!!!!!!----> ", *jay)
			}
		}
		fmt.Println("bye")
		wg.Done()
	}(&wg, &c, &jay)

	// 工人二
	wg.Add(1)
	go func(wg *sync.WaitGroup, c *chan int, jay *int64, l *sync.Mutex) {

		time.Sleep(10 * time.Second)
		*c <- 1
		wg.Done()
	}(&wg, &c, &jay, l)

	wg.Wait()
	fmt.Println("all the tasks done＠＠＠")

	time.Sleep(10 * time.Second)
	fmt.Println("all the tasks done...")
}
