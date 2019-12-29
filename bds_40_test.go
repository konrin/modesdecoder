package modesdecoder

import (
	"testing"
	"time"
)

func TestBDS40_Alt(t *testing.T) {
	msg := NewMessage("A000029C85E42F313000007047D3", time.Now())

	bds := BDS40{}

	if mcp, fms := bds.Alt(msg.Bin); mcp != 3008 || fms != 3008 {
		t.Error("Номер рейса не распарсен")
	}
}

func TestBDS40_Baro(t *testing.T) {
	msg := NewMessage("A000029C85E42F313000007047D3", time.Now())

	bds := BDS40{}

	if b := bds.Baro(msg.Bin); b != 1020.0 {
		t.Error("Номер рейса не распарсен")
	}
}
