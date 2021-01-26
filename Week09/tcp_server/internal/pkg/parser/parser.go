package parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
)

const CloseMsg = "close"
const CloseAckMsg = "ack:close"

type RequestResponse struct {
	Serial  uint32 // 序号
	Payload string // 内容
}

func ReadFrom(conn *net.TCPConn) (*RequestResponse, error) {
	ret := &RequestResponse{}
	buf := make([]byte, 4)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, fmt.Errorf("读长度故障：%s", err.Error())
	}
	length := binary.BigEndian.Uint32(buf)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, fmt.Errorf("读Serial故障：%s", err.Error())
	}
	ret.Serial = binary.BigEndian.Uint32(buf)
	payloadBytes := make([]byte, length-4)
	if _, err := io.ReadFull(conn, payloadBytes); err != nil {
		return nil, fmt.Errorf("读Payload故障：%s", err.Error())
	}
	ret.Payload = string(payloadBytes)
	return ret, nil
}

// 序列化RequestResponse，并发送
// 序列化后的结构如下：
//   长度  4字节
//   Serial 4字节
//   PayLoad 变长
func WriteTo(r *RequestResponse, conn *net.TCPConn, lock *sync.Mutex) {
	lock.Lock()
	defer lock.Unlock()
	payloadBytes := []byte(r.Payload)
	serialBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(serialBytes, r.Serial)
	length := uint32(len(payloadBytes) + len(serialBytes))
	lengthByte := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthByte, length)

	_, _ = conn.Write(lengthByte)
	_, _ = conn.Write(serialBytes)
	_, _ = conn.Write(payloadBytes)
	fmt.Println("发送: " + r.Payload)
}
