package xlog

import (
	"fmt"
	"io"
	"testing"
)

func TestLog(t *testing.T) {
	Printf("test: %v", true)
	Print("test", true, false, 1.5)

	Infof("test: %v", true)
	Info("test", true, false, 1.5)

	fmt.Println(true, false, 1.5)
}

func BenchmarkLog(b *testing.B) {
	ll.wr = io.Discard
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("very long string test:", true, false, 1.5)
	}
}

func BenchmarkLogf(b *testing.B) {
	ll.wr = io.Discard
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Infof("very long string test: %v", true)
	}
}
