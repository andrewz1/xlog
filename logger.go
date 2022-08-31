package xlog

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"
)

type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})

	Msg(lvl Level, v ...interface{})
	Msgf(lvl Level, format string, v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
}

type logger struct {
	wr  io.Writer // base output
	lvl Level     // base level
}

var (
	seq uint32
	ll  = &logger{
		wr:  os.Stderr,
		lvl: defLevel,
	}
)

func (l *logger) skip(lvl Level) bool {
	if lvl < 0 { // don't skip noLevel
		return false
	}
	return lvl > l.lvl
}

func (l *logger) Print(v ...interface{}) {
	l.Msg(noLevel, v...)
}

func (l *logger) Printf(format string, v ...interface{}) {
	l.Msgf(noLevel, format, v...)
}

func (l *logger) Msg(lvl Level, v ...interface{}) {
	if l.skip(lvl) || len(v) == 0 {
		return
	}
	b := bufGet()
	defer bufPut(b)
	if lvl >= PanicLevel {
		b.buf = fmt.Appendf(b.buf, "%s", lvl.LogString())
	}
	b.putStr(v...)
	b.addSpaceOpt()
	b.buf = fmt.Appendf(b.buf, "seq=%d\n", atomic.AddUint32(&seq, 1))
	l.wr.Write(b.buf)
	b.checkFatalPanic(lvl)
}

func (l *logger) Msgf(lvl Level, format string, v ...interface{}) {
	if l.skip(lvl) || len(format) == 0 {
		return
	}
	b := bufGet()
	defer bufPut(b)
	if lvl >= PanicLevel {
		b.buf = fmt.Appendf(b.buf, "%s", lvl.LogString())
	}
	b.addSpaceOpt()
	b.buf = fmt.Appendf(b.buf, format, v...)
	b.addSpaceOpt()
	b.buf = fmt.Appendf(b.buf, "seq=%d\n", atomic.AddUint32(&seq, 1))
	l.wr.Write(b.buf)
	b.checkFatalPanic(lvl)
}

func (l *logger) Panic(v ...interface{}) {
	l.Msg(PanicLevel, v...)
}

func (l *logger) Panicf(format string, v ...interface{}) {
	l.Msgf(PanicLevel, format, v...)
}

func (l *logger) Fatal(v ...interface{}) {
	l.Msg(FatalLevel, v...)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.Msgf(FatalLevel, format, v...)
}

func (l *logger) Error(v ...interface{}) {
	l.Msg(ErrorLevel, v...)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.Msgf(ErrorLevel, format, v...)
}

func (l *logger) Warn(v ...interface{}) {
	l.Msg(WarnLevel, v...)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.Msgf(WarnLevel, format, v...)
}

func (l *logger) Info(v ...interface{}) {
	l.Msg(InfoLevel, v...)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.Msgf(InfoLevel, format, v...)
}

func (l *logger) Debug(v ...interface{}) {
	l.Msg(DebugLevel, v...)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.Msgf(DebugLevel, format, v...)
}

func Print(v ...interface{}) {
	ll.Print(v...)
}

func Printf(format string, v ...interface{}) {
	ll.Printf(format, v...)
}

func Msg(lvl Level, v ...interface{}) {
	ll.Msg(lvl, v...)
}

func Msgf(lvl Level, format string, v ...interface{}) {
	ll.Msgf(lvl, format, v...)
}

func Panic(v ...interface{}) {
	ll.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	ll.Panicf(format, v...)
}

func Fatal(v ...interface{}) {
	ll.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	ll.Fatalf(format, v...)
}

func Error(v ...interface{}) {
	ll.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	ll.Errorf(format, v...)
}

func Warn(v ...interface{}) {
	ll.Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	ll.Warnf(format, v...)
}

func Info(v ...interface{}) {
	ll.Info(v...)
}

func Infof(format string, v ...interface{}) {
	ll.Infof(format, v...)
}

func Debug(v ...interface{}) {
	ll.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	ll.Debugf(format, v...)
}
