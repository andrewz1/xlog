package xlog

import (
	"fmt"
	"log/syslog"
	"testing"
)

func TestMsg(t *testing.T) {
	//var err error
	//opt.log, err = syslog.New(syslog.LOG_DEBUG, "")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//defer opt.log.Close()
	//
	//opt.log.Info("test123")

	for i := 0; i <= int(syslog.LOG_DEBUG); i++ {
		msg(Level(i), noOp, fmt.Sprintf("test level %d", i), nil)
	}
}
