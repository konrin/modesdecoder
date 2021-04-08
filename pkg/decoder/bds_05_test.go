package decoder

import (
	"testing"
	"time"

	"github.com/konrin/modesdecoder/pkg/common"
)

func TestBDS05_AirbornePosition(t *testing.T) {
	msgEven, _ := common.NewMessage("8D40058B58C904A87F402D3B8C59", time.Now())
	msgOdd, _ := common.NewMessage("8D40058B58C901375147EFD09357", time.Now().Add(time.Second*5))

	bds := BDS05{}

	lat, lon, err := bds.AirbornePosition(msgEven.GetBinRaw(), msgEven.ReceiptAt, msgOdd.GetBinRaw(), msgOdd.ReceiptAt)
	if err != nil {
		t.Error(err)
		return
	}

	if lat != 49.81755 || lon != 6.08442 {
		t.Error("Позиция определена не корректно")
	}
}

func TestBDS05_PositionRef(t *testing.T) {
	type posRef struct {
		Message string
		Lat     float64
		Lon     float64
		ResLat  float64
		ResLon  float64
	}

	list := []posRef{
		{
			Message: "8D40058B58C901375147EFD09357",
			Lat:     49.0,
			Lon:     6.0,
			ResLat:  49.8241,
			ResLon:  6.06785,
		},
		{
			Message: "8D40058B58C904A87F402D3B8C59",
			Lat:     49.0,
			Lon:     6.0,
			ResLat:  49.81755,
			ResLon:  6.08442,
		},
	}

	bds := BDS05{}

	for _, pos := range list {
		msg, _ := common.NewMessage(pos.Message, time.Now())

		lat, lon, err := bds.AirbornePositionWithRef(msg.GetBinRaw(), msg.OE, pos.Lat, pos.Lon)
		if err != nil {
			t.Error(err)
			return
		}

		if lat != pos.ResLat || lon != pos.ResLon {
			t.Error("Позиция определена не корректно")
		}
	}
}

func TestBDS05_Altitude(t *testing.T) {
	msg, _ := common.NewMessage("8D40058B58C901375147EFD09357", time.Now())

	bds := BDS05{}

	if alt, err := bds.Altitude(msg.GetBinRaw(), 18); err != nil || alt != 39000 {
		if err != nil {
			t.Error(err)
			return
		}

		t.Error("Alt error")
	}
}
