package actorType

type ActorType int32

var GetScenePlayerAgent func() ActorType
var GetSceneCenter func() ActorType
var GetSceneCenterDiscovery func() ActorType
var GetScenePlayer func() ActorType
var GetSceneArea func() ActorType
var GetSceneServer func() ActorType

var GetSceneNil func() ActorType
var GetSceneActorNoExist func() ActorType
