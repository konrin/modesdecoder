package adsbdecoder

import "testing"

func TestBDS08_Callsign(t *testing.T) {
	ctx := NewMessageContext("8D406B902015A678D4D220AA4BDA")

	bds := BDS08{}

	if data, err := bds.Callsign(ctx); err != nil || data[CALLSING] != "EZY85MH_" {
		t.Error("Номер рейса не распарсен", err)
	}
}
