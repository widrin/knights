package actorI

import (
	"sync"
	"tzgit.kaixinxiyou.com/engine/hf/engineI"
	"tzgit.kaixinxiyou.com/engine/hf/enum/actorState"
	"tzgit.kaixinxiyou.com/engine/hf/enum/actorType"
	"tzgit.kaixinxiyou.com/engine/hf/enum/eventcode"
	"tzgit.kaixinxiyou.com/engine/hf/enum/msgType"
	"tzgit.kaixinxiyou.com/engine/hf/enum/tickType"
	"tzgit.kaixinxiyou.com/engine/message"
)

type IActor interface {
	GetId() uint32
	Init(self IActor) bool
	SetRegisterExt(d map[uint16]interface{}) //设置注册的消息
	GetName() string
	GetType() actorType.ActorType
	GetSelf() IActor
	GetState() actorState.ActorState
	SetState(s actorState.ActorState)
	SetChannelNum(n int)
	OnDayChange()
	OnSecondChange()
	CheckClose()
	CloseByScene()
	OnMinuteChange()
	OnWeekChange()
	OnSecond60Change()
	Update(diff int32)
	OnReceiveMessageAfter(_ IActorMessage, tag int32)
	OnReceiveMessageBefore(_ IActorMessage)
	Go(id uint16, args ...interface{})
	//只传1个参数的接口
	Go1(id uint16, args interface{})
	GoMsg(message message.IMessage)
	Call(id uint16, args ...interface{}) interface{}
	CallMsg(message message.IReqMessage) message.IRepMessage
	Register(msgId uint16, f interface{})
	RegisterEventCode(msgId eventcode.EventCode, f interface{})
	OnDestroy()
	Start()
	OnInInit()
	OnHandleMessage(args []interface{})
	SetScene(scene IScene)
	CallToMessage(actorType actorType.ActorType, id uint32, message message.IReqMessage) message.IRepMessage
	ReadToMessage(actorType actorType.ActorType, id uint32, message message.IReqMessage) message.IRepMessage
	GoNetMsg(t msgType.MsgType, data []byte)
	OnCloseBefore()
	Close(closeWg *sync.WaitGroup)
	GetFreeTime() int64
	SetFreeTime(freeTime int64)
	GetServer() engineI.IServer
	GetAction(id uint16) (interface{}, bool)
	SetTagName(tags ...string)
	// OnHookMessage actor hook 消息处理，返回出true，说明hook成功 则不继续执行，
	OnHookMessage(message IActorMessage) bool
	SetActionCapacity(cap uint16) //设置actor的action容量
	GetAsyncWait() *sync.WaitGroup
}
type IScene interface {
	IActor
	GetActor(actorId uint32) IActor
	OnInitOk()
	SetCheck(check bool)
	GetActorNum() int32
	GetChanNum() int32
	SetTick(tick tickType.TickType)
	CloseScene()
	GoRemoveActor(actorId uint32)
}
type ActorRetInfo struct {
	Ret interface{}
	Err error
}
type IActorMessage interface {
	GetMsgType() msgType.MsgType
	SetMsgType(t msgType.MsgType)
	SetBeforeMsgType(t msgType.MsgType)
	GetBeforeMsgType() msgType.MsgType
	SetMessageType(t int32)
	SetMessage(message message.IMessage)
	SetMessageData(data []byte)
	GetMessageType() int32
	SetRpcId(rpcId uint32)
	SetUserId(userId uint32)
	SetTargetSceneId(actorType actorType.ActorType)
	SetTargetActorId(actorId uint32)
	GetTargetSceneId() actorType.ActorType
	GetTargetActorId() uint32
	GetMsgOnlyId() uint64
	GetServerId() uint32
	SetServerId(serverId uint32)
	GetMessageId() uint16
	GetMessage() message.IMessage
	GetRpcId() uint32
	GetMessageData() []byte
	UnmarshalBytes() error
	MarshalBytes()
	ToBytes(data []byte) []byte
	ToBytes1() []byte
	GetUserIds() []uint32
	CopyMessageData(data []byte)
	GetUserId() uint32
	AddUserId(userId uint32)
	GetCreateServerId() uint32         //获取创建的serverId
	SetCreateServerId(serverId uint32) //设置创建的serverId，目前主要用于路由消息
	GetServerRpcId() int64
	SetServerRpcId(rpcId int64)
	Reset()
}

type IServerActor interface {
	IActor
	GoSendMessage(message IActorMessage) bool

	GoReceiveRepMsg(message IActorMessage)
}

type ICenterDiscoveryScene interface {
	CallGetServerId(actorType actorType.ActorType, actorId uint32) uint32
	CallGetCenterServerId(actorType actorType.ActorType) uint32
	CallSetServerId(actorType actorType.ActorType, actorId uint32, serverId uint32)
	CallClearServerId(actorType actorType.ActorType, actorId uint32)

	CallGetServerIdByGroup(group string, actorType actorType.ActorType, actorId uint32) uint32
	CallGetCenterServerIdByGroup(group string, actorType actorType.ActorType) uint32
	RequestDistributionServer(group string, actorType actorType.ActorType, actorId uint32) uint32
	RequestHeartServer(actorType actorType.ActorType, actorId []uint32, serverId uint32)
}
type ICenterActor interface {
	IActor
	GetAllServer() []IServerData
	GetServerData(serverId uint32) IServerData
	GoServerOnline(serverId uint32)
	GoServerOffline(serverId uint32)
}
