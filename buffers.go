package xlog

import (
	"fmt"
	"os"
	"sync"
)

const (
	defLen = 64
)

type lbuf struct {
	buf []byte
}

var (
	lbp = sync.Pool{
		New: func() interface{} {
			return &lbuf{
				buf: make([]byte, 0, defLen),
			}
		},
	}
)

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

//func (b *lbuf) Write(p []byte) (n int, err error) {
//	b.buf = append(b.buf, p...)
//	return len(p), nil
//}

func (b *lbuf) addSpaceOpt() {
	n := len(b.buf)
	if n == 0 {
		return
	}
	if b.buf[n-1] == ' ' { // last is space
		return
	}
	b.buf = append(b.buf, ' ')
}

func (b *lbuf) putStr(v ...interface{}) {
	for _, x := range v {
		b.addSpaceOpt()
		b.buf = fmt.Appendf(b.buf, "%v", x)
	}
}

func (b *lbuf) checkFatalPanic(lvl Level) {
	switch lvl {
	case FatalLevel:
		os.Exit(1)
	case PanicLevel:
		panic(string(b.buf))
	}
}

func (b *lbuf) needQuote() bool {
	for _, v := range b.buf {
		if v <= ' ' {
			return true
		}
	}
	return false
}
