package xlog

import (
	"fmt"
	"log/syslog"
)

const (
	noLevel Level = iota - 1
	PanicLevel
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel

	defLevel = InfoLevel
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
	lvl2log = map[Level]syslog.Priority{
		DebugLevel: syslog.LOG_DEBUG,
		InfoLevel:  syslog.LOG_INFO,
		WarnLevel:  syslog.LOG_WARNING,
		ErrorLevel: syslog.LOG_ERR,
		FatalLevel: syslog.LOG_CRIT,
		PanicLevel: syslog.LOG_EMERG,
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
