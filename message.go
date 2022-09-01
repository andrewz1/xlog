package xlog

import (
	"time"
)

type message struct {
	msg       []byte
	timestamp time.Time
	level     Level
	fields    map[string][]byte
}
