package xlog

import (
	"fmt"
	"log/syslog"
	"strings"
)

const (
	noLevel      Level = -1
	EmergLevel   Level = Level(syslog.LOG_EMERG)
	AlertLevel   Level = Level(syslog.LOG_ALERT)
	CritLevel    Level = Level(syslog.LOG_CRIT)
	ErrLevel     Level = Level(syslog.LOG_ERR)
	WarningLevel Level = Level(syslog.LOG_WARNING)
	NoticeLevel  Level = Level(syslog.LOG_NOTICE)
	InfoLevel    Level = Level(syslog.LOG_INFO)
	DebugLevel   Level = Level(syslog.LOG_DEBUG)
	defLevel           = InfoLevel
)

type Level int

var (
	lvl2str = map[Level]string{
		EmergLevel:   "emerg",
		AlertLevel:   "alert",
		CritLevel:    "crit",
		ErrLevel:     "error",
		WarningLevel: "warn",
		NoticeLevel:  "notice",
		InfoLevel:    "info",
		DebugLevel:   "debug",
	}
	str2lvl = map[string]Level{
		"emerg":  EmergLevel,
		"alert":  AlertLevel,
		"crit":   CritLevel,
		"error":  ErrLevel,
		"warn":   WarningLevel,
		"notice": NoticeLevel,
		"info":   InfoLevel,
		"debug":  DebugLevel,
	}
)

func (l Level) String() string {
	if str, ok := lvl2str[l]; ok {
		return str
	}
	return "unknown"
}

func (l Level) logString() string {
	return strings.ToUpper(l.String())
}

func (l Level) isValid() bool {
	return l >= EmergLevel && l <= DebugLevel
}

func parseLevel(lvlStr string) (Level, error) {
	if lvl, ok := str2lvl[lvlStr]; ok {
		return lvl, nil
	}
	return defLevel, fmt.Errorf("unknown level: '%s', defaulting to '%s'", lvlStr, defLevel)
}
