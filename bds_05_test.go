package adsbdecoder

import (
	"testing"
	"time"
)

func TestBDS05_AirbornePosition(t *testing.T) {
	msgEven := NewMessageContext("8D40058B58C904A87F402D3B8C59")
	msgOdd := NewMessageContext("8D40058B58C901375147EFD09357")

	msgOdd.SetTime(time.Now().Add(time.Second * 5))

	msgEven.SetLastAirPositionMessage(msgOdd, false)

	bds := NewBDS05()

	data, err := bds.AirbornePosition(msgEven)
	if err != nil {
		t.Error(err)
		return
	}

	if len(data) != 2 {
		t.Errorf("len data != 2 => %d", len(data))
		return
	}

	if data[AIRBORN_LAT] != 49.81755 || data[AIRBORN_LON] != 6.08442 {
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
		posRef{
			Message: "8D40058B58C901375147EFD09357",
			Lat:     49.0,
			Lon:     6.0,
			ResLat:  49.8241,
			ResLon:  6.06785,
		},
		posRef{
			Message: "8D40058B58C904A87F402D3B8C59",
			Lat:     49.0,
			Lon:     6.0,
			ResLat:  49.81755,
			ResLon:  6.08442,
		},
	}

	bds := BDS05{}

	for _, pos := range list {
		ctx := NewMessageContext(pos.Message)
		ctx.SetLastAirPosition(&GeoPosition{pos.Lat, pos.Lon})

		data, err := bds.AirbornePositionWithRef(ctx)
		if err != nil {
			t.Error(err)
			return
		}

		if data[AIRBORN_LAT] != pos.ResLat || data[AIRBORN_LON] != pos.ResLon {
			t.Error("Позиция определена не корректно")
		}
	}
}

func TestBDS05_Altitude(t *testing.T) {
	ctx := NewMessageContext("8D40058B58C901375147EFD09357")

	bds := BDS05{}

	if data, err := bds.Altitude(ctx); err != nil || data[ALTITUDE] != 39000 {
		if err != nil {
			t.Error(err)
			return
		}

		t.Error("Высота распарсена не корректно")
	}
}
