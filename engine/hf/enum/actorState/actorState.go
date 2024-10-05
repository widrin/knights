package actorState

type ActorState int32

const (
	Normal    ActorState = iota //正常
	NeedClose                   //需要关闭
	Closeing                    //关闭中
	Closed                      //关闭完成
)
