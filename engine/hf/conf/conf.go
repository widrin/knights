package conf

import (
	"strconv"
	"time"
	"tzgit.kaixinxiyou.com/engine/hf/enum/actorType"
)

var LogicalFrameTime = time.Millisecond * 33 //1秒10帧 毫秒
var DebugAddr = ""                           //性能调试的地址
var CenterAddr = ""                          //中心服地址
var Centers []int32                          //中心服务器监控的类型
var ServerId = uint32(0)                     //服务器编号
var AppId = uint32(0)                        //appId
var OpenDebugMessage = false                 //是否开启追踪消息
var WarnExecTimeInterval = int64(0)          //监控从发送消息到开始执行的间隔时间，超过这个值，则报警
var WarnExecTime = int64(0)                  //监控消息的执行时间，超过这个值，则报警
var WsMaxConnNum = 10000                     //websocket最大连接数
var WSAddr = "0.0.0.0:0"                     //websocket监听端口
var ServerFreeTime = 60                      //服务器场景的FreeTime(单位秒)

// 服务器组
var WSGroup = ""
var WSGroupId = uint32(0)   //服务器组编号
var IsService = true        //是否服务器服务
var GetActorRetryTimes = 10 //不是服务器的actor，为空的时候重试的次数
func GetServerOnlineServerKey() string {
	return WSGroup + ":server:online" //online:actorType
}
func GetServerOnlineServerField() string {
	return strconv.Itoa(int(ServerId)) //online:actorType
}

// 中心服务器在线
func GetCenterOnlineKey() string {
	return "All:center:online" //online:actorType
}

// 获取分配数据
func GetCenterRuleKey(actorType actorType.ActorType) string {
	return WSGroup + ":center:rule:" + strconv.Itoa(int(actorType))
}

// 场景负载人数
func GetLoadNumKey(actorType actorType.ActorType) string {
	return WSGroup + ":load:" + strconv.Itoa(int(actorType)) //online:actorType
}
