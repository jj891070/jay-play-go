package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	var (
		graceShotdown int
		addr          string
		help          bool
	)
	flag.BoolVar(&help, "h", false, "this help")
	flag.IntVar(&graceShotdown, "grace", 15, "grace shotdown saecond")
	flag.StringVar(&addr, "addr", ":8099", "http listen on")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	c := make(chan error, 1)

	sv := &http.Server{
		Addr:    ":8098",
		Handler: SetUpRoute(),
	}

	go func(c chan error) {
		// 绑定端口，然后启动应用
		err := sv.ListenAndServe()
		if err != nil {
			c <- err
		}
	}(c)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-sig:
		log.Println("os signal:", s)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(graceShotdown)*time.Second)
		defer cancel()
		if err := sv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
	case e := <-c:
		log.Fatal("Server Start:", e)
	}

	log.Println("server exit")
}
