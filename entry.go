package xlog

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	entName = "entry"
)

type Entry struct {
	sync.RWMutex
	fld  map[string]any    // entry fields
	gfld map[string]string // gelf rendered fields
	cur  time.Time         // alloc time
	sil  atomic.Bool       // no log on put
	buf  *lbuf             // rendered fields buffer
}

var (
	ep = sync.Pool{New: newEntryAny}
)

func newEntryAny() any {
	return &Entry{
		fld:  make(map[string]any),
		gfld: make(map[string]string),
	}
}

func getEntry() *Entry {
	e := ep.Get().(*Entry)
	e.cur = time.Now()
	e.buf = getBuf()
	return e
}

func putEntry(e *Entry) time.Duration {
	rv := time.Since(e.cur)
	e.reset()
	putBuf(e.buf)
	ep.Put(e)
	return rv
}

func (e *Entry) reset() { // reset entry for reuse
	for k := range e.fld {
		delete(e.fld, k)
	}
	for k := range e.gfld {
		delete(e.gfld, k)
	}
	e.cur = time.Time{}
	e.sil.Store(false)
	e.buf.reset()
}

func (e *Entry) Reset() { // reset entry for reuse
	e.Lock()
	defer e.Unlock()
	e.reset()
}

func GetEntry(name string) *Entry {
	e := getEntry()
	e.addField(entName, name)
	return e
}

func GetEmptyEntry() *Entry {
	return getEntry()
}

func PutEntryMsg(e *Entry, m string) time.Duration {
	if !e.sil.Load() {
		msg(defLevel, noOp, m, e)
	}
	return putEntry(e)
}

func PutEntry(e *Entry) time.Duration {
	return PutEntryMsg(e, "")
}

func PutEntrySilent(e *Entry) time.Duration {
	return putEntry(e)
}

func (e *Entry) Used() { // mark entry used - don't print on put
	e.sil.Store(true)
}

func (e *Entry) UnUsed() { // mark entry unused - print on put
	e.sil.Store(false)
}

func (e *Entry) Caller(name string) {
	e.Lock()
	defer e.Unlock()
	e.AddField(entName, name)
}

func (e *Entry) toBuf(bb *lbuf) {
	tb := getBuf()
	defer putBuf(tb)
	for k, v := range e.fld {
		tb.setAny(v)
		bb.putStr(tb.kvString(k))
	}
}

func (e *Entry) rebuildBuf() {
	e.buf.reset()
	e.toBuf(e.buf)
}

func (e *Entry) addField(k string, v any) {
	tb := getBuf()
	defer putBuf(tb)
	_, rebuild := e.fld[k] // is field exists?
	e.fld[k] = v
	tb.setAny(v)
	if rebuild {
		e.rebuildBuf()
	} else {
		e.buf.putStr(tb.kvString(k))
	}
}

func (e *Entry) AddField(k string, v any) {
	e.Lock()
	defer e.Unlock()
	e.addField(k, v)
}

func (e *Entry) delField(k string) {
	if _, exists := e.fld[k]; !exists {
		return
	}
	delete(e.fld, k)
	e.rebuildBuf()
}

func (e *Entry) DelField(k string) {
	e.Lock()
	defer e.Unlock()
	e.delField(k)
}
