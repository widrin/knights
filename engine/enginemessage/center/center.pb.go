package center

import (
	"encoding/json"
	"tzgit.kaixinxiyou.com/engine/message"
	"tzgit.kaixinxiyou.com/engine/message/msglist"
)

func init() {
	msglist.Register(PCK_S2SActorHeart_ID, func() message.IMessage { return new(S2SActorHeart) }, 0)
	msglist.Register(PCK_S2SDistributionRuleServerReq_ID, func() message.IMessage { return new(S2SDistributionRuleServerReq) }, 0)
	msglist.Register(PCK_S2SDistributionRuleServerRep_ID, func() message.IMessage { return new(S2SDistributionRuleServerRep) }, 0)
	msglist.Register(PCK_S2SClearRuleCacheReq_ID, func() message.IMessage { return new(S2SClearRuleCacheReq) }, 0)
	msglist.Register(PCK_S2SDistributionCenterServerReq_ID, func() message.IMessage { return new(S2SDistributionCenterServerReq) }, 0)
	msglist.Register(PCK_S2SResetCenterServerReq_ID, func() message.IMessage { return new(S2SResetCenterServerReq) }, 0)
	msglist.Register(PCK_S2SCloseCenterActor_ID, func() message.IMessage { return new(S2SCloseCenterActor) }, 0)
	msglist.Register(PCK_S2SCloseCenterActorSuccess_ID, func() message.IMessage { return new(S2SCloseCenterActorSuccess) }, 0)
	msglist.Register(PCK_S2SDistributionCenterServerRep_ID, func() message.IMessage { return new(S2SDistributionCenterServerRep) }, 0)
	msglist.Register(PCK_S2SHeartCenterServerReq_ID, func() message.IMessage { return new(S2SHeartCenterServerReq) }, 0)
}

const PCK_S2SActorHeart_ID = 5002 //活动的actor心跳
// 活动的actor心跳
type S2SActorHeart struct {
	//集群编号
	GroupId uint32 `protobuf:"varint,1,opt,name=GroupId,proto3" json:"GroupId,omitempty"`
	//ActorType
	ActorType uint32 `protobuf:"varint,2,opt,name=ActorType,proto3" json:"ActorType,omitempty"`
	//ServerId
	ServerId uint32 `protobuf:"varint,3,opt,name=ServerId,proto3" json:"ServerId,omitempty"`
	//actorId
	ActorIds []uint32 `protobuf:"varint,4,rep,name=ActorIds,proto3" json:"ActorIds,omitempty"`
}

func (m *S2SActorHeart) String() string { a, _ := json.Marshal(m); return string(a) }
func (*S2SActorHeart) ProtoMessage()    {}
func (m *S2SActorHeart) GetId() uint16  { return PCK_S2SActorHeart_ID }
func CreateInitS2SActorHeart() *S2SActorHeart {
	info := &S2SActorHeart{}
	info.ActorIds = make([]uint32, 0, 8)
	return info
}

const PCK_S2SDistributionRuleServerReq_ID = 5003 //获取分配的serverId
// 获取分配的serverId
type S2SDistributionRuleServerReq struct {
	//actorId
	ActorId uint32 `protobuf:"varint,1,opt,name=ActorId,proto3" json:"ActorId,omitempty"`
	//ActorType
	ActorType uint32 `protobuf:"varint,2,opt,name=ActorType,proto3" json:"ActorType,omitempty"`
	//集群编号
	GroupId uint32 `protobuf:"varint,3,opt,name=GroupId,proto3" json:"GroupId,omitempty"`
}

func (m *S2SDistributionRuleServerReq) String() string { a, _ := json.Marshal(m); return string(a) }
func (*S2SDistributionRuleServerReq) ProtoMessage()    {}
func (m *S2SDistributionRuleServerReq) GetId() uint16  { return PCK_S2SDistributionRuleServerReq_ID }
func CreateInitS2SDistributionRuleServerReq() *S2SDistributionRuleServerReq {
	info := &S2SDistributionRuleServerReq{}
	return info
}

const PCK_S2SDistributionRuleServerRep_ID = 5004 //获取分配的serverId
// 获取分配的serverId
type S2SDistributionRuleServerRep struct {
	//结果
	Tag int32 `protobuf:"varint,1,opt,name=Tag,proto3" json:"Tag,omitempty"`
	//服务器编号
	ServerId uint32 `protobuf:"varint,2,opt,name=ServerId,proto3" json:"ServerId,omitempty"`
}

func (m *S2SDistributionRuleServerRep) GetTag() int32 {
	return m.Tag
}
func (m *S2SDistributionRuleServerRep) String() string { a, _ := json.Marshal(m); return string(a) }
func (*S2SDistributionRuleServerRep) ProtoMessage()    {}
func (m *S2SDistributionRuleServerRep) GetId() uint16  { return PCK_S2SDistributionRuleServerRep_ID }
func CreateInitS2SDistributionRuleServerRep() *S2SDistributionRuleServerRep {
	info := &S2SDistributionRuleServerRep{}
	return info
}

const PCK_S2SClearRuleCacheReq_ID = 5005 //获取分配的serverId
// 获取分配的serverId
type S2SClearRuleCacheReq struct {
	//清理缓存rule，后台用
	ServerId uint32 `protobuf:"varint,1,opt,name=ServerId,proto3" json:"ServerId,omitempty"`
}

func (m *S2SClearRuleCacheReq) String() string { a, _ := json.Marshal(m); return string(a) }
func (*S2SClearRuleCacheReq) ProtoMessage()    {}
func (m *S2SClearRuleCacheReq) GetId() uint16  { return PCK_S2SClearRuleCacheReq_ID }
func CreateInitS2SClearRuleCacheReq() *S2SClearRuleCacheReq {
	info := &S2SClearRuleCacheReq{}
	return info
}

