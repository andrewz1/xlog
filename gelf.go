package xlog

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/valyala/fastjson"
)

const (
	gchLen    = 4096
	noDataMsg = "[empty]"
	sendTmo   = time.Second / 2
)

type gelfData struct {
	shortMsg string
	fullMsg  string
	fld      map[string]string
	cur      time.Time
	dur      time.Duration
	seq      uint64
	lvl      int // log level
	op       int // exit or panic
}

var (
	gp    = sync.Pool{New: newGelfData}
	gCh   chan *gelfData
	gSock *net.UDPConn
	gHost string
	ap    fastjson.ArenaPool
)

func getArena() *fastjson.Arena {
	return ap.Get()
}

func putArena(a *fastjson.Arena) {
	a.Reset()
	ap.Put(a)
}

func newGelfData() any {
	return &gelfData{fld: make(map[string]string)}
}

func (gd *gelfData) reset() {
	gd.shortMsg = ""
	gd.fullMsg = ""
	gd.dur = 0
	gd.op = noOp
	for k := range gd.fld {
		delete(gd.fld, k)
	}
}

func gelfOk() bool {
	return gSock != nil
}

func getGelf(lvl Level, op int, short string, seq uint64) *gelfData {
	if !gelfOk() {
		return nil
	}
	gd := gp.Get().(*gelfData)
	gd.shortMsg = short
	gd.seq = seq
	gd.lvl = int(lvl)
	gd.op = op
	return gd
}

func (gd *gelfData) send() {
	if gd != nil {
		gCh <- gd
	}
}

func putGelf(gd *gelfData) {
	if gd == nil {
		return
	}
	gd.reset()
	gp.Put(gd)
}

func (gd *gelfData) setFullMsg(m string) {
	if gd != nil {
		gd.fullMsg = m
	}
}

func (gd *gelfData) setTime(cur time.Time, dur time.Duration) {
	if gd != nil {
		gd.cur = cur
		gd.dur = dur
	}
}

func (gd *gelfData) addField(k, v string) {
	if gd == nil || k == "" {
		return
	}
	var kv string
	if k == "id" {
		kv = "_ID"
	} else {
		kv = "_" + k
	}
	gd.fld[kv] = v
}

func gelfInit(host, dst string) error {
	da, err := net.ResolveUDPAddr("udp", dst)
	if err != nil {
		return err
	}
	gSock, err = net.DialUDP("udp", nil, da)
	if err != nil {
		return err
	}
	gHost = host
	gCh = make(chan *gelfData, gchLen)
	go gchLog()
	return nil
}

func gchLog() {
	defer func() {
		_ = gSock.Close()
		gSock = nil
	}()
	for v := range gCh {
		v.process()
		switch v.op {
		case noOp:
			putGelf(v)
		case opFatal:
			time.Sleep(sendTmo) // for send
			os.Exit(1)
		case opPanic:
			time.Sleep(sendTmo) // for send
			panic(v.shortMsg)
		}
	}
}

func (gd *gelfData) getTS() string {
	ts := gd.cur.UnixMilli()
	return fmt.Sprintf("%d.%d", ts/1000, ts%1000)
}

func (gd *gelfData) process() {
	a := getArena()
	bb := bufGet()
	defer func() {
		bufPut(bb)
		putArena(a)
	}()
	v := a.NewObject()
	v.Set("version", a.NewString("1.1"))
	v.Set("host", a.NewString(gHost))
	if gd.shortMsg == "" {
		if en, ok := gd.fld["_"+entName]; ok {
			v.Set("short_message", a.NewString(fmt.Sprintf(entName+"=%s", en)))
		} else {
			v.Set("short_message", a.NewString(noDataMsg))
		}
	} else {
		v.Set("short_message", a.NewString(gd.shortMsg))
	}
	if gd.fullMsg != "" {
		v.Set("full_message", a.NewString(gd.fullMsg))
	}
	v.Set("timestamp", a.NewNumberString(gd.getTS()))
	v.Set("level", a.NewNumberInt(gd.lvl))
	for fk, fv := range gd.fld {
		v.Set(fk, a.NewString(fv))
	}
	if gd.dur != 0 {
		v.Set("_duration", a.NewString(gd.dur.String()))
	}
	v.Set("_seq", a.NewNumberString(fmt.Sprintf("%d", gd.seq)))
	_, _ = gSock.Write(v.MarshalTo(bb.buf))
}
