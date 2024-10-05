package tzplugin

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"tzgit.kaixinxiyou.com/engine/tzplugin/internal/inpub"
	"tzgit.kaixinxiyou.com/engine/tzplugin/internal/inwjx"
	"tzgit.kaixinxiyou.com/utils/common/log"
	"tzgit.kaixinxiyou.com/utils/common/mysql"
)

const (
	EmptyJSONStr = "{}"
)

var (
	guc  inpub.UserCallback
	gkey string
)

//Init 初始化 需要传递实现了调用下方PlayerAction的对象
//一般在公网暴露的api服务器上调用
//giftkey 礼包发放签名key
//wjxAppkey 问卷星验证签名key(推送密钥)，如果不使用可传递空字符串
func Init(uc inpub.UserCallback, dbAddr string, giftkey string, wjxAppkey string) {
	guc = uc
	gkey = giftkey
	inpub.WjxAppkey = wjxAppkey
	inpub.InitGiftDB(dbAddr)
}

// TzpPlayerCall 当收到Init的TzpCallToUserActor后，需要在玩家actor中调用此函数
// 需要在玩家actor上调用
// 需要在调用此函数后置脏玩家数据。
func TzpPlayerCall(p inpub.IPlayer, param1, param2, param3 string) string {

	switch param1 {
	case "wjx":
		return inwjx.PrivCall(p, param2, param3)
	}

	log.Debug("TzpPlayerCall意外的参数1:%v", param1)
	return EmptyJSONStr
}

// SendGiftData 发送礼包的json格式
type SendGiftData struct {
	UseridList string //角色ID列表，使用|来分割
	GiftID     int32  //礼包ID
	Timestamp  int64  //秒级时间戳
	Sign       string //签名=Sha1(UseridList+GiftID+Timestamp+Key)
}

