package xlog

import (
	"fmt"
	"log/syslog"
)

const (
	noLevel    = Level(-1)
	DebugLevel = Level(syslog.LOG_DEBUG)
	InfoLevel  = Level(syslog.LOG_INFO)
	WarnLevel  = Level(syslog.LOG_WARNING)
	ErrorLevel = Level(syslog.LOG_ERR)
	FatalLevel = Level(syslog.LOG_CRIT)
	PanicLevel = Level(syslog.LOG_EMERG)
	defLevel   = InfoLevel
)

type Level int

var (
	lvl2strL = map[Level]string{
		DebugLevel: "debug",
		InfoLevel:  "info",
		WarnLevel:  "warn",
		ErrorLevel: "error",
		FatalLevel: "fatal",
		PanicLevel: "panic",
	}
	lvl2strU = map[Level]string{
		DebugLevel: "DEBUG",
		InfoLevel:  "INFO",
		WarnLevel:  "WARN",
		ErrorLevel: "ERROR",
		FatalLevel: "FATAL",
		PanicLevel: "PANIC",
	}
	str2lvl = map[string]Level{
		"debug":   DebugLevel,
		"info":    InfoLevel,
		"warn":    WarnLevel,
		"warning": WarnLevel,
		"error":   ErrorLevel,
		"fatal":   FatalLevel,
		"panic":   PanicLevel,
	}
)

func (l Level) String() string {
	if str, ok := lvl2strL[l]; ok {
		return str
	}
	return "unknown"
}

func (l Level) LogString() string {
	if str, ok := lvl2strU[l]; ok {
		return str
	}
	return "UNKNOWN"
}

func ParseLevel(lvlStr string) (Level, error) {
	if lvl, ok := str2lvl[lvlStr]; ok {
		return lvl, nil
	}
	return defLevel, fmt.Errorf("unknown level: '%s', defaulting to '%s'", lvlStr, defLevel)
}
