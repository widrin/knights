package msglist

import "tzgit.kaixinxiyou.com/engine/message"

var msgInfo = make(map[uint16]func() message.IMessage)
var pckLevel = make(map[uint16]int32)

func Register(id uint16, f func() message.IMessage, level int32) {
	msgInfo[id] = f
	if level > 0 {
		pckLevel[id] = level
	}
}

func GetLen() int32 {
	return int32(len(msgInfo))
}
func GetLevel(msgId uint16) int32 {
	return pckLevel[msgId]
}
func GetMsgList() map[uint16]func() message.IMessage {
	return msgInfo
}
func ClearData() {
	msgInfo = nil
}
