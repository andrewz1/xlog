package xlog

import (
	"fmt"
)

func emsg2(lvl Level, op int, e *Entry, v ...interface{}) {
	tb := getBuf()
	defer putBuf(tb)
	for _, x := range v {
		tb.putAny(x)
	}
	msg(lvl, op, tb.String(), e)
}

func emsgf2(lvl Level, op int, e *Entry, format string, v ...interface{}) {
	tb := getBuf()
	defer putBuf(tb)
	fmt.Fprintf(tb, format, v...)
	msg(lvl, op, tb.String(), e)
}

func msg2(lvl Level, op int, v ...interface{}) {
	emsg2(lvl, op, nil, v...)
}

func msgf2(lvl Level, op int, format string, v ...interface{}) {
	emsgf2(lvl, op, nil, format, v...)
}

func Msg(lvl Level, v ...interface{}) {
	msg2(lvl, noOp, v...)
}

func Msgf(lvl Level, format string, v ...interface{}) {
	msgf2(lvl, noOp, format, v...)
}

func Print(v ...interface{}) {
	Msg(noLevel, v...)
}

func Printf(format string, v ...interface{}) {
	Msgf(noLevel, format, v...)
}

func Emerg(v ...interface{}) {
	Msg(EmergLevel, v...)
}

func Emergf(format string, v ...interface{}) {
	Msgf(EmergLevel, format, v...)
}

func Alert(v ...interface{}) {
	Msg(AlertLevel, v...)
}

func Alertf(format string, v ...interface{}) {
	Msgf(AlertLevel, format, v...)
}

func Crit(v ...interface{}) {
	Msg(CritLevel, v...)
}

func Critf(format string, v ...interface{}) {
	Msgf(CritLevel, format, v...)
}

func Error(v ...interface{}) {
	Msg(ErrLevel, v...)
}

func Errorf(format string, v ...interface{}) {
	Msgf(ErrLevel, format, v...)
}

func Warn(v ...interface{}) {
	Msg(WarningLevel, v...)
}

func Warnf(format string, v ...interface{}) {
	Msgf(WarningLevel, format, v...)
}

func Notice(v ...interface{}) {
	Msg(NoticeLevel, v...)
}

func Noticef(format string, v ...interface{}) {
	Msgf(NoticeLevel, format, v...)
}

func Info(v ...interface{}) {
	Msg(InfoLevel, v...)
}

func Infof(format string, v ...interface{}) {
	Msgf(InfoLevel, format, v...)
}

func Debug(v ...interface{}) {
	Msg(InfoLevel, v...)
}

func Debugf(format string, v ...interface{}) {
	Msgf(InfoLevel, format, v...)
}

func Fatal(v ...interface{}) {
	msg2(EmergLevel, opFatal, v...)
}

func Fatalf(format string, v ...interface{}) {
	msgf2(EmergLevel, opFatal, format, v...)
}

func Panic(v ...interface{}) {
	msg2(EmergLevel, opPanic, v...)
}

func Panicf(format string, v ...interface{}) {
	msgf2(EmergLevel, opPanic, format, v...)
}

///

func (e *Entry) Msg(lvl Level, v ...interface{}) {
	emsg2(lvl, noOp, e, v...)
}

func (e *Entry) Msgf(lvl Level, format string, v ...interface{}) {
	emsgf2(lvl, noOp, e, format, v...)
}

func (e *Entry) Print(v ...interface{}) {
	e.Msg(noLevel, v...)
}

func (e *Entry) Printf(format string, v ...interface{}) {
	e.Msgf(noLevel, format, v...)
}

func (e *Entry) Emerg(v ...interface{}) {
	e.Msg(EmergLevel, v...)
}

func (e *Entry) Emergf(format string, v ...interface{}) {
	e.Msgf(EmergLevel, format, v...)
}

func (e *Entry) Alert(v ...interface{}) {
	e.Msg(AlertLevel, v...)
}

func (e *Entry) Alertf(format string, v ...interface{}) {
	e.Msgf(AlertLevel, format, v...)
}

func (e *Entry) Crit(v ...interface{}) {
	e.Msg(CritLevel, v...)
}

func (e *Entry) Critf(format string, v ...interface{}) {
	e.Msgf(CritLevel, format, v...)
}

func (e *Entry) Error(v ...interface{}) {
	e.Msg(ErrLevel, v...)
}

func (e *Entry) Errorf(format string, v ...interface{}) {
	e.Msgf(ErrLevel, format, v...)
}

func (e *Entry) Warn(v ...interface{}) {
	e.Msg(WarningLevel, v...)
}

func (e *Entry) Warnf(format string, v ...interface{}) {
	e.Msgf(WarningLevel, format, v...)
}

func (e *Entry) Notice(v ...interface{}) {
	e.Msg(NoticeLevel, v...)
}

func (e *Entry) Noticef(format string, v ...interface{}) {
	e.Msgf(NoticeLevel, format, v...)
}

func (e *Entry) Info(v ...interface{}) {
	e.Msg(InfoLevel, v...)
}

func (e *Entry) Infof(format string, v ...interface{}) {
	e.Msgf(InfoLevel, format, v...)
}

func (e *Entry) Debug(v ...interface{}) {
	e.Msg(InfoLevel, v...)
}

func (e *Entry) Debugf(format string, v ...interface{}) {
	e.Msgf(InfoLevel, format, v...)
}

func (e *Entry) Fatal(v ...interface{}) {
	emsg2(EmergLevel, opFatal, e, v...)
}

func (e *Entry) Fatalf(format string, v ...interface{}) {
	emsgf2(EmergLevel, opFatal, e, format, v...)
}

func (e *Entry) Panic(v ...interface{}) {
	emsg2(EmergLevel, opPanic, e, v...)
}

func (e *Entry) Panicf(format string, v ...interface{}) {
	emsgf2(EmergLevel, opPanic, e, format, v...)
}
