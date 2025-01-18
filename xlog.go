package xlog

import (
	"fmt"
	"log/syslog"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/andrewz1/xtoml"
)

const (
	noOp = iota
	opFatal
	opPanic
)

type Conf struct {
	Level    string `conf:"log.level"`    // log level
	SysLog   bool   `conf:"log.syslog"`   // use syslog
	Duration bool   `conf:"log.duration"` // log duration

	lvl Level
	tag string
	log *syslog.Writer
}

var (
	opt = &Conf{
		Level: "info",
		//SysLog:   true,
		Duration: true,
		lvl:      InfoLevel,
	}
	seq atomic.Uint64
)

func Init(xc *xtoml.XConf) (err error) {
	if err = xc.LoadConf(opt); err != nil {
		return err
	}
	opt.baseInit()
	return nil
}

func Init2(xc *xtoml.XConf, sysLog bool) error {
	opt.SysLog = sysLog
	return Init(xc)
}

func nextSeq() uint64 {
	return seq.Add(1)
}

func logTag() string {
	if n := strings.LastIndexByte(os.Args[0], '/'); n >= 0 && n < len(os.Args[0]) {
		return os.Args[0][n+1:]
	}
	return os.Args[0]
}

func (c *Conf) baseInit() {
	var err error
	if c.lvl, err = parseLevel(c.Level); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	c.tag = logTag()
	if c.SysLog {
		c.log, err = syslog.New(syslog.Priority(c.lvl)|syslog.LOG_DAEMON, c.tag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "use stderr, syslog open err: %v\n", err)
			c.log = nil
		}
	}
}

func msgSysLog(lvl Level, m string) {
	switch lvl {
	case EmergLevel:
		opt.log.Emerg(m)
	case AlertLevel:
		opt.log.Alert(m)
	case CritLevel:
		opt.log.Crit(m)
	case ErrLevel:
		opt.log.Err(m)
	case WarningLevel:
		opt.log.Warning(m)
	case NoticeLevel:
		opt.log.Notice(m)
	case InfoLevel:
		opt.log.Info(m)
	case DebugLevel:
		opt.log.Debug(m)
	default:
		opt.log.Write([]byte(m))
	}
}

func msgStdErr(lvl Level, m string) {
	if lvl.isValid() {
		fmt.Fprintf(os.Stderr, "%s %s\n", lvl.logString(), m)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", m)
	}
}

func msg(lvl Level, op int, m string, e *Entry) {
	nseq := nextSeq()
	bb := getBuf()
	defer putBuf(bb)
	bb.setStr(m)
	if e != nil {
		bb.putStr(e.buf.String())
		if opt.Duration {
			bb.putStr(fmt.Sprintf("duration=%v", time.Since(e.cur)))
		}
	}
	bb.putStr(fmt.Sprintf("seq=%d", nseq))
	s := bb.String()
	if opt.log != nil {
		msgSysLog(lvl, s)
	} else {
		msgStdErr(lvl, s)
	}
	switch op {
	case opFatal:
		os.Exit(1)
	case opPanic:
		panic(s)
	}
}
