package engineI

type RetInfo struct {
	Ret interface{}
	Err error
	Cb  interface{}
}

type ICallInfo interface {
	Reset() //重置，回收到对象池前调用
}
type IServer interface {
	SetIsStart(b bool)
	GetIsStart() bool
	SetName(name string)
	GetName() string
	Close()
	GetChanCall() chan ICallInfo
	GetReadCall() chan ICallInfo
	Exec(ci ICallInfo)
	Register(id interface{}, f interface{})
	Call(id interface{}, args ...interface{}) (interface{}, error)
	Go(id interface{}, args ...interface{})
	Go1(id interface{}, args interface{})
	GoRead(id interface{}, args []byte)
	GoBytes(id interface{}, args []byte)
	GetChanNum() int
	SetCurrMsgId(msgId uint16) //设置当前处理的消息编号，输出日志用
}
type IClient interface {
	Idle() bool
	Close()
	GetChanAsyncRet() chan *RetInfo
	Cb(ri *RetInfo)
}
type ISkeletonManage interface {
	Create(rpcServer IServer) ISkeleton
}
type ISkeleton interface {
	Init()
	Run(closeSig chan bool)
	RegisterChanRPC(id interface{}, f interface{})

	SetServerName(n string)
	GetServerName() string
	GetServer() IServer
}
type IChanRpcManage interface {
	NewServer(l int) IServer
	NewClient(l int) IClient
}
