package xlog

import (
	"fmt"
	"sort"
	"sync"
)

type fields struct {
	sync.RWMutex
	m map[string]interface{}
}

var (
	fp = sync.Pool{
		New: func() interface{} {
			return newFields()
		},
	}
)

func newFields() *fields {
	return &fields{m: make(map[string]interface{})}
}

func getFields() *fields {
	return fp.Get().(*fields)
}

func putFields(f *fields) {
	f.clear()
	fp.Put(f)
}

func (f *fields) add(k string, v interface{}) {
	f.Lock()
	f.m[k] = v
	f.Unlock()
}

func (f *fields) del(k string) {
	f.Lock()
	delete(f.m, k)
	f.Unlock()
}

func (f *fields) clear() {
	f.Lock()
	for k := range f.m {
		delete(f.m, k)
	}
	f.Unlock()
}

func (f *fields) writeToBuf(b *lbuf) {
	f.RLock()
	bb := bufGet()
	defer func() {
		bufPut(bb)
		f.RUnlock()
	}()
	keys := make([]string, 0, len(f.m))
	for k := range f.m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var q bool
	for _, k := range keys {
		b.addSpaceOpt()
		bb.reset()
		bb.buf = fmt.Appendf(bb.buf, k)
		if q = bb.needQuote(); q {
			b.buf = append(b.buf, '"')
		}
		b.buf = append(b.buf, bb.buf...)
		if q {
			b.buf = append(b.buf, '"')
		}
		b.buf = append(b.buf, '=')
		bb.reset()
		bb.buf = fmt.Append(bb.buf, f.m[k])
		if q = bb.needQuote(); q {
			b.buf = append(b.buf, '"')
		}
		b.buf = append(b.buf, bb.buf...)
		if q {
			b.buf = append(b.buf, '"')
		}
	}
}
