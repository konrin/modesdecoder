package modesdecoder

import (
	"testing"
	"time"
)

func TestBDS02_CS(t *testing.T) {
	msg, _ := NewMessage("A000083E202CC371C31DE0AA1CCF", time.Now())

	bds := BDS20{}

	if data := bds.CS(msg.Bin); data != "KLM1017_" {
		t.Error("Номер рейса не распарсен")
	}
}
