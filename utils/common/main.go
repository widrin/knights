package main

import (
	"tzgit.kaixinxiyou.com/utils/common/bithelper"
	"tzgit.kaixinxiyou.com/utils/common/bytebuf"
	_ "tzgit.kaixinxiyou.com/utils/common/kafka_producer"
	_ "tzgit.kaixinxiyou.com/utils/common/mongoDb"
	"tzgit.kaixinxiyou.com/utils/common/queue"
	_ "tzgit.kaixinxiyou.com/utils/common/redis"
	"tzgit.kaixinxiyou.com/utils/common/util"
)

func main() {
	_ = bytebuf.NewByteBuf(0)
	bithelper.BitsNum(1)
	//rzip.Default.Unzip("aa", "aa")
	_ = util.GoID()
	_ = queue.Create()

}
