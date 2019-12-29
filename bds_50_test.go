package modesdecoder

import (
	"testing"
	"time"
)

func TestBDS50_Roll(t *testing.T) {
	bds := BDS50{}

	for m, v := range map[string]float32{
		"A000139381951536E024D4CCF6B5": 2.1,
		"A0001691FFD263377FFCE02B2BF9": -0.4,
	} {
		msg := NewMessage(m, time.Now())

		val := bds.Roll(msg.Bin)
		if val != v {
			t.Error()
		}
	}
}

func TestBDS50_TRK(t *testing.T) {
	msg := NewMessage("A000139381951536E024D4CCF6B5", time.Now())

	bds := BDS50{}

	val := bds.TRK(msg.Bin)
	if val != 114.258 {
		t.Error()
	}
}

func TestBDS50_GS(t *testing.T) {
	msg := NewMessage("A000139381951536E024D4CCF6B5", time.Now())

	bds := BDS50{}

	val := bds.GS(msg.Bin)
	if val != 438 {
		t.Error()
	}
}

func TestBDS50_RTRK(t *testing.T) {
	msg := NewMessage("A000139381951536E024D4CCF6B5", time.Now())

	bds := BDS50{}

	val := bds.RTRK(msg.Bin)
	if val != .125 {
		t.Error()
	}
}

func TestBDS50_TAS(t *testing.T) {
	msg := NewMessage("A000139381951536E024D4CCF6B5", time.Now())

	bds := BDS50{}

	val := bds.TAS(msg.Bin)
	if val != 424 {
		t.Error()
	}
}
