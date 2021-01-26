package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"tcp_server/internal/pkg/parser"

	"golang.org/x/sync/errgroup"
)

var clientWriteLock sync.Mutex
var seq uint32 = 0

func startClient(msgChan chan string) error {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:8888")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	go clientMsgListen(conn)
	go clientMsgSend(conn, msgChan)
	return nil
}

func clientMsgListen(conn *net.TCPConn) {
	for {
		r, err := parser.ReadFrom(conn)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		if r.Payload == parser.CloseAckMsg {
			conn.Close()
			break
		}
		fmt.Println("收到:" + r.Payload)
	}
}

func clientMsgSend(conn *net.TCPConn, msgChan chan string) {
	for v := range msgChan {
		parser.WriteTo(&parser.RequestResponse{seq, v}, conn, &clientWriteLock)
		seq++
	}
}

func main() {
	ctx := context.Background()
	group, _ := errgroup.WithContext(ctx)
	msgChan := make(chan string, 10)
	group.Go(func() error {
		return startClient(msgChan)
	})
	group.Go(func() error {
		for {
			var paylod string
			fmt.Scanln(&paylod)
			msgChan <- paylod
			if paylod == parser.CloseMsg {
				fmt.Println("recieve close command")
				break
			}
		}
		return errors.New("exit")
	})
	err := group.Wait()
	fmt.Println(err)
}
