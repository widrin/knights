package tzidWorkId

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"tzgit.kaixinxiyou.com/utils/tzid/tzid"
)

const RenewalDuration = time.Hour //续约时间

var serverKey string
var serverUrl = "http://127.0.0.1:999/get_work_id"

func InitWorkId(_serverKey string, _serverUrl ...string) {
	serverKey = _serverKey
	if len(_serverUrl) > 0 {
		serverUrl = _serverUrl[0]
	}
	workId := getWorkId()
	if workId <= 0 {
		log.Fatal("初始化workId 失败")
	}
	tzid.Init(workId)
	//开启续期
	go renewal()
}
func renewal() {
	//续约
	t := time.NewTicker(RenewalDuration)
	for {
		select {

		case <-t.C:
			getWorkId()
		}
	}
}
func getWorkId() int64 {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
		}
	}()
	client := &http.Client{}
	d := make(url.Values)
	d.Add("server_key", serverKey)
	r, e := client.PostForm(serverUrl, d)
	if e != nil {
		log.Println(e)
		return -1
	}
	if r.StatusCode != 200 {
		log.Printf("服务器内部错误：%v\n", r.StatusCode)
		return -1
	}
	a, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Println(e)
		return -1
	}
	v, e := strconv.ParseInt(string(a), 10, 64)
	if e != nil {
		log.Println(e)
		return -1
	}
	return v
}
