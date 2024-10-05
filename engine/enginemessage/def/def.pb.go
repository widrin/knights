package def

import (
	"encoding/json"
	"tzgit.kaixinxiyou.com/engine/message"
	"tzgit.kaixinxiyou.com/engine/message/msglist"
)

func init() {
	msglist.Register(PCK_S2CRep_ID, func() message.IMessage { return new(S2CRep) }, 0)
}

const PCK_S2CRep_ID = 5001 //通用的，没有额外数据的rep返回
// 通用的，没有额外数据的rep返回
type S2CRep struct {
	//返回结果
	Tag int32 `protobuf:"varint,1,opt,name=Tag,proto3" json:"Tag,omitempty"`
}

func (m *S2CRep) GetTag() int32 {
	return m.Tag
}
func (m *S2CRep) String() string { a, _ := json.Marshal(m); return string(a) }
func (*S2CRep) ProtoMessage()    {}
func (m *S2CRep) GetId() uint16  { return PCK_S2CRep_ID }
func CreateInitS2CRep() *S2CRep {
	info := &S2CRep{}
	return info
}
