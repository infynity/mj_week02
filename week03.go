package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main(){

	g, ctx := errgroup.WithContext(context.Background())

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(9e9)
		writer.Write([]byte("are you ok"))

	})

	serverOut:=make(chan struct{})

	mux.HandleFunc("/simulateShutdown", func(writer http.ResponseWriter, request *http.Request) {
		serverOut<- struct{}{}
	})
	server := http.Server{
		Handler: mux,
		Addr:    ":8088",
	}

	g.Go(func() error {
		return server.ListenAndServe()
	})


	g.Go(func() error {
		fmt.Println(123)
		select {
			case <-ctx.Done():
				fmt.Println("g2 errgrp exit")
			case <-serverOut:
				fmt.Println("server out")
		}

		timeoutCtx, _ := context.WithTimeout(context.Background(), 3e9)
		return server.Shutdown(timeoutCtx)
	})


	g.Go(func() error {
		sgchan := make(chan os.Signal)

		signal.Notify(sgchan,os.Interrupt,syscall.SIGTERM)

		select {
		case <-ctx.Done():
			fmt.Println("g3 now is quit")
			return ctx.Err()
		case sig:=<-sgchan:
			return errors.Errorf("accept signal %s",sig)
		}
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("errgroup exiting: %+v\n", err)
	}
}
