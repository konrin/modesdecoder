package decoder

import (
	"testing"
	"time"

	"github.com/konrin/modesdecoder/pkg/common"
)

func TestBDS08_Callsign(t *testing.T) {
	msg, _ := common.NewMessage("8D406B902015A678D4D220AA4BDA", time.Now())

	bds := BDS08{}

	if data := bds.Callsign(msg.GetBin()); data != "EZY85MH" {
		t.Error("Flight number is not parsed " + data)
	}
}
