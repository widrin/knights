package msgType

type MsgType int32

const (
	Default                     MsgType = 0
	Reply                               = 1  //回复消息
	PlayerActor                         = 2  //客户端发送给 playerActor 发送的消息
	ForwardPlayer                       = 3  //直接回复给玩家
	ActorNOExist                        = 4  //发送失败，actor 不存在
	Read                                = 6  //只读消息
	BroadcastForwardPlayer              = 7  //广播给一批玩家
	BroadcastOnlinePlayer               = 8  //中转到区管理，然后广播在线用户
	BroadcastOnlinePlayerActor          = 9  //中转到区管理，然后广播在线用户Actor
	RouteMessage                        = 10 //player路由消息
	BroadcastForwardPlayerActor         = 11 //广播给一批玩家actor

)
