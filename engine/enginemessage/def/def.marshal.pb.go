package def

import (
	math_bits "math/bits"
)

func (m *S2CRep) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Tag != 0 {
		n += 1 + sov(uint64(m.Tag))
	}
	return n
}

func (m *S2CRep) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *S2CRep) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *S2CRep) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Tag != 0 {
		i = encodeVarint(dAtA, i, uint64(m.Tag))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarint(dAtA []byte, offset int, v uint64) int {
	offset -= sov(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func sov(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func ByteSov(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func soz(x uint64) (n int) {
	return sov(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
