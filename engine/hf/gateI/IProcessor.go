package gateI

import "tzgit.kaixinxiyou.com/engine/message"

type IProcessor interface {
	ReadId(data []byte) (id uint16)
	Register(id uint16, f func() message.IMessage)
	// must goroutine safe
	Unmarshal(data []byte) (message.IMessage, error)
	// must goroutine safe
	Marshal(msg message.IMessage) ([]byte, error)
	GetMsgType() MsgType
	ReadTag(number int32, b []byte) int64
	GetMsg(id uint16) message.IMessage
}
type IProcessorManage interface {
	Create() IProcessor
}
