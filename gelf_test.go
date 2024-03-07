package xlog

import (
	"testing"
	"time"
)

const (
	gelfHost = "172.31.82.222:12201"
)

func TestGelf(t *testing.T) {
	err := gelfInit("test-host[test-service]", gelfHost)
	if err != nil {
		t.Fatal(err)
	}
	gd := getGelfData(InfoLevel, noOp, nextSeq(), nil)
	if gd == nil {
		t.Fatal("getGelfData fail")
	}
	gd.setMsg("this is message")
	gd.addField("test1", "value1")
	gd.addField("test2", "value2")
	gd.addField("test3", "value3")
	gd.addField("test4", "value4")
	gd.send()
	time.Sleep(time.Second)
}

func TestGelf2(t *testing.T) {
	opt.GLogDst = gelfHost
	opt.baseInit()

	e := GetEntry("test entry")
	e.AddField("test1", "value1")
	e.AddField("test2", 2)
	e.AddField("test2", time.Now())
	e.AddField("test4", "value4")
	e.DelField("test4")

	e.Infof("test info message")

	e.Emerg("test info message", "tail")

	PutEntry(e)

	time.Sleep(time.Second)
}
