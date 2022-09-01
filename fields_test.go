package xlog

import (
	"net"
	"os"
	"testing"
)

func TestFields(t *testing.T) {
	b := bufGet()
	defer bufPut(b)
	f := newFields()
	f.add("test_bool", true)
	f.add("test_ip", net.IPv4(1, 2, 3, 4))
	f.writeToBuf(b)
	t.Logf("%s\n", string(b.buf))
	h, err := os.Hostname()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(h)
}

func BenchmarkFields(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		bb := bufGet()
		f := getFields()
		f.add("test_bool", true)
		f.add("test_ip", net.IPv4(1, 2, 3, 4))
		f.writeToBuf(bb)
		putFields(f)
		bufPut(bb)
	}
}
