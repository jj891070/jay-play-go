package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func work(stop *bool) error {
	r := rand.New(rand.NewSource(int64(time.Now().Second())))
	for !*stop {
		if r.Int()%10 >= 9 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("work done")
	return nil
}

func doWork(ctx context.Context) error {
	c := make(chan error, 1)
	stop := false
	go func() { c <- work(&stop) }()
	select {
	case <-ctx.Done():
		stop = true
		<-c //wait for work
		return ctx.Err()
	case err := <-c:
		return err
	}
}

func main() {
	timeout := 1 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	go func() {
		nouse := false
		work(&nouse)
		cancel()
	}()
	err := doWork(ctx)
	if err != nil {
		fmt.Println(err)
	}

}
