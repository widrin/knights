package hf

import (
	"tzgit.kaixinxiyou.com/engine/hf/actorI"
	"tzgit.kaixinxiyou.com/engine/hf/engineI"
	"tzgit.kaixinxiyou.com/engine/hf/gateI"
	"tzgit.kaixinxiyou.com/engine/hf/tzDataI"
)

type I interface {
	GetGateManage() gateI.IGateManage
	GetProcessorManage() gateI.IProcessorManage
	GetSkeletonManage() engineI.ISkeletonManage
	GetChanRpcManage() engineI.IChanRpcManage
	GetTzDataManage() tzDataI.ITZManage
	GetSceneManage() actorI.ISceneManage
	GetEngineManage() actorI.IEngineManage
}

var Instance I
