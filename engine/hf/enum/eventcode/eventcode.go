package eventcode

type EventCode int32

const (
	ActorMsgHandle EventCode = 50001 + iota //正常
	ActorGoNetMessage
	ActorMsgHandleRet
	ServerActorGetReqMessage //获取请求的消息

	TickChange
	SecondChange
	MinuteChange
	HourChange
	DayChange
	WeekChange
	Second60Change //60秒事件

	InitOk //初始化完成
	ServerOnline
	ServerOffline

	// 内置server,playerAgent

	PlayerAgentS2A_GoNetMessage
	PlayerAgentS2A_ChangeAgent
	PlayerAgentS2A_AgentDestroy
	PlayerAgentS2A_SendCachePck

	PlayerAgentScene_AgentInit

	ServerActor_SendMessage
	ServerActor_AgentInit
	ServerActor_GoReceiveRepMsg
	ServerActor_AgentDestroy //客户端断开连接

	NilActor_NotExist //处理不存在的actor
	ActorNoExistSend  //处理不存在的actor
	RemoveActor       //删除已经关闭的Actor
	Scene_CloseActor  //场景关闭Actor
)

var (
	PlayerAgentActor_Forward uint16 //直接转发给客户端。使用引擎的server,playerAgent,noExist需要设置该值
)
