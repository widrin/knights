package bithelper

import (
	"runtime/debug"
	"tzgit.kaixinxiyou.com/utils/common/log"
)

func HasBit(v int32, nIdx int32) bool {
	return v&(1<<uint32(nIdx)) != 0
}
func HasBit64(v int64, nIdx int32) bool {
	return v&(1<<uint32(nIdx)) != 0
}
func SetBit(v int32, nIdx int32) int32 {
	v = v | (1 << uint32(nIdx))
	return v
}
func SetBit64(v int64, nIdx int32) int64 {
	v = v | (1 << uint32(nIdx))
	return v
}
func ClrBit(v int32, nIdx int32) int32 {
	v = v & (^(1 << uint32(nIdx)))
	return v
}
func OrBits(v1 int32, vr int32) int32 {
	v1 |= vr
	return v1
}
func AddBits(v1 int32, vr int32) int32 {
	v1 &= vr
	return v1
}
func BitsNum(v1 int32) int32 {
	count := int32(0)
	for i := int32(0); i < 32; i++ {
		if HasBit(v1, i) {
			count++
		}
	}
	return count
}

//获取一个字节
func GetB(v uint32, nIdx uint32) uint32 {
	if nIdx == 0 {
		return v & 0x000000ff
	}
	if nIdx == 1 {
		return v & 0x0000ff00
	}
	if nIdx == 2 {
		return v & 0x00ff0000
	}
	if nIdx == 3 {
		return v & 0xff000000
	}
	log.Error("GetB 参数错误:%v  %s", nIdx, debug.Stack())
	return 0
}
func SetB(v uint32, nIdx uint32, vr uint32) uint32 {
	if nIdx == 0 {
		return v&0xffffff00 | vr
	}
	if nIdx == 1 {
		return v&0xffff00ff | (vr << 8)
	}
	if nIdx == 2 {
		return v&0x00ff0000 | (vr << 16)
	}
	if nIdx == 3 {
		return v&0xff000000 | (vr << 24)
	}
	log.Error("SetB 参数错误:%v  %s", nIdx, debug.Stack())
	return v
}
