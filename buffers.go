package xlog

import (
	"fmt"
	"sync"
)

const (
	defLen   = 64
	fmtPlain = `%s=%s`
	fmtQuote = `%s="%s"`
)

type lbuf struct {
	buf []byte
}

var (
	lbp = sync.Pool{New: newBuf}
)

func newBuf() any {
	return &lbuf{buf: make([]byte, 0, defLen)}
}

func bufGet() *lbuf {
	return lbp.Get().(*lbuf)
}

func bufPut(b *lbuf) {
	b.reset()
	lbp.Put(b)
}

func (b *lbuf) reset() {
	b.buf = b.buf[:0]
}

func (b *lbuf) lastIsSpace() bool {
	n := len(b.buf)
	if n == 0 {
		return true
	}
	return b.buf[n-1] == ' '
}

func (b *lbuf) addSpaceOpt() {
	if !b.lastIsSpace() {
		b.buf = append(b.buf, ' ')
	}
}

func (b *lbuf) putStr(s string) {
	b.addSpaceOpt()
	b.buf = append(b.buf, s...)
}

func (b *lbuf) putAny(v any) {
	b.addSpaceOpt()
	b.buf = fmt.Append(b.buf, v)
}

func (b *lbuf) setStr(s string) {
	b.buf = append(b.buf[:0], s...)
}

func (b *lbuf) setAny(v any) {
	b.buf = fmt.Append(b.buf[:0], v)
}

func (b *lbuf) needQuote() bool { // todo utf8
	for _, v := range b.buf {
		if v <= ' ' {
			return true
		}
	}
	return false
}

func (b *lbuf) String() string {
	return string(b.buf)
}

func (b *lbuf) kvString(k string) string {
	var curFmt string
	if b.needQuote() {
		curFmt = fmtQuote
	} else {
		curFmt = fmtPlain
	}
	return fmt.Sprintf(curFmt, k, b.String())
}

func (b *lbuf) Write(data []byte) (int, error) {
	if len(data) > 0 {
		b.buf = append(b.buf, data...)
	}
	return len(data), nil
}
