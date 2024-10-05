package m

import (
	"tzgit.kaixinxiyou.com/engine/hf/gateI"
	"tzgit.kaixinxiyou.com/engine/message"
)

var Processor gateI.IProcessor
var NewS2SDistributionRuleServerReq func() IS2SDistributionRuleServerReq
var NewS2SDistributionRuleServerRep func() IS2SDistributionRuleServerRep

type IS2SDistributionRuleServerReq interface {
	message.IReqMessage

	SetActorType(actorType uint32)
	GetActorType() uint32
	SetActorId(actorId uint32)
	GetActorId() uint32
}

var PCK_S2SDistributionRuleServerReq_ID uint16
var PCK_S2SClearRuleCacheReq_ID uint16
var PCK_S2SActorHeart_ID uint16

type IS2SDistributionRuleServerRep interface {
	message.IRepMessage

	GetTag() int32
	GetServerId() uint32
	SetServerId(serverId uint32)
}

var NewS2SActorHeart func() IS2SActorHeart

// actdor的心跳
type IS2SActorHeart interface {
	message.IMessage

	SetActorType(actorType uint32)
	GetActorType() uint32
	SetServerId(serverId uint32)
	GetServerId() uint32
	SetActorIds(ids []uint32)
	GetActorIds() []uint32
}

type IS2CRep interface {
	message.IRepMessage
	SetTag(tag int32)
}

var NewS2CRep func() IS2CRep

type IS2SClearRuleCacheReq interface {
	message.IRepMessage
	GetServerId() uint32
}

// 内置server,playerAgent

var NewS2CErrorRpc func() IS2CErrorRpc

type IS2CErrorRpc interface {
	message.IRepMessage
	SetTag(uint32)
}

var NewG2SPlayerOnlineState func() IG2SPlayerOnlineState

type IG2SPlayerOnlineState interface {
	message.IMessage

	SetState(int32)
	SetLoginAreaId(uint32)
}

var NewS2CLogin func() IS2CLogin

type IS2CLogin interface {
	message.IRepMessage
	SetUserId(int64)
	SetTag(uint32)
}

var NewS2CReLogin func() IS2CReLogin

type IS2CReLogin interface {
	message.IRepMessage
	SetTag(uint32)
}

var NewC2SLogin func() IC2SLogin

type IC2SLogin interface {
	message.IMessage

	GetAccountID() int64
	GetToken() string
	GetAreaID() uint32
}

var NewC2SReLogin func() IC2SReLogin

type IC2SReLogin interface {
	message.IMessage

	GetToken() string
	GetUserId() uint32
}

var NewS2SKill func() IS2SKill

type IS2SKill interface {
	message.IMessage

	GetIsFCM() uint32
	GetFcmText() string
	GetTag() uint32
}

var NewS2CKill func() IS2CKill

type IS2CKill interface {
	message.IRepMessage
	SetTag(uint32)
	SetKillType(uint32)
	SetKillMsg(string)
}

var (
	SYS_成功        int
	SYS_Actor不存在  int
	SYS_重连失败      int
	SYS_超时        int
	SYS_帐号在其他地方登录 int
	SYS_服务器维护中    int
	SYS_TOKEN错误   int
	SYS_您被踢下线     int
	SYS_服务器内部错误   int
)

var PCK_C2SLogin_ID uint16
var PCK_C2SReLogin_ID uint16
var PCK_S2SKill_ID uint16
var PCK_C2SLoginEnd_ID uint16
