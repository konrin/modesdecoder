package decoder

import (
	"testing"
	"time"

	"github.com/konrin/modesdecoder/pkg/common"
)

func TestBDS40_Alt(t *testing.T) {
	msg, _ := common.NewMessage("A000029C85E42F313000007047D3", time.Now())

	bds := BDS40{}

	if mcp, fms := bds.Alt(msg.GetBin()); mcp != 3008 || fms != 3008 {
		t.Error("")
	}
}

func TestBDS40_Baro(t *testing.T) {
	msg, _ := common.NewMessage("A000029C85E42F313000007047D3", time.Now())

	bds := BDS40{}

	if b := bds.Baro(msg.GetBin()); b != 1020.0 {
		t.Error("")
	}
}
