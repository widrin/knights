package logCls

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
	"tzgit.kaixinxiyou.com/utils/common/buffer"
	"tzgit.kaixinxiyou.com/utils/common/log"
	"tzgit.kaixinxiyou.com/utils/common/log/syncWriter"
)

// 腾讯云cls 输出
type ClsOut struct {
	syncWriter.SyncWriter
	producerInstance *cls.AsyncProducerClient
	clsTopic         string
	cb               *callback
	addParam         map[string]string
}

func Create(clsAddr string, clsPass string, clsTopic string, level log.Level, param ...string) log.ISyncWriter {
	r := &ClsOut{}
	producerConfig := cls.GetDefaultAsyncProducerClientConfig()
	producerConfig.LingerMs = 200
	producerConfig.Endpoint = clsAddr
	t := strings.Split(clsPass, "#")
	producerConfig.AccessKeyID = t[0]
	producerConfig.AccessKeySecret = t[1]
	var err error
	r.producerInstance, err = cls.NewAsyncProducerClient(producerConfig)
	if err != nil {
		log.Error(" %v", err)
		return nil
	}
	r.clsTopic = clsTopic
	r.cb = &callback{}
	r.addParam = make(map[string]string)
	if len(param) > 0 {
		for i := 0; i < len(param)-1; i += 2 {
			r.addParam[param[i]] = param[i+1]
		}
	}

	r.producerInstance.Start()
	//注册
	r.InitBufferWriter(r, 0, 0, level)

	return r
}
func (c *ClsOut) SyncWrite(b *buffer.Buffer) {
	stu := map[string]string{}
	err := jsoniter.Unmarshal(b.Bytes(), &stu)
	for k, v := range c.addParam {
		stu[k] = v
	}
	aa := cls.NewCLSLog(time.Now().UnixMilli(), stu)
	err = c.producerInstance.SendLog(c.clsTopic, aa, c.cb)
	if err != nil {
		fmt.Println(err)
	}
	//log.InitCls(zap.String("serverid", convert.ToString(conf.ServerId)), zap.String("appid", convert.ToString(conf.AppId)))
	//level := stu["L"]
	//for _, v := range tbCfg.GetCfg().GetErrorWarnLevel() {
	//	if level == v {
	//		//发送预警
	//		bb := cls.NewCLSLog(time.Now().Unix(), stu)
	//		producerInstance.SendLog(tbCls.GetTopicId("log_error"), bb, clbCallBack)
	//	}
	//}
}
func (c *ClsOut) GetEncoder() zapcore.Encoder {
	kafkaEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	return kafkaEncoder
}

// 关闭事件
func (c *ClsOut) OnClose() {
	c.producerInstance.Close(1000 * 60)
}

type callback struct {
}

func (callback *callback) Success(_ *cls.Result) {
}

func (callback *callback) Fail(result *cls.Result) {
	fmt.Println(result.IsSuccessful())
	fmt.Println(result.GetErrorCode())
	fmt.Println(result.GetErrorMessage())
	for _, v := range result.GetReservedAttempts() {
		fmt.Println("err", v)
	}
	fmt.Println(result.GetReservedAttempts())
	fmt.Println(result.GetRequestId())
	fmt.Println(result.GetTimeStampMs())
}
