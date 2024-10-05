package actorI

import "tzgit.kaixinxiyou.com/engine/hf/enum/actorType"

//引擎管理
type IEngineManage interface {
	CreateActor(t actorType.ActorType, actorId uint32) IActor
	CreateScene(t actorType.ActorType) IScene
	GetActorNumByTag(tag string) int
	GetActorNumList() string
}
