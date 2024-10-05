package inpub

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"tzgit.kaixinxiyou.com/engine/tzplugin/tzppub"
	"tzgit.kaixinxiyou.com/utils/common/log"
	"tzgit.kaixinxiyou.com/utils/common/mysql"
)

const (
	// MAXINT32 int32最大值
	MAXINT32 = 2147483647
)

var (
	WjxAppkey string
)

// SendMailCondition 发放邮件条件
type SendMailCondition int32

const (
	SendMailCondition_不限接收次数           SendMailCondition = 0
	SendMailCondition_单角色每天只能接受一次该ID邮件 SendMailCondition = 1
	SendMailCondition_单角色每周只能接受一次该ID邮件 SendMailCondition = 2
	SendMailCondition_单角色每月只能接受一次该ID邮件 SendMailCondition = 3
	SendMailCondition_单角色永久只能接受一次该ID邮件 SendMailCondition = 4
)

func (smc SendMailCondition) GetString() string {
	switch smc {
	case SendMailCondition_单角色每天只能接受一次该ID邮件:
		return "d" + time.Now().Format("20060102")
	case SendMailCondition_单角色每月只能接受一次该ID邮件:
		return "m" + time.Now().Format("200601")
	case SendMailCondition_单角色每周只能接受一次该ID邮件:
		y, w := time.Now().ISOWeek()
		return fmt.Sprintf("w%04d%02d", y, w)
	case SendMailCondition_单角色永久只能接受一次该ID邮件:
		return "forever"
	}
	return ""
}

var (
	TzpGiftdata = &TzpGiftData{}
	TzpWjxdata  = &TzpWjxData{}

	isGiftInit int32
	isWjxInit  int32
)

// UserCallback 调用到userid的actor
type UserCallback interface {
	//routerid是用来协助路由到userid的，如果不需要可以忽略,param参数需要透传到tzplugin.TzpPlayerCall
	//tzplugin.TzpPlayerCall 必须在player actor中调用
	//无法执行tzplugin.TzpPlayerCall 则返回空字符串
	//执行tzplugin.TzpPlayerCall 后则返回TzpPlayerCall的返回结果
	TzpCallToUserActor(routerid, userid uint32, param1, param2, param3 string) string

	// TzpSendGift 发送奖励
	// useridList 游戏角色ID列表
	// title 邮件标题
	// content 邮件内容
	// items 邮件道具
	TzpSendGift(useridList []uint32, title string, content string, items []*tzppub.ItemInfo)

	TzpGetPlatformId(userid uint32) int64 //获取玩家的平台ID
	TzpGetAccountId(userid uint32) uint32 //获取玩家的账号ID
	TzpGetChannelId(userid uint32) uint32 //获取玩家的登录渠道ID
}

// IPlayer 玩家基础接口
type IPlayer interface {
	TzpGetAreaOpenDay() int32   //获取开服天数
	TzpGetLevel() int32         //获取玩家等级
	TzpGetChargeValue() int32   //获取总充值金额（分）
	TzpGetUserID() uint32       //获取玩家角色ID
	TzpGetPlatformId() int64    //获取玩家的平台ID
	TzpGetAccountId() uint32    //获取玩家的账号ID
	TzpGetChannelId() uint32    //获取玩家的登录渠道ID
	TzpGetWjxData() interface{} //获取wjx.CreatePlayerWjxData创建出来的指针
}

// TzpGift 天纵插件库礼包数据
type TzpGift struct {
	Id        int32
	ChannelId map[uint32]bool
	Name      string
	Condition SendMailCondition
	StartTime int64
	EndTime   int64
	Title     string
	Content   string
	Status    int32
	Attach    []*tzppub.ItemInfo
	Creator   string
}

// TzpWjx 天纵插件问卷星数据
type TzpWjx struct {
	Wjid       int32
	ChannelId  map[uint32]bool
	Reward     int32
	Url        string
	Content    string
	Name       string
	MinLev     int32
	MaxLev     int32
	MinPay     int32
	MaxPay     int32
	MinOpenDay int32
	MaxOpenDay int32
	StartTime  int64
	EndTime    int64
	Status     int32
	Creator    string
	Version    int64
}

// TzpWjxData 问卷星数据
type TzpWjxData struct {
	Wjxs []*TzpWjx
}

// TzpGiftData 礼包数据
type TzpGiftData struct {
	Gifts []*TzpGift
}

// Get 获取礼包数据
func (td *TzpGiftData) Get(Id int32) *TzpGift {
	for _, v := range td.Gifts {
		if v.Id == Id {
			return v
		}
	}
	return nil
}

