package gateI

import (
	"net"
)

type IConn interface {
	ReadMsg() ([]byte, error)
	WriteMsg(args []byte) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	SetMsgType(v int)
	RealIp() string //真实ip，websocket 可以获取
}
