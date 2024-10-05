package redisSubscribe

import (
	"context"
	"time"
	"tzgit.kaixinxiyou.com/utils/common/log"
	"tzgit.kaixinxiyou.com/utils/common/redis"
)

// redis订阅
func init() {
	go start()
}

const KEY = "redis_subscribe"

func start() {
	t := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-t.C:
			secondChange()
		}
	}
}

var currVer int64

// GetVer 获取版本
func GetVer() int64 {
	return currVer
}

var _context = context.Background()

// IncrVer 增加版本
func IncrVer() {
	if redis.Redis != nil {
		redis.Redis.IncrBy(_context, KEY, 1)
	}
}

// 每秒检测版本
func secondChange() {
	if redis.Redis != nil {
		_currVer, err := redis.Redis.Get(_context, KEY).Int64()
		if err != nil {
			log.Warn("%v", err)
		} else {
			currVer = _currVer
		}
	}
}
