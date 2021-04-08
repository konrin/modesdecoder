package decoder

import (
	"testing"
	"time"

	"github.com/konrin/modesdecoder/pkg/common"
)

func TestBDS40_Alt(t *testing.T) {
	msg, _ := common.NewMessage("A000029C85E42F313000007047D3", time.Now())

	bds := BDS40{}

	if bds.AltMcp(msg.GetBin()) != 3008 {
		t.Error("AltMcp error")
	}

	if bds.AltFms(msg.GetBin()) != 3008 {
		t.Error("AltFms error")
	}
}

func TestBDS40_Baro(t *testing.T) {
	msg, _ := common.NewMessage("A000029C85E42F313000007047D3", time.Now())

	bds := BDS40{}

	if b := bds.Baro(msg.GetBin()); b != 1020.0 {
		t.Error("")
	}
}
