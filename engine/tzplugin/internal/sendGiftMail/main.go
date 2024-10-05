package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"tzgit.kaixinxiyou.com/engine/tzplugin"
)

var (
	URL        = flag.String("url", "", "发放礼包地址")
	UseridList = flag.String("user", "", "角色ID|角色ID")
	GiftID     = flag.Int64("gif", 0, "礼包ID")
	Key        = flag.String("key", "", "签名Key")
)

func main() {

	flag.Parse()

	p := &tzplugin.SendGiftData{}
	p.GiftID = int32(*GiftID)
	p.Timestamp = time.Now().Unix()
	p.UseridList = *UseridList

	sha1h := sha1.New()
	_, _ = sha1h.Write([]byte(fmt.Sprintf("%v%v%v%v", p.GiftID, p.GiftID, p.Timestamp, *Key)))
	p.Sign = hex.EncodeToString(sha1h.Sum(nil))

	jsonData, _ := json.Marshal(p)

	request, err := http.NewRequest("POST", *URL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("NewRequest err", err.Error())
		return
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Do err", err.Error())
		return
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	// fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

}
