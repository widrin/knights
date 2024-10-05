package gateI

import "net"

type NetEvent struct {
	OnAgentInit    func(IAgent)
	OnAgentDestroy func(IAgent)
	OnReceiveMsg   func(agent IAgent, data []byte)
}
type IAgent interface {
	UserData() interface{}
	WriteBuff(buff []byte)
	SetUserData(data interface{})
	GetId() int32
	RemoteAddr() net.Addr
	RealIp() string
	Close()
	Run()
	OnClose()
	SetAutoReconnect(v bool)
	GetAutoReconnect() bool
	SetInitTime(v int64)
	GetInitTime() int64
}
type IGate interface {
	SetMaxConnNum(num int)
	Run(closeSig chan bool)
	SetWSAddr(wsAddr string)
	SetTCPAddr(tcpAddr string)
	SetNetEvent(event *NetEvent)
	Wait()
}
type IGateManage interface {
	Create() IGate
	Connect(addr string, netEvent *NetEvent, tag int64, autoReconnect bool) IWSClient
	ConnectTcp(addr string, netEvent *NetEvent, tag int64, autoReconnect bool) ITCPClient
}
