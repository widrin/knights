package jsonEncoder

import (
	"encoding/base64"
	"math"
	"sync"
	"time"
	"tzgit.kaixinxiyou.com/utils/common/buffer"
	"unicode/utf8"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

// For JSON-escaping; see jsonEncoder.safeAddString below.
const _hex = "0123456789abcdef"

var _jsonPool = sync.Pool{New: func() interface{} {
	v := &jsonEncoder{}
	v.buf = buffer.Get()
	return v
}}

func Get() ObjectEncoder {
	ret := _jsonPool.Get().(*jsonEncoder)
	ret.buf.AppendByte('{')
	return ret
}

func Put(enc ObjectEncoder) {
	enc.Reset()
	_jsonPool.Put(enc)
}

func (enc *jsonEncoder) Reset() {
	enc.buf.Reset()
	enc.isMarshal = false
	//enc.buf = nil
	enc.spaced = false
}

type jsonEncoder struct {
	buf    *buffer.Buffer
	spaced bool // include spaces after colons and commas
	// for encoding generic values by reflection
	openNamespaces int
	isMarshal      bool
}

func (enc *jsonEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.AppendTimeLayout(val, timeFormat)
}

func (enc *jsonEncoder) OpenNamespace(key string) {
	if enc.isMarshal {
		return
	}
	enc.addKey(key)
	enc.buf.AppendByte('{')
	enc.openNamespaces++
}

//func (enc *jsonEncoder) AddArray(key string, arr ArrayMarshaler) error {
//	enc.addKey(key)
//	return enc.AppendArray(arr)
//}
//
//func (enc *jsonEncoder) AddObject(key string, obj ObjectMarshaler) error {
//	enc.addKey(key)
//	return enc.AppendObject(obj)
//}

func (enc *jsonEncoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (enc *jsonEncoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.AppendByteString(val)
}

func (enc *jsonEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
}

func (enc *jsonEncoder) AddComplex128(key string, val complex128) {
	enc.addKey(key)
	enc.AppendComplex128(val)
}

func (enc *jsonEncoder) AddComplex64(key string, val complex64) {
	enc.addKey(key)
	enc.AppendComplex64(val)
}

func (enc *jsonEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(val)
}

func (enc *jsonEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.AppendFloat64(val)
}

func (enc *jsonEncoder) AddFloat32(key string, val float32) {
	enc.addKey(key)
	enc.AppendFloat32(val)
}

func (enc *jsonEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.AppendInt64(val)
}

func (enc *jsonEncoder) addKey(key string) {
	if enc.isMarshal {
		return
	}
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddString(key)
	enc.buf.AppendByte('"')
	enc.buf.AppendByte(':')
	if enc.spaced {
		enc.buf.AppendByte(' ')
	}
}
func (enc *jsonEncoder) addElementSeparator() {
	if enc.isMarshal {
		return
	}
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}
	switch enc.buf.Bytes()[last] {
	case '{', '[', ':', ',', ' ':
		return
	default:
		enc.buf.AppendByte(',')
		if enc.spaced {
			enc.buf.AppendByte(' ')
		}
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (enc *jsonEncoder) safeAddString(s string) {
	if enc.isMarshal {
		return
	}
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (enc *jsonEncoder) tryAddRuneSelf(b byte) bool {
	if enc.isMarshal {
		return false
	}
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}
func (enc *jsonEncoder) tryAddRuneError(r rune, size int) bool {
	if enc.isMarshal {
		return false
	}
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

func (enc *jsonEncoder) AppendDuration(val time.Duration) {
	if enc.isMarshal {
		return
	}
	cur := enc.buf.Len()
	//if e := enc.EncodeDuration; e != nil {
	//	e(val, enc)
	//}
	if cur == enc.buf.Len() {
		// User-supplied EncodeDuration is a no-op. Fall back to nanoseconds to keep
		// JSON valid.
		enc.AppendInt64(int64(val))
	}
}

func (enc *jsonEncoder) AppendInt64(val int64) {
	if enc.isMarshal {
		return
	}
	enc.addElementSeparator()
	enc.buf.AppendInt(val)
}

//func (enc *jsonEncoder) AppendReflected(val interface{}) error {
//	valueBytes, err := enc.encodeReflected(val)
//	if err != nil {
//		return err
//	}
//	enc.addElementSeparator()
//	_, err = enc.buf.Write(valueBytes)
//	return err
//}

func (enc *jsonEncoder) AppendString(val string) {
	if enc.isMarshal {
		return
	}
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddString(val)
	enc.buf.AppendByte('"')
}

func (enc *jsonEncoder) AppendTimeLayout(time time.Time, layout string) {
	if enc.isMarshal {
		return
	}
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.buf.AppendTime(time, layout)
	enc.buf.AppendByte('"')
}

//func (enc *jsonEncoder) AppendTime(val time.Time) {
//	cur := enc.buf.Len()
//	if e := enc.EncodeTime; e != nil {
//		e(val, enc)
//	}
//	if cur == enc.buf.Len() {
//		// User-supplied EncodeTime is a no-op. Fall back to nanos since epoch to keep
//		// output JSON valid.
//		enc.AppendInt64(val.UnixNano())
//	}
//}

func (enc *jsonEncoder) AppendUint64(val uint64) {
	if enc.isMarshal {
		return
	}
	enc.addElementSeparator()
	enc.buf.AppendUint(val)
}

func (enc *jsonEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.AppendUint64(val)
}

// appendComplex appends the encoded form of the provided complex128 value.
// precision specifies the encoding precision for the real and imaginary
// components of the complex number.
func (enc *jsonEncoder) appendComplex(val complex128, precision int) {
	if enc.isMarshal {
		return
	}
	enc.addElementSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, precision)
	// If imaginary part is less than 0, minus (-) sign is added by default
	// by AppendFloat.
	if i >= 0 {
		enc.buf.AppendByte('+')
	}
	enc.buf.AppendFloat(i, precision)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}
func (enc *jsonEncoder) appendFloat(val float64, bitSize int) {
	if enc.isMarshal {
		return
	}
	enc.addElementSeparator()
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

func (enc *jsonEncoder) AddInt(k string, v int)         { enc.AddInt64(k, int64(v)) }
func (enc *jsonEncoder) AddInt32(k string, v int32)     { enc.AddInt64(k, int64(v)) }
func (enc *jsonEncoder) AddInt16(k string, v int16)     { enc.AddInt64(k, int64(v)) }
func (enc *jsonEncoder) AddInt8(k string, v int8)       { enc.AddInt64(k, int64(v)) }
func (enc *jsonEncoder) AddUint(k string, v uint)       { enc.AddUint64(k, uint64(v)) }
func (enc *jsonEncoder) AddUint32(k string, v uint32)   { enc.AddUint64(k, uint64(v)) }
func (enc *jsonEncoder) AddUint16(k string, v uint16)   { enc.AddUint64(k, uint64(v)) }
func (enc *jsonEncoder) AddUint8(k string, v uint8)     { enc.AddUint64(k, uint64(v)) }
func (enc *jsonEncoder) AddUintptr(k string, v uintptr) { enc.AddUint64(k, uint64(v)) }
func (enc *jsonEncoder) AppendComplex64(v complex64)    { enc.appendComplex(complex128(v), 32) }
func (enc *jsonEncoder) AppendComplex128(v complex128)  { enc.appendComplex(complex128(v), 64) }
func (enc *jsonEncoder) AppendFloat64(v float64)        { enc.appendFloat(v, 64) }
func (enc *jsonEncoder) AppendFloat32(v float32)        { enc.appendFloat(float64(v), 32) }
func (enc *jsonEncoder) AppendInt(v int)                { enc.AppendInt64(int64(v)) }
func (enc *jsonEncoder) AppendInt32(v int32)            { enc.AppendInt64(int64(v)) }
func (enc *jsonEncoder) AppendInt16(v int16)            { enc.AppendInt64(int64(v)) }
func (enc *jsonEncoder) AppendInt8(v int8)              { enc.AppendInt64(int64(v)) }
func (enc *jsonEncoder) AppendUint(v uint)              { enc.AppendUint64(uint64(v)) }
func (enc *jsonEncoder) AppendUint32(v uint32)          { enc.AppendUint64(uint64(v)) }
func (enc *jsonEncoder) AppendUint16(v uint16)          { enc.AppendUint64(uint64(v)) }
func (enc *jsonEncoder) AppendUint8(v uint8)            { enc.AppendUint64(uint64(v)) }
func (enc *jsonEncoder) AppendUintptr(v uintptr)        { enc.AppendUint64(uint64(v)) }
func (enc *jsonEncoder) AppendBool(val bool) {
	if enc.isMarshal {
		return
	}
	enc.addElementSeparator()
	enc.buf.AppendBool(val)
}
func (enc *jsonEncoder) AppendByteString(val []byte) {
	if enc.isMarshal {
		return
	}
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddByteString(val)
	enc.buf.AppendByte('"')
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *jsonEncoder) safeAddByteString(s []byte) {
	if enc.isMarshal {
		return
	}
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.Write(s[i : i+size])
		i += size
	}
}
func (enc *jsonEncoder) AddString(key, val string) {
	enc.addKey(key)
	enc.AppendString(val)
}
func (enc *jsonEncoder) Marshal() []byte {
	if enc.isMarshal == false {
		enc.buf.AppendString("}")
		enc.isMarshal = true
	}
	return enc.buf.Bytes()
}
func (enc *jsonEncoder) Clone() []byte {
	d := enc.Marshal()
	ret := make([]byte, len(d))
	copy(ret, d)
	return d
}
