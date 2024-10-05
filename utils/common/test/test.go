package main

import (
	"tzgit.kaixinxiyou.com/utils/common/log"
	"tzgit.kaixinxiyou.com/utils/common/log/logFile"
	"tzgit.kaixinxiyou.com/utils/common/log/logStd"
)

func main() {

	log.RegisterLog(logStd.Create(0, log.DebugLevel))
	log.RegisterLogWithField(logFile.Create("test", 0, log.DebugLevel), "server_id", 1)
	//log.RegisterLog(logFile.Create("test_err", 0, log.ErrorLevel))
	for i := 0; i < 100; i++ {
		log.Debug("aaaaaaaaaaa:%v", i)
		//log.Error("aaaaaaaaaaa:%v", i)
	}
	log.Close()
}
