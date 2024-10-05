package f

import (
	"time"
	"tzgit.kaixinxiyou.com/engine/hf/enum/enginelog"

	"tzgit.kaixinxiyou.com/engine/hf/actorI"
	"tzgit.kaixinxiyou.com/engine/hf/enum/actorType"
	"tzgit.kaixinxiyou.com/engine/hf/enum/msgType"
	"tzgit.kaixinxiyou.com/engine/hf/m"
	"tzgit.kaixinxiyou.com/engine/message"
)

// 引擎需要的函数

// 获取本服务的中心服的地址 ip+端口
var GetCurrCenterAddr func() string

var GetGroupName func(actorType actorType.ActorType, id uint32) string
var GetGroupIdByName func(name string) uint32

// 消息执行统计
var OnExecutionStatistics func(msgId uint16, useTime int64)

var CallToMessage func(actorType actorType.ActorType, id uint32, message message.IReqMessage) message.IRepMessage
var CallToMessage1 func(serverId uint32, t actorType.ActorType, id uint32, message message.IReqMessage) message.IRepMessage
var RequestDistributionServer func(groupName string, t actorType.ActorType, actorId uint32) uint32
var RequestHeartServer func(t actorType.ActorType, actorId []uint32, serverId uint32)

// 返回指定区在线玩家id列表
var GetOnlineUserIds func(areaId uint32) []uint32

// actor 是否关闭
var ActorIsClose func(actorType actorType.ActorType, actorId uint32) bool

// 获取初始化actor的channel 的数量
var GetChanChannelNum func(actorType actorType.ActorType) int

var GetOffsetTime = func() time.Duration {
	return 0
}

// 获取检测actor是否关闭的版本
var GetActorCloseVersion func() int64

// 发送消息到指定的服务器
var GoSendMessage func(message actorI.IActorMessage) bool

// 发送消息到指定的服务器
var GoSendMessage1 func(serverId uint32, message actorI.IActorMessage) bool

// 发送messgae 到指定的actor里面
var GoSendMessage2 func(t actorType.ActorType, actorId uint32, message message.IMessage) bool

// 发送messgae 到指定的actor里面
var GoSendMessage3 func(serverId uint32, t actorType.ActorType, actorId uint32, message message.IMessage) bool
var GoSendMessage4 func(t actorType.ActorType, actorId uint32, message actorI.IActorMessage) bool

var SendToPlayerRouteMsgData func(userId uint32, t actorType.ActorType, id uint32, rpcId uint32, message []byte)
var SendToMessageData func(actorType actorType.ActorType, id uint32, rpcId uint32, message []byte, messageType msgType.MsgType)
var SendForwardPlayer func(id uint32, message message.IMessage)
var SendToPlayer func(id uint32, message message.IMessage)
var BroadcastForwardPlayer func(userIds []uint32, message message.IMessage, data []byte)

// 广播指定区的在线用户(直接发送给网关agent)
var BroadcastArea func(areaId uint32, pck message.IMessage)

// 广播指定战区的在线用户(直接发送给网关agent)
var BroadcastFightArea func(fightAreaId uint32, pck message.IMessage)
var GetPlayer func(userId uint32) actorI.ISimplePlayer

// 广播指定用户(发送到PlayerActor)
var BroadcastPlayerActor func(userIds []uint32, message message.IMessage, data []byte)

// 重新补发消息
var RetrySendMessage func(message actorI.IActorMessage)

// 清理actor所在的server缓存缓存
var CallClearServerId func(actorType actorType.ActorType, id uint32)

// 获取服务器编号
var GetServerId func(_actorType actorType.ActorType, actorId uint32) uint32
var NewActorMessage1 func([]byte) (actorI.IActorMessage, error)
var NewActorMessage func() actorI.IActorMessage
var PutActorMessage func(actorI.IActorMessage) // 回收actorMessage

var GetAllAreaIdByFightAreaId func(fightAreaId uint32) []uint32

var SendActorNum func(serverId uint32, actorType actorType.ActorType, num int32)

var GetNowTime func() int64
var GetNow func() time.Time

// 获取机器的时间，一般不要使用
var GetSystemUnix func() int64

// 返回系统的day ，比如2023-01-06 返回 6
var GetDay func() int

var GetScene func(actorType.ActorType) actorI.IScene
var OnStart func()
var OnCloseBefore func()
var GetStatusInfo func() []byte
var GetTotalChanNum func() int32
var ExecFuncByServer func(name, param1, param2, param3 string)

// 内置playerAgent
var CheckLogin func(m.IC2SLogin) int
var CheckPckIdIsOpen func(uint16) bool
var CheckPckIdFlag func(uint16) (isReq bool, Route int32)
var CheckRoute func(routeId int32, d []byte, actorID uint32, loginAreaId uint32) (actorType.ActorType, uint32)
var CheckMaintain func(loginAreaId uint32) (isMaintain bool)
var CheckUserID func(AccountID int64, AreaID uint32) (userid uint32)

// 内置serverActor
var CheckServerConnectAddr func(actorID uint32) string
var CheckServerListenAddr func() string

// 获取服务器配置的的版本
var GetServerConfigVer func() int64

// 获取伴生携程数据
var GetAssociatedActor func() map[actorType.ActorType]actorType.ActorType

// 上报执行时间
var SendExecTime func(funcName string, execTime int64, params ...interface{})

// 检测是否要记录log
var CheckLog func(enginelog.EngineLog) bool

// 增加包的数量统计
var AddPckNum func(pckId int32)

//请求分配服务器
