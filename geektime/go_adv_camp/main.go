package main

import (
	"context"
	"fmt"
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
	return server.ListenAndServe()
}

func (m *MyHandler) Close()  error{
	return m.Svc.Close()
}
func handleSignal()  {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()
	signal.Stop(sigs)
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
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
		return handler.Run(s)
	})

	group.Go(func() error {
		<-errCtx.Done()
		return svc.Close()
	})

	chanel := make(chan os.Signal, 1)
	signal.Notify(chanel, syscall.SIGINT, syscall.SIGTERM)

	group.Go(func() error {
		for {
			select {
			case <-errCtx.Done():
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