const PCK_S2SDistributionCenterServerReq_ID = 5006 //获取中心服所在的地址
// 获取中心服所在的地址
type S2SDistributionCenterServerReq struct {
	//groupId*100000+actorType
	ActorId uint32 `protobuf:"varint,1,opt,name=ActorId,proto3" json:"ActorId,omitempty"`
}

func (m *S2SDistributionCenterServerReq) String() string { a, _ := json.Marshal(m); return string(a) }
func (*S2SDistributionCenterServerReq) ProtoMessage()    {}
func (m *S2SDistributionCenterServerReq) GetId() uint16  { return PCK_S2SDistributionCenterServerReq_ID }
func CreateInitS2SDistributionCenterServerReq() *S2SDistributionCenterServerReq {
	info := &S2SDistributionCenterServerReq{}
	return info
}

const PCK_S2SResetCenterServerReq_ID = 5009 //重新分配中心的serverId
// 重新分配中心的serverId
type S2SResetCenterServerReq struct {
	//groupId*100000+actorType
	ActorId uint32 `protobuf:"varint,1,opt,name=ActorId,proto3" json:"ActorId,omitempty"`
}

func (m *S2SResetCenterServerReq) String() string { a, _ := json.Marshal(m); return string(a) }
func (*S2SResetCenterServerReq) ProtoMessage()    {}
func (m *S2SResetCenterServerReq) GetId() uint16  { return PCK_S2SResetCenterServerReq_ID }
func CreateInitS2SResetCenterServerReq() *S2SResetCenterServerReq {
	info := &S2SResetCenterServerReq{}
	return info
}

const PCK_S2SCloseCenterActor_ID = 5010 //关闭Actor
// 关闭Actor
type S2SCloseCenterActor struct {
	//groupId*100000+actorType
	ActorId uint32 `protobuf:"varint,1,opt,name=ActorId,proto3" json:"ActorId,omitempty"`
}

func (m *S2SCloseCenterActor) String() string { a, _ := json.Marshal(m); return string(a) }
func (*S2SCloseCenterActor) ProtoMessage()    {}
func (m *S2SCloseCenterActor) GetId() uint16  { return PCK_S2SCloseCenterActor_ID }
func CreateInitS2SCloseCenterActor() *S2SCloseCenterActor {
	info := &S2SCloseCenterActor{}
	return info
}

const PCK_S2SCloseCenterActorSuccess_ID = 5011 //关闭Actor成功
// 关闭Actor成功
type S2SCloseCenterActorSuccess struct {
	//groupId*100000+actorType
	ActorId uint32 `protobuf:"varint,1,opt,name=ActorId,proto3" json:"ActorId,omitempty"`
}

func (m *S2SCloseCenterActorSuccess) String() string { a, _ := json.Marshal(m); return string(a) }
func (*S2SCloseCenterActorSuccess) ProtoMessage()    {}
func (m *S2SCloseCenterActorSuccess) GetId() uint16  { return PCK_S2SCloseCenterActorSuccess_ID }
func CreateInitS2SCloseCenterActorSuccess() *S2SCloseCenterActorSuccess {
	info := &S2SCloseCenterActorSuccess{}
	return info
}

const PCK_S2SDistributionCenterServerRep_ID = 5007 //获取中心服所在的地址
// 获取中心服所在的地址
type S2SDistributionCenterServerRep struct {
	//结果
	Tag int32 `protobuf:"varint,1,opt,name=Tag,proto3" json:"Tag,omitempty"`
	//服务器编号
	ServerId uint32 `protobuf:"varint,2,opt,name=ServerId,proto3" json:"ServerId,omitempty"`
}

func (m *S2SDistributionCenterServerRep) GetTag() int32 {
	return m.Tag
}
func (m *S2SDistributionCenterServerRep) String() string { a, _ := json.Marshal(m); return string(a) }
func (*S2SDistributionCenterServerRep) ProtoMessage()    {}
func (m *S2SDistributionCenterServerRep) GetId() uint16  { return PCK_S2SDistributionCenterServerRep_ID }
func CreateInitS2SDistributionCenterServerRep() *S2SDistributionCenterServerRep {
	info := &S2SDistributionCenterServerRep{}
	return info
}

const PCK_S2SHeartCenterServerReq_ID = 5008 //活动的中心服心跳
// 活动的中心服心跳
type S2SHeartCenterServerReq struct {
	//集群编号
	GroupId uint32 `protobuf:"varint,1,opt,name=GroupId,proto3" json:"GroupId,omitempty"`
	//ActorType
	ActorType uint32 `protobuf:"varint,2,opt,name=ActorType,proto3" json:"ActorType,omitempty"`
	//ServerId
	ServerId uint32 `protobuf:"varint,3,opt,name=ServerId,proto3" json:"ServerId,omitempty"`
	//actorId
	ActorIds []uint32 `protobuf:"varint,4,rep,name=ActorIds,proto3" json:"ActorIds,omitempty"`
}

func (m *S2SHeartCenterServerReq) String() string { a, _ := json.Marshal(m); return string(a) }
func (*S2SHeartCenterServerReq) ProtoMessage()    {}
func (m *S2SHeartCenterServerReq) GetId() uint16  { return PCK_S2SHeartCenterServerReq_ID }
func CreateInitS2SHeartCenterServerReq() *S2SHeartCenterServerReq {
	info := &S2SHeartCenterServerReq{}
	info.ActorIds = make([]uint32, 0, 8)
	return info
}
