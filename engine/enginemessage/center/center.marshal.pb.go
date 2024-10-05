package center

import (
	math_bits "math/bits"
)

func (m *S2SActorHeart) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ActorIds) > 0 {
		for _, e := range m.ActorIds {
			n += 1 + sov(uint64(e))
		}
	}
	if m.ServerId != 0 {
		n += 1 + sov(uint64(m.ServerId))
	}
	if m.ActorType != 0 {
		n += 1 + sov(uint64(m.ActorType))
	}
	if m.GroupId != 0 {
		n += 1 + sov(uint64(m.GroupId))
	}
	return n
}

func (m *S2SActorHeart) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *S2SActorHeart) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *S2SActorHeart) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ActorIds) > 0 {
		for iNdEx := len(m.ActorIds) - 1; iNdEx >= 0; iNdEx-- {
			i = encodeVarint(dAtA, i, uint64(m.ActorIds[iNdEx]))
			i--
			dAtA[i] = 0x20
		}
	}
	if m.ServerId != 0 {
		i = encodeVarint(dAtA, i, uint64(m.ServerId))
		i--
		dAtA[i] = 0x18
	}
	if m.ActorType != 0 {
		i = encodeVarint(dAtA, i, uint64(m.ActorType))
		i--
		dAtA[i] = 0x10
	}
	if m.GroupId != 0 {
		i = encodeVarint(dAtA, i, uint64(m.GroupId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *S2SDistributionRuleServerReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.GroupId != 0 {
		n += 1 + sov(uint64(m.GroupId))
	}
	if m.ActorType != 0 {
		n += 1 + sov(uint64(m.ActorType))
	}
	if m.ActorId != 0 {
		n += 1 + sov(uint64(m.ActorId))
	}
	return n
}

func (m *S2SDistributionRuleServerReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *S2SDistributionRuleServerReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *S2SDistributionRuleServerReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.GroupId != 0 {
		i = encodeVarint(dAtA, i, uint64(m.GroupId))
		i--
		dAtA[i] = 0x18
	}
	if m.ActorType != 0 {
		i = encodeVarint(dAtA, i, uint64(m.ActorType))
		i--
		dAtA[i] = 0x10
	}
	if m.ActorId != 0 {
		i = encodeVarint(dAtA, i, uint64(m.ActorId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *S2SDistributionRuleServerRep) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ServerId != 0 {
		n += 1 + sov(uint64(m.ServerId))
	}
	if m.Tag != 0 {
		n += 1 + sov(uint64(m.Tag))
	}
	return n
}

func (m *S2SDistributionRuleServerRep) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *S2SDistributionRuleServerRep) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *S2SDistributionRuleServerRep) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ServerId != 0 {
		i = encodeVarint(dAtA, i, uint64(m.ServerId))
		i--
		dAtA[i] = 0x10
	}
	if m.Tag != 0 {
		i = encodeVarint(dAtA, i, uint64(m.Tag))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *S2SClearRuleCacheReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ServerId != 0 {
		n += 1 + sov(uint64(m.ServerId))
	}
	return n
}

func (m *S2SClearRuleCacheReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *S2SClearRuleCacheReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *S2SClearRuleCacheReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ServerId != 0 {
		i = encodeVarint(dAtA, i, uint64(m.ServerId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *S2SDistributionCenterServerReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ActorId != 0 {
		n += 1 + sov(uint64(m.ActorId))
	}
	return n
}

func (m *S2SDistributionCenterServerReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *S2SDistributionCenterServerReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *S2SDistributionCenterServerReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ActorId != 0 {
		i = encodeVarint(dAtA, i, uint64(m.ActorId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *S2SResetCenterServerReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ActorId != 0 {
		n += 1 + sov(uint64(m.ActorId))
	}
	return n
}

func (m *S2SResetCenterServerReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *S2SResetCenterServerReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *S2SResetCenterServerReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ActorId != 0 {
		i = encodeVarint(dAtA, i, uint64(m.ActorId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *S2SCloseCenterActor) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ActorId != 0 {
		n += 1 + sov(uint64(m.ActorId))
	}
	return n
}

func (m *S2SCloseCenterActor) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *S2SCloseCenterActor) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *S2SCloseCenterActor) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ActorId != 0 {
		i = encodeVarint(dAtA, i, uint64(m.ActorId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *S2SCloseCenterActorSuccess) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ActorId != 0 {
		n += 1 + sov(uint64(m.ActorId))
	}
	return n
}

func (m *S2SCloseCenterActorSuccess) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *S2SCloseCenterActorSuccess) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *S2SCloseCenterActorSuccess) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ActorId != 0 {
		i = encodeVarint(dAtA, i, uint64(m.ActorId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *S2SDistributionCenterServerRep) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ServerId != 0 {
		n += 1 + sov(uint64(m.ServerId))
	}
	if m.Tag != 0 {
		n += 1 + sov(uint64(m.Tag))
	}
	return n
}

func (m *S2SDistributionCenterServerRep) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *S2SDistributionCenterServerRep) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *S2SDistributionCenterServerRep) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ServerId != 0 {
		i = encodeVarint(dAtA, i, uint64(m.ServerId))
		i--
		dAtA[i] = 0x10
	}
	if m.Tag != 0 {
		i = encodeVarint(dAtA, i, uint64(m.Tag))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *S2SHeartCenterServerReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ActorIds) > 0 {
		for _, e := range m.ActorIds {
			n += 1 + sov(uint64(e))
		}
	}
	if m.ServerId != 0 {
		n += 1 + sov(uint64(m.ServerId))
	}
	if m.ActorType != 0 {
		n += 1 + sov(uint64(m.ActorType))
	}
	if m.GroupId != 0 {
		n += 1 + sov(uint64(m.GroupId))
	}
	return n
}

func (m *S2SHeartCenterServerReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *S2SHeartCenterServerReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *S2SHeartCenterServerReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ActorIds) > 0 {
		for iNdEx := len(m.ActorIds) - 1; iNdEx >= 0; iNdEx-- {
			i = encodeVarint(dAtA, i, uint64(m.ActorIds[iNdEx]))
			i--
			dAtA[i] = 0x20
		}
	}
	if m.ServerId != 0 {
		i = encodeVarint(dAtA, i, uint64(m.ServerId))
		i--
		dAtA[i] = 0x18
	}
	if m.ActorType != 0 {
		i = encodeVarint(dAtA, i, uint64(m.ActorType))
		i--
		dAtA[i] = 0x10
	}
	if m.GroupId != 0 {
		i = encodeVarint(dAtA, i, uint64(m.GroupId))
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
