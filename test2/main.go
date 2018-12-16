package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func work(stop *bool, name string) error {
	r := rand.New(rand.NewSource(int64(time.Now().Second())))
	for !*stop {
		if r.Int()%10 >= 9 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("work done ---> ", name)
	return nil
}

func doWork(ctx context.Context) error {
	c := make(chan error, 1)
	stop := false
	go func() { c <- work(&stop, "sandy") }()
	select {
	case <-ctx.Done():
		stop = true
		a := <-c //wait for work
		fmt.Println("do work for sandy ---> ", a)
		return ctx.Err()
	case err := <-c:
		fmt.Println("sandy go back (error)")
		return err
	}
}

func main() {
	// timeout := 1 * time.Second
	// ctx, _ := context.WithTimeout(context.Background(), timeout)
	// err := doWork(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		nouse := false
		work(&nouse, "jay")
		cancel()
		fmt.Println("jay over")
	}()
	err := doWork(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