// SendGiftHandler 邮件礼包功能,需要先Init
func SendGiftHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Debug("发送礼包回调读取错误：%v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	time.Now().Unix()

	if len(body) == 0 {
		log.Debug("发送礼包回数据为空")
		http.Error(w, "body is empty", http.StatusBadRequest)
		return
	}

	log.Debug("发送礼包回调数据:%v", string(body))

	var p SendGiftData
	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Debug("发送礼包回调反序列化出错：%v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nt := time.Now().Unix()
	if p.Timestamp+30 < nt {
		log.Debug("无效的时间戳")
		http.Error(w, "无效的时间戳", http.StatusBadRequest)
		return
	}

	lowSign := strings.ToLower(p.Sign)

	sha1h := sha1.New()
	_, _ = sha1h.Write([]byte(fmt.Sprintf("%v%v%v%v", p.GiftID, p.GiftID, p.Timestamp, gkey)))
	result := hex.EncodeToString(sha1h.Sum(nil))

	if lowSign != result {
		log.Debug("签名不正确")
		http.Error(w, "签名不正确", http.StatusBadRequest)
		return
	}

	sUseridList := strings.Split(p.UseridList, "|")
	var iUseridList []uint32
	for _, suid := range sUseridList {
		iuid, _ := strconv.Atoi(suid)
		if iuid != 0 {
			iUseridList = append(iUseridList, uint32(iuid))
		}
	}

	sendGift(p.GiftID, iUseridList, false)
	w.Write([]byte("OK"))

}

func logGift(reason string, rewardid int, chanid int, gifName string, platId int64, accountId uint32, userid int, gifCreator string, status int, isSystem bool) {

	log.Debug(reason)
	if isSystem {
		gifCreator = "系统"
	}
	_, err := mysql.Exec("insert into tzp_gift_log (gift_id,channel_id,name,plat_id,account_id,role_id,send_time,status,creator,reason)values(?,?,?,?,?,?,?,?,?,?)",
		rewardid, chanid, gifName, strconv.Itoa(int(platId)), strconv.Itoa(int(accountId)), userid, time.Now(), status, gifCreator, reason)
	if err != nil {
		log.Warn("insert tzp_gift_log err %v", err)
	}
}

func sendGift(giftid int32, useridList []uint32, isSystem bool) {

	log.Debug("天纵插件库sendGift %v %v", giftid, useridList)

	gif := inpub.TzpGiftdata.Get(giftid)
	if gif == nil {
		log.Debug("找不到礼包ID:%v", giftid)
		return
	}

	if gif.Status == 0 {
		log.Debug("礼包ID:%v 未启用", giftid)
		return
	}

	nt := time.Now().Unix()

	if nt < gif.StartTime || gif.EndTime < nt {
		log.Debug("礼包ID:%v 不在有效期", giftid)
		return
	}

	var sendUseridList []uint32

	cycleDS := gif.Condition.GetString()

	for _, userid := range useridList {

		cid := guc.TzpGetChannelId(userid)
		platId := guc.TzpGetPlatformId(userid)
		accountId := guc.TzpGetAccountId(userid)

		if len(gif.ChannelId) != 0 {
			_, ok := gif.ChannelId[cid]
			if !ok {
				reason := fmt.Sprintf("玩家：%v，礼包ID:%v 玩家与礼包的渠道不匹配", userid, giftid)
				logGift(reason, int(gif.Id), int(cid), gif.Name, platId, accountId, int(userid), gif.Creator, 0, isSystem)
				continue
			}
		}

		if cycleDS != "" {

			dkey := fmt.Sprintf("%v_%v_%v", userid, gif.Id, cycleDS)
			_, err := mysql.Exec("insert into tzp_gift_record (dkey,ts)values(?,?)", dkey, time.Now())
			if err != nil {

				reason := fmt.Sprintf("insert tzp_gift_record err %v", err.Error())
				if strings.Contains(err.Error(), "Duplicate") {
					reason = "礼包在时间段内已经发送过:" + dkey
				}

				logGift(reason, int(gif.Id), int(cid), gif.Name, platId, accountId, int(userid), gif.Creator, 0, isSystem)
				continue
			}
		}

		logGift("", int(gif.Id), int(cid), gif.Name, platId, accountId, int(userid), gif.Creator, 1, isSystem)

		sendUseridList = append(sendUseridList, userid)

	}

	log.Debug("天纵插件库 TzpSendGift %v %v", giftid, sendUseridList)
	guc.TzpSendGift(sendUseridList, gif.Title, gif.Content, gif.Attach)
}

// WjxHandler 问卷星回调处理函数,需要先Init
func WjxHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Debug("问卷星推送回调读取错误：%v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		log.Debug("问卷星推送回数据为空")
		http.Error(w, "body is empty", http.StatusBadRequest)
		return
	}

	log.Debug("问卷星推送回调数据:%v", string(body))

	var p inwjx.WjxPushJson
	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Debug("问卷星推送回调反序列化出错：%v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sha1h := sha1.New()
	_, _ = sha1h.Write([]byte(p.Activity + p.Index + inpub.WjxAppkey))
	result := hex.EncodeToString(sha1h.Sum(nil))

	if result != p.Sign {
		log.Debug("问卷星签名错误，Activity:%v Index:%v Sign:%v", p.Activity, p.Index, p.Sign)
		http.Error(w, "问卷星签名错误", http.StatusForbidden)
		return
	}

	WjidRouteridUserid := strings.SplitN(p.SoJumpParm, "|", 3)

	if len(WjidRouteridUserid) != 3 {
		w.Write([]byte("测试参数错误"))
		http.Error(w, "测试参数错误", http.StatusMethodNotAllowed)
		return
	}

	routerid, _ := strconv.Atoi(WjidRouteridUserid[1])
	userid, _ := strconv.Atoi(WjidRouteridUserid[2])

	//Activity: p.Activity, Index: p.Index, Sign: p.Sign
	data := inwjx.CallToUserData{WjID: WjidRouteridUserid[0]}
	bdata, _ := json.Marshal(&data)

	resp := guc.TzpCallToUserActor(uint32(routerid), uint32(userid), "wjx", "wjidCallback", string(bdata))

	log.Debug("TzpCallToUserActor返回数据:%v", resp)

	if len(resp) != 0 {
		var p inwjx.CallPlayerResp
		_ = json.Unmarshal([]byte(resp), &p)
		sendGift(p.GiftID, []uint32{uint32(userid)}, true)
	}

	w.Write([]byte("感谢您的调查问卷。"))

}

// WjxTestHandler ！！问卷星【测试】回调处理函数,需要先Init！！
func WjxTestHandler(w http.ResponseWriter, r *http.Request) {

	values := r.URL.Query()
	sojumpparm := values.Get("sojumpparm")
	if sojumpparm == "" {
		w.Write([]byte("测试参数为空"))
		return
	}

	WjidRouteridUserid := strings.SplitN(sojumpparm, "|", 3)
	if len(WjidRouteridUserid) != 3 {
		w.Write([]byte("测试参数错误"))
		return
	}

	routerid, _ := strconv.Atoi(WjidRouteridUserid[1])
	userid, _ := strconv.Atoi(WjidRouteridUserid[2])
	data := inwjx.CallToUserData{WjID: WjidRouteridUserid[0]}
	bdata, _ := json.Marshal(&data)

	resp := guc.TzpCallToUserActor(uint32(routerid), uint32(userid), "wjx", "wjidCallback", string(bdata))

	log.Debug("TzpCallToUserActor返回数据:%v", resp)

	if len(resp) != 0 {
		var p inwjx.CallPlayerResp
		_ = json.Unmarshal([]byte(resp), &p)
		sendGift(p.GiftID, []uint32{uint32(userid)}, true)
	}

	w.Write([]byte("测试调查问卷加载完成,数据:" + sojumpparm))

}
