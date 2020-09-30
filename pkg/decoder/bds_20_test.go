package decoder

import (
	"testing"
	"time"

	"github.com/konrin/modesdecoder/pkg/common"
)

func TestBDS02_CS(t *testing.T) {
	msg, _ := common.NewMessage("A000083E202CC371C31DE0AA1CCF", time.Now())

	bds := BDS20{}

	if data := bds.CS(msg.GetBin()); data != "KLM1017_" {
		t.Error("Номер рейса не распарсен")
	}
}
