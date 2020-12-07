package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

const HOST1 string = "127.0.0.1:8080"
const HOST2 string = "127.0.0.1:8081"

type IndexHandler struct {
	name string
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(h.name))
}

type CloseHandler struct {
	Server *http.Server
}

func (h *CloseHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	_, _ = w.Write([]byte("closing"))

	err := h.Server.Close()
	if err != nil {
		fmt.Printf("service close err %+v \n", err)
	}
}

func startServer() {

	// index1 handler
	mux1 := http.NewServeMux()
	mux1.Handle("/", &IndexHandler{name: "index1"})
	closeHandler1 := &CloseHandler{}
	mux1.Handle("/close", closeHandler1)

	s1Ch := make(chan error)

	// index2 handler
	mux2 := http.NewServeMux()
	mux2.Handle("/", &IndexHandler{name: "index2"})
	closeHandler2 := &CloseHandler{}
	mux2.Handle("/close", closeHandler1)

	s2Ch := make(chan error)

	// 监听系统信号
	ch := make(chan os.Signal, 10)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	ctx := context.Background()
	group, _ := errgroup.WithContext(ctx)

	server1 := &http.Server{Addr: HOST1, Handler: mux1}
	closeHandler1.Server = server1
	server2 := &http.Server{Addr: HOST2, Handler: mux2}
	closeHandler2.Server = server2

	group.Go(func() error {
		// 收到任何一个信号就关闭服务
		select {
		case <-ch:
			fmt.Println("receive close signal!")
		case <-s1Ch:
			fmt.Println("receive server1 close!")
		case <-s2Ch:
			fmt.Println("receive server2 close!")
		}

		signal.Stop(ch)
		err1 := server1.Close()
		fmt.Printf("server1 close %+v \n", err1)

		err2 := server2.Close()
		fmt.Printf("server2 close %+v \n", err2)

		return nil
	})

	group.Go(func() error {
		err := server1.ListenAndServe()
		close(s1Ch)
		fmt.Printf("server1 closed %+v \n", err)
		return nil
	})

	group.Go(func() error {
		err := server2.ListenAndServe()
		close(s2Ch)
		fmt.Printf("server2 closed %+v \n", err)
		return nil
	})

	err := group.Wait()

	fmt.Printf("group err %+v \n", err)
}

func main() {
	startServer()

	// sleep to wait all logs
	time.Sleep(time.Second * 5)
}
