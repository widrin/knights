package inwjx

import (
	"encoding/json"
	"strconv"

	"tzgit.kaixinxiyou.com/engine/tzplugin/internal/inpub"
	"tzgit.kaixinxiyou.com/engine/tzplugin/wjx"
	"tzgit.kaixinxiyou.com/utils/common/log"
)

const (
	EmptyJSONStr = "{}"
)

// WjxPushJson 问卷星回调数据结构
type WjxPushJson struct {
	SoJumpParm string `json:"sojumpparm"` //传递的参数
	Activity   string `json:"activity"`   //问卷ID（问卷星侧ID)
	Index      string `json:"index"`      //作答序号
	Sign       string `json:"sign"`       //sign=sha1(activity+index+推送密钥)
}

type CallToUserData struct {
	WjID string `json:"wjid"` //问卷ID（数据库ID)
	// Activity string `json:"activity"` //问卷ID（问卷星侧ID)
	// Index    string `json:"index"`    //作答序号
	// Sign     string `json:"sign"`     //sign=sha1(activity+index+推送密钥)
}

type CallPlayerResp struct {
	GiftID int32
}

// PrivCall 插件库内部调用,实现玩家功能执行
func PrivCall(p inpub.IPlayer, _, wjdatas string) string {

	var wjdata CallToUserData
	_ = json.Unmarshal([]byte(wjdatas), &wjdata)

	//Activity:%v Index:%v Sign:%v wjdata.Activity, wjdata.Index, wjdata.Sign
	log.Debug("问卷星系统回调，PrivCall，问卷ID:%v", wjdata.WjID)

	userid := p.TzpGetUserID()

	iwjid, _ := strconv.Atoi(wjdata.WjID)

	wi := p.TzpGetWjxData()
	if wi == nil {
		log.Debug("问卷星系统回调，找不到%v玩家的问卷星数据", userid)
		return EmptyJSONStr
	}

	w, ok := wi.(*wjx.WjxData)
	if !ok {
		log.Debug("问卷星系统回调，游戏侧返回错误的数据类型,%v", userid)
		return EmptyJSONStr
	}

	wj := inpub.TzpWjxdata.Get(int32(iwjid))
	if wj == nil {
		log.Debug("问卷星系统回调，玩家：%v，找不到问卷ID:%v", userid, iwjid)
		return EmptyJSONStr
	}

	isAlready, ok := w.AlreadyPrize[int32(iwjid)]
	if ok && isAlready {
		//防刷
		log.Debug("问卷星系统回调，问卷邮件已经发送过。 玩家:%v，问卷ID:%v的礼包ID:%v ", userid, iwjid, wj.Reward)
		return EmptyJSONStr
	}
	w.AlreadyPrize[int32(iwjid)] = true

	w.CacheInfo = nil
	log.Debug("完成问卷，关闭入口。玩家:%v，问卷ID:%v", userid, iwjid)
	w.Playerwjx.OnClose()
	w.ToCheck(false)

	b, _ := json.Marshal(&CallPlayerResp{GiftID: wj.Reward})

	return string(b)
}
