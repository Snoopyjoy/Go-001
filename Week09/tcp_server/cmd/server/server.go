package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"tcp_server/internal/pkg/parser"
)

var serverWriteLock sync.Mutex

func startServer() {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8888")
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	defer tcpListener.Close()
	fmt.Println("server start and waiting ...")
	go closeListener(tcpListener)
	for {
		conn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println("AcceptTCP err", err)
			break
		}
		msgChan := make(chan *parser.RequestResponse)
		fmt.Println("a new connection:" + conn.RemoteAddr().String())
		go serverMsgListen(conn, msgChan)
		go serverMsgSend(conn, msgChan)
	}
}

func closeListener(tcpListener *net.TCPListener) {
	ch := make(chan os.Signal, 10)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-ch
	tcpListener.Close()
}

func serverMsgListen(conn *net.TCPConn, msgChan chan *parser.RequestResponse) {
	for {
		r, err := parser.ReadFrom(conn)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		if r.Payload == parser.CloseMsg {
			msgChan <- r
			conn.Close()
			close(msgChan)
			break
		}
		msgChan <- r
	}
}

func serverMsgSend(conn *net.TCPConn, msgChan chan *parser.RequestResponse) {
	for v := range msgChan {
		if v == nil {
			break
		}
		parser.WriteTo(&parser.RequestResponse{v.Serial, fmt.Sprintf("ack:%s", v.Payload)}, conn, &serverWriteLock)
	}
}

func main() {
	startServer()
}
