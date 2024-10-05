package workid

import (
	_redis "github.com/go-redis/redis/v8"
	"net/http"
	"time"
	"tzgit.kaixinxiyou.com/utils/common/convert"
	"tzgit.kaixinxiyou.com/utils/common/log"
	"tzgit.kaixinxiyou.com/utils/common/util"
	"unsafe"
)

const (
	RedisKey     = "tzid_service"
	RedisWorkKey = "tzid_service:work" //workId 对应的 key
	RedisSortKey = "tzid_service:sort"
)

func Create() *WorkId {
	r := &WorkId{}
	return r
}

type WorkId struct {
}

var Redis _redis.Cmdable

//自动去获取MACHINEid
func (w *WorkId) Init(redisAddr string, redisUser string, redisPass string, redisDb int) {

	cfg := &_redis.Options{
		Addr:     redisAddr,
		Username: redisUser,
		Password: redisPass, // no password set
		DB:       redisDb,   // use default DB
	}
	cfg.ReadTimeout = 500 * time.Millisecond
	cfg.WriteTimeout = 500 * time.Millisecond
	//cfg.MinIdleConns=10
	cfg.MaxConnAge = 30 * time.Minute
	cfg.PoolSize = 64
	cfg.DB = int(0)
	Redis = _redis.NewClient(cfg)
}
func (w *WorkId) Start() error {
	mux := http.NewServeMux()
	httpServer := &http.Server{Addr: "0.0.0.0:999", Handler: mux} //修改为0.0.0.0 为了统计服务器状态
	mux.HandleFunc("/get_work_id", w.GetWorkId)
	return httpServer.ListenAndServe()

}

func (w *WorkId) saveData() {
}

var script = `
	local timeOut=60*60*6
	local workId = redis.call("hget", KEYS[1],ARGV[1])
	if not workId  then
		workId = redis.call("zcard", KEYS[2])+1
		if workId > 16383 then --查找最新的一个
			workId= redis.call("zrange", KEYS[2],0,0)[1]
			--判断是否超过24小时，没有，则返回错误，表示没有分配的id
			local oldTime = redis.call("zscore", KEYS[2],workId)
			if ARGV[2]-oldTime<timeOut then
				return -1
			end
			--删除原来的对应key
			local oldKey =redis.call("hget", KEYS[3],workId)			
			redis.call("hdel", KEYS[1],oldKey)
		end
		redis.call("hset", KEYS[1],ARGV[1],workId)		
		redis.call("hset", KEYS[3],workId,ARGV[1])
	end
	redis.call("zadd" ,KEYS[2] ,ARGV[2], workId) -- 增加时间
	return workId
`

func (w *WorkId) GetWorkId(rep http.ResponseWriter, req *http.Request) {
	now := time.Now().Unix()
	err := req.ParseForm()
	if err != nil {
		rep.WriteHeader(500)
		rep.Write(StringToBytes(err.Error()))
		return
	}
	reqMap := req.Form
	serverKey := reqMap.Get("server_key")
	if serverKey == "" {
		log.Debug("server_key不能为空")
		rep.WriteHeader(500)
		rep.Write(StringToBytes("server_key不能为空"))
		return
	}
	ip := util.ClientIP(req)
	serverKey = ip + "#" + serverKey

	c := Redis.Eval(req.Context(), script, []string{RedisKey, RedisSortKey, RedisWorkKey}, serverKey, convert.ToString(now))
	if c.Err() != nil {
		log.Debug("GetWorkId:%v  %v", serverKey, c.Err().Error())
		rep.WriteHeader(500)
		rep.Write(StringToBytes(c.Err().Error()))
	} else {
		log.Debug("GetWorkId:%v  %v", serverKey, c.Val())
		rep.Write(StringToBytes(convert.ToString(c.Val())))
	}
}

func StringToBytes(data string) []byte {
	return *(*[]byte)(unsafe.Pointer(&data))
}
