package tickType

type TickType int32

const (
	Update   TickType = 1  //更新，1秒10次
	Second   TickType = 2  //秒tick
	Minute   TickType = 4  //分钟tick
	Hour     TickType = 8  //小时tick
	Day      TickType = 16 //天tick
	Week     TickType = 32 //周tick
	Second60 TickType = 64 //60tick

	ALL TickType = Update | Second | Minute | Second60 | Hour | Day | Week
)
