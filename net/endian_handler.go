package net

import (
	"encoding/binary"
	"errors"

	"github.com/widrin/knights/logger"
)

type EndianHandler struct {
	order binary.ByteOrder
}

func NewEndianHandler(endianType string) *EndianHandler {
	switch endianType {
	case "big":
		return &EndianHandler{order: binary.BigEndian}
	case "little":
		return &EndianHandler{order: binary.LittleEndian}
	default:
		logger.Warn("Unsupported endian type, defaulting to little endian")
		return &EndianHandler{order: binary.LittleEndian}
	}
}

func (e *EndianHandler) DecodeUint32(b []byte) (uint32, error) {
	if len(b) < 4 {
		return 0, errors.New("insufficient bytes for uint32 decoding")
	}
	return e.order.Uint32(b), nil
}

func (e *EndianHandler) EncodeUint32(v uint32) []byte {
	buf := make([]byte, 4)
	e.order.PutUint32(buf, v)
	return buf
}

// 实现其他基本类型的编解码方法
