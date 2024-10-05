package timeoutWaitGroup

import (
	"fmt"
	"sync"
	"time"
)

func Create() *TimeoutWaitGroup {
	wait := &TimeoutWaitGroup{}
	return wait
}

type TimeoutWaitGroup struct {
	sync.WaitGroup
}

func (t *TimeoutWaitGroup) GetWaitGroup() *sync.WaitGroup {
	return &t.WaitGroup
}

//timeout 毫秒
func (t *TimeoutWaitGroup) Wait(timeout int64) {
	if timeout == 0 {
		t.WaitGroup.Wait()
	} else {
		ch := make(chan bool, 1)
		go func() {
			t.WaitGroup.Wait()
			ch <- true
		}()
		select {
		case <-time.After(time.Millisecond * time.Duration(timeout)):
			fmt.Println("我超时了")

		case <-ch:
			fmt.Println("我结束了")

		}

	}
}
