package actorI

type IServerData interface {
	GetNum() uint32
	GetServerId() uint32
	AddNum()
}

//规则
type IRule interface {
	Init(actor ICenterActor)
	GetServerId(actorId uint32) uint32
	OnServerOnline(serverList uint32)
	OnServerOffline(serverList uint32)
	CheckTimeOut()
	ActorHeartServerId(actorId uint32, serverId uint32)
	GetNumByServer(serverId uint32) uint32
}