// Get 获取问卷数据
func (td *TzpWjxData) Get(Id int32) *TzpWjx {
	for _, v := range td.Wjxs {
		if v.Wjid == Id {
			return v
		}
	}
	return nil
}

// InitGiftDB 提供给插件库内部使用
func InitGiftDB(dbAddr string) {

	if atomic.LoadInt32(&isGiftInit) != 0 {
		return
	}
	atomic.StoreInt32(&isGiftInit, 1)

	mysql.Init(dbAddr, 400)

	go func() {
		for {
			initGiftTab()
			time.Sleep(10 * time.Second)
		}

	}()

}

func initGiftTab() {
	gifts := make([]*TzpGift, 0)
	rows, err := mysql.Query("select * from tzp_gift")
	if err != nil {
		log.Error("select tzp_gift err:%v", err)
		return
	}

	for _, row := range rows {
		info := &TzpGift{}
		info.Id = row.GetInt32("id")

		chid := row.GetString("channel_id")
		if chid != "" {
			channelIdList := strings.Split(chid, ",")
			info.ChannelId = make(map[uint32]bool)
			for _, v := range channelIdList {
				cid, _ := strconv.Atoi(v)
				info.ChannelId[uint32(cid)] = true
			}
		}

		info.Name = row.GetString("name")
		info.Condition = SendMailCondition(row.GetInt32("condi"))
		info.StartTime = row.GetTime("start_time").Unix()
		info.EndTime = row.GetTime("end_time").Unix()
		info.Title = row.GetString("title")
		info.Content = row.GetString("content")
		info.Status = row.GetInt32("status")

		attach := row.GetString("attach")
		if len(attach) == 0 {
			continue
		}

		attachList := strings.Split(attach, "|")
		for _, v := range attachList {
			if len(v) == 0 {
				continue
			}
			kv := strings.Split(v, ":")
			if len(kv) != 2 {
				continue
			}
			ki, _ := strconv.Atoi(kv[0])
			vi, _ := strconv.Atoi(kv[1])
			info.Attach = append(info.Attach, &tzppub.ItemInfo{ItemId: int32(ki), ItemNum: int64(vi)})
		}

		info.Creator = row.GetString("creator")

		gifts = append(gifts, info)
	}

	TzpGiftdata = &TzpGiftData{Gifts: gifts}

}

// InitWjxDB 提供给插件库内部使用
func InitWjxDB(dbAddr string) {

	if atomic.LoadInt32(&isWjxInit) != 0 {
		return
	}
	atomic.StoreInt32(&isWjxInit, 1)

	mysql.Init(dbAddr, 400)

	go func() {
		for {
			initWjxTab()
			time.Sleep(10 * time.Second)
		}

	}()

}

func initWjxTab() {

	wjxs := make([]*TzpWjx, 0)
	rows, err := mysql.Query("select * from tzp_wjx")
	if err != nil {
		log.Error("select tzp_wjx err:%v", err)
		return
	}

	for _, row := range rows {
		info := &TzpWjx{}
		info.Wjid = row.GetInt32("wjid")

		chid := row.GetString("channel_id")
		if chid != "" {
			channelIdList := strings.Split(chid, ",")
			info.ChannelId = make(map[uint32]bool)
			for _, v := range channelIdList {
				cid, _ := strconv.Atoi(v)
				info.ChannelId[uint32(cid)] = true
			}
		}

		info.Reward = row.GetInt32("reward")
		info.Url = row.GetString("url")
		info.Content = row.GetString("content")
		info.Name = row.GetString("name")
		info.MinLev = row.GetInt32("min_lev")
		info.MaxLev = row.GetInt32("max_lev")
		info.MinPay = row.GetInt32("min_pay")
		info.MaxPay = row.GetInt32("max_pay")
		info.MinOpenDay = row.GetInt32("min_open_day")
		info.MaxOpenDay = row.GetInt32("max_open_day")
		info.StartTime = row.GetTime("start_time").Unix()
		info.EndTime = row.GetTime("end_time").Unix()
		info.Status = row.GetInt32("status")
		info.Creator = row.GetString("creator")
		info.Version = row.GetInt64("version")

		if info.MaxLev == 0 {
			info.MaxLev = MAXINT32
		}
		if info.MaxPay == 0 {
			info.MaxPay = MAXINT32
		}
		if info.MaxOpenDay == 0 {
			info.MaxOpenDay = MAXINT32
		}

		wjxs = append(wjxs, info)
	}

	sort.Slice(wjxs, func(i, j int) bool {
		return wjxs[i].Wjid < wjxs[j].Wjid
	})

	TzpWjxdata = &TzpWjxData{Wjxs: wjxs}

}
