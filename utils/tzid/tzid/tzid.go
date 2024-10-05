package tzid

import (
	"fmt"
	"log"
	"runtime"
	"sync/atomic"
	"time"
)

var haveInit = false
var currIdx = int64(0)

const (
	StartTime = 1640966400000 //2022-1-1 00:00:00  毫秒
)

func Init(workId int64) {
	//判断WorkID 是否正确
	if workId < 1 || workId > 16383 {
		log.Fatal("workId 范围[1-16383] ", workId)
	}
	currIdx = workId << 49
	currIdx += (time.Now().UnixNano()/1e6 - StartTime) << 10
	haveInit = true
}
func GetNewId() int64 {
	if haveInit == false {
		fmt.Println(stack())
		log.Fatal("没有初始化WorkId")
	}
	return atomic.AddInt64(&currIdx, 1)
}
func UnmarshalId(id int64) (int64, int64, int64) {
	workId := id >> 49
	timestamp := id>>10&0x7fffffffff + StartTime
	idx := id & 0x3ff
	return workId, timestamp, idx
}
func stack() string {
	var buf [2 << 10]byte
	return string(buf[:runtime.Stack(buf[:], true)])
}
