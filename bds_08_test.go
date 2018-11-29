package adsbdecoder

import (
	"testing"
	"time"
)

func TestBDS08_Callsign(t *testing.T) {
	msg := NewMessage("8D406B902015A678D4D220AA4BDA", time.Now())

	bds := BDS08{}

	if data := bds.Callsign(msg); data != "EZY85MH_" {
		t.Error("Номер рейса не распарсен")
	}
}
