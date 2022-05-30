package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)




type MyHandler struct{
	Svc http.Server
}

func (m *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func (m *MyHandler)Run(server *http.Server) error {
	fmt.Printf("listening %s\n",server.Addr)
	return server.ListenAndServe()
}

func (m *MyHandler) Close()  error{
	return m.Svc.Close()
}

func main() {
	handler := &MyHandler{}
	svc := &http.Server{
		Addr:   ":8080",
		Handler: handler,
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	group, errCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return handler.Run(svc)
	})

	chanel := make(chan os.Signal, 1)
	signal.Notify(chanel, syscall.SIGINT, syscall.SIGTERM)

	group.Go(func() error {
		for {
			select {
			case <-errCtx.Done():
				e := svc.Close()
				if e!=nil{
					println(e.Error())
				}
				return errCtx.Err()
			case <-chanel:
				cancel()

			}
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}
	println("over")
}

