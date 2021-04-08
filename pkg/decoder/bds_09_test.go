package decoder

import (
	"testing"
	"time"

	"github.com/konrin/modesdecoder/pkg/common"
)

func TestBDS09_AirborneVelocity(t *testing.T) {
	msg, _ := common.NewMessage("8D485020994409940838175B284F", time.Now())

	bds := BDS09{}

	speed, track, rocd, tag, _, err := bds.AirborneVelocity(msg.GetBin())
	if err != nil {
		t.Error(err)
		return
	}

	if speed != 159.2 || track != 182.9 || rocd != -832 || tag != "GS" {
		t.Error("oops")
	}

}
