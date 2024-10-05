package rrand

import (
	"math/rand"
	"sync"
	"time"
)

var Default *RRand

func init() {
	rand.Seed(time.Now().UnixNano())
	Default = New(time.Now().UnixNano())
}
func New(seed int64) *RRand {
	r := &RRand{}
	r.r = rand.New(rand.NewSource(seed))
	return r
}

type RRand struct {
	lk sync.Mutex
	r  *rand.Rand
}

func (r *RRand) ResetSeed(seed int64) {
	r.r.Seed(seed)
}
func (r *RRand) RandGroup(p ...uint32) int {
	if p == nil {
		panic("args not found")
	}
	r.lk.Lock()
	defer r.lk.Unlock()

	ret := make([]uint32, len(p))
	for i := 0; i < len(p); i++ {
		if i == 0 {
			ret[0] = p[0]
		} else {
			ret[i] = ret[i-1] + p[i]
		}
	}

	rl := ret[len(ret)-1]
	if rl == 0 {
		return 0
	}

	rn := uint32(r.r.Int63n(int64(rl)))
	for i := 0; i < len(ret); i++ {
		if rn < ret[i] {
			return i
		}
	}

	panic("bug")
}

func (r *RRand) RandInterval(b1, b2 int32) int32 {
	if b1 == b2 {
		return b1
	}
	r.lk.Lock()
	defer r.lk.Unlock()

	min, max := int64(b1), int64(b2)
	if min > max {
		min, max = max, min
	}
	return int32(r.r.Int63n(max-min+1) + min)
}
func (r *RRand) RandInterval64(b1, b2 int64) int64 {
	if b1 == b2 {
		return b1
	}
	r.lk.Lock()
	defer r.lk.Unlock()
	min, max := b1, b2
	if min > max {
		min, max = max, min
	}
	return r.r.Int63n(max-min+1) + min
}

func (r *RRand) RandIntervalN(b1, b2 int32, n uint32) []int32 {
	if b1 == b2 {
		return []int32{b1}
	}
	r.lk.Lock()
	defer r.lk.Unlock()
	min, max := int64(b1), int64(b2)
	if min > max {
		min, max = max, min
	}
	l := max - min + 1
	if int64(n) > l {
		n = uint32(l)
	}

	ret := make([]int32, n)
	m := make(map[int32]int32)
	for i := uint32(0); i < n; i++ {
		v := int32(r.r.Int63n(l) + min)

		if mv, ok := m[v]; ok {
			ret[i] = mv
		} else {
			ret[i] = v
		}

		lv := int32(l - 1 + min)
		if v != lv {
			if mv, ok := m[lv]; ok {
				m[v] = mv
			} else {
				m[v] = lv
			}
		}

		l--
	}

	return ret
}
