package actorI

import "tzgit.kaixinxiyou.com/engine/hf/enum/actorType"

type ISceneManage interface {
	RegisterScene(scene IScene)
	GetScene(actorType actorType.ActorType) IScene
	RegisterCreateActor(t actorType.ActorType, f func(id uint32) IActor)
	Start()
	RegisterSystemScene(t actorType.ActorType)
	GetIsClose() bool
}
