package buffer

import "errors"

const maxInt = int(^uint(0) >> 1)

var ErrTooLarge = errors.New("bytes.Buffer: too large")

// tryGrowByReslice is a inlineable version of grow for the fast-case where the
// internal buffer only needs to be resliced.
// It returns the index where bytes should be written and whether it succeeded.
func tryGrowByReslice(buf []byte, n int) ([]byte, bool) {
	if l := len(buf); n <= cap(buf)-l {
		buf = buf[:l+n]
		return buf, true
	}
	return buf, false
}

// grow grows the buffer to guarantee space for n more bytes.
// It returns the index where bytes should be written.
// If the buffer can't grow it will panic with ErrTooLarge.
func grow(buf []byte, n int) []byte {

	// Try to grow by means of a reslice.
	if i, ok := tryGrowByReslice(buf, n); ok {
		return i
	}

	if buf == nil && n <= 64 {
		buf = make([]byte, n, 64)
		return buf
	}
	c := cap(buf)
	m := len(buf)
	if n <= c/2-m {
		// We can slide things down instead of allocating a new
		// slice. We only need m+n <= c to slide, but
		// we instead let capacity get twice as large so we
		// don't spend all our time copying.
	} else if c > maxInt-c-n {
		panic(ErrTooLarge)
	} else {
		// Not enough space anywhere, we need to allocate.
		_buf := makeSlice(2*c + n)
		copy(_buf, buf[:])
		buf = _buf
	}
	// Restore b.off and len(b.buf).
	buf = buf[:m+n]
	return buf
}

// Grow grows the buffer's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to the
// buffer without another allocation.
// If n is negative, Grow will panic.
// If the buffer can't grow it will panic with ErrTooLarge.
func Grow(buf []byte, n int) []byte {
	if n < 0 {
		panic("bytes.Buffer.Grow: negative count")
	}
	buf = grow(buf, n)
	return buf
}

// makeSlice allocates a slice of size n. If the allocation fails, it panics
// with ErrTooLarge.
func makeSlice(n int) []byte {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
	return make([]byte, n)
}
