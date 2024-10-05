package wjx

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"tzgit.kaixinxiyou.com/engine/tzplugin/internal/inpub"
	"tzgit.kaixinxiyou.com/engine/tzplugin/tzppub"
	"tzgit.kaixinxiyou.com/utils/common/log"
)

// IPlayerWjx 游戏中问卷星模块需要实现的接口
type IPlayerWjx interface {
	OnStart(info *Info) //开启图标回调
	OnClose()           //关闭图标回调
}

// 问卷星展示信息
type Info struct {
	Wjid      int32              //数据库中的问卷ID
	Reward    []*tzppub.ItemInfo //奖励内容
	Url       string             //需要加载的URL
	Content   string             //展示内容
	Name      string             //问卷名
	StartTime int64              //问卷开始时间
	EndTime   int64              //问卷结束时间
	Version   int64              //问卷版本号
}

// WjxData 问卷星数据结构
type WjxData struct {
	Playerwjx    IPlayerWjx     //玩家问卷星接口
	Player       inpub.IPlayer  //玩家接口
	CacheInfo    *Info          //优化缓存
	AlreadyPrize map[int32]bool //已经领取的问卷调查奖励Id,该数据需要落地到游戏数据库
}

// Init 初始化问卷星功能
// 在游戏服中初始化
// dbAddr sql地址
func Init(dbAddr string) {
	inpub.InitGiftDB(dbAddr)
	inpub.InitWjxDB(dbAddr)
}

// CreatePlayerWjxData 创建玩家问卷星模块数据
// player 玩家接口
// playerwjx 玩家问卷星接口
func CreatePlayerWjxData(player inpub.IPlayer, playerwjx IPlayerWjx) *WjxData {
	w := &WjxData{}
	w.Playerwjx = playerwjx
	w.Player = player
	w.AlreadyPrize = make(map[int32]bool)
	return w
}

// ToCheck 游戏侧定期（一般是1分钟)调用该该方法来检测问卷星状态
// isLogin 需要实现在登录时候isLogin设置为true。如果不设置，玩家在重登的时候不会显示问卷图标
func (w *WjxData) ToCheck(isLogin bool) {

	if isLogin {
		w.CacheInfo = nil
	}

	cid := w.Player.TzpGetChannelId()
	lv := w.Player.TzpGetLevel()
	openday := w.Player.TzpGetAreaOpenDay()
	chargefen := w.Player.TzpGetChargeValue()
	platId := w.Player.TzpGetPlatformId()
	accountId := w.Player.TzpGetAccountId()
	userid := w.Player.TzpGetUserID()

	nt := time.Now().Unix()
	for _, v := range inpub.TzpWjxdata.Wjxs {

		if v.Status != 0 && v.StartTime <= nt && nt <= v.EndTime &&
			v.MinLev <= lv && lv <= v.MaxLev &&
			v.MinOpenDay <= openday && openday <= v.MaxOpenDay &&
			v.MinPay <= chargefen && chargefen <= v.MaxPay {

			if len(v.ChannelId) != 0 {
				_, ok := v.ChannelId[cid]
				if !ok {
					continue
				}
			}

			reward := inpub.TzpGiftdata.Get(v.Reward)
			if reward == nil {
				continue
			}

			_, already := w.AlreadyPrize[v.Wjid]
			if already {
				continue
			}

			if w.CacheInfo != nil && w.CacheInfo.Version == v.Version {
				log.Debug("=====版本相同=跳过==玩家：%v 问卷ID:%v", userid, w.CacheInfo.Wjid)
				return
			}

			info := &Info{}
			info.Wjid = v.Wjid

			suserid := strconv.Itoa(int(userid))
			wjid := strconv.Itoa(int(v.Wjid))
			if strings.Contains(v.Url, "?") {
				info.Url = fmt.Sprintf("%v&sojumpparm=%v|%v|%v&q60001=%v&q60002=%v&q60003=%v&q60004=%v",
					v.Url, wjid, 0, suserid, cid, platId, userid, accountId)
			} else {
				info.Url = fmt.Sprintf("%v?sojumpparm=%v|%v|%v&q60001=%v&q60002=%v&q60003=%v&q60004=%v",
					v.Url, wjid, 0, suserid, cid, platId, userid, accountId)
			}

			info.Content = v.Content
			info.Name = v.Name
			info.StartTime = v.StartTime
			info.EndTime = v.EndTime
			info.Version = v.Version

			info.Reward = reward.Attach
			w.CacheInfo = info

			log.Debug("发现问卷，打开入口。玩家:%v,问卷:%v", userid, wjid)
			w.Playerwjx.OnStart(info)
			return

		}
	}

	if w.CacheInfo != nil {
		w.CacheInfo = nil
		log.Debug("没有问卷，关闭入口。玩家:%v", userid)
		w.Playerwjx.OnClose()
	}
}

// SetData 设置问卷星领取状态
// 游戏侧在加载问卷星模块数据后设置插件库数据
func (w *WjxData) SetData(data []int32) {

	for _, v := range data {
		w.AlreadyPrize[v] = true
	}

}

// GetData 获取问卷星领取状态
// 游戏侧在落地玩家数据时候获取该值保存到游戏侧的数据库中
func (w *WjxData) GetData() (AlreadyPrizeList []int32) {
	for k := range w.AlreadyPrize {
		AlreadyPrizeList = append(AlreadyPrizeList, k)
	}
	return
}
