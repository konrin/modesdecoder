package adsbdecoder

import (
	"errors"
	"math"
	"strings"
)

type BDS05 struct {
	BDS
}

func NewBDS05() BDS05 {
	return BDS05{}
}

func (BDS05) Is(ctx *MessageContext) bool {
	return true
}

func (BDS05) Decode(ctx *MessageContext) (FlightData, error) {
	data := make(FlightData)

	return data, nil
}

func (BDS05) AirbornePosition(ctx *MessageContext) (FlightData, error) {
	data := make(FlightData)

	msg0, msg1, err := ShuffleFlagMessage(ctx, ctx.LastAirEvenMsg, ctx.LastAirOddMsg)
	if err != nil {
		return data, err
	}

	bin2 := strings.Join(msg0.GetBin(), "")
	bin1 := strings.Join(msg1.GetBin(), "")

	cprlatEven := float64(MustBinToInt(bin1[54:71])) / 131072.0
	cprlonEven := float64(MustBinToInt(bin1[71:88])) / 131072.0
	cprlatOdd := float64(MustBinToInt(bin2[54:71])) / 131072.0
	cprlonOdd := float64(MustBinToInt(bin2[71:88])) / 131072.0

	airDLatEven := 360.0 / 60
	airDLatOdd := 360.0 / 59

	j := math.Floor(59*cprlatEven - 60*cprlatOdd + 0.5)

	latEven := (airDLatEven * (Mod(j, 60) + cprlatEven))
	latOdd := (airDLatOdd * (Mod(j, 59) + cprlatOdd))

	if latEven >= 270 {
		latEven = latEven - 360
	}

	if latOdd >= 270 {
		latOdd = latOdd - 360
	}

	if cprNLEven, cprNLOdd := CprNL(latEven), CprNL(latOdd); cprNLEven != cprNLOdd {
		return data, errors.New("cprNLEven != cprNLOdd")
	}

	var lat, lon float64

	if msg0.GetTime().Unix() > msg1.GetTime().Unix() {
		lat = latEven
		nl := CprNL(lat)
		ni := math.Max(nl-0, 1)
		w := math.Floor(cprlonEven*(nl-1) - cprlonOdd*nl + 0.5)
		lon = (360.0 / ni) * (Mod(w, ni) + cprlonEven)
	} else {
		lat = latOdd
		nl := CprNL(lat)
		ni := math.Max(nl-1, 1)
		w := math.Floor(cprlonEven*(nl-1) - cprlonOdd*nl + 0.5)
		lon = (360.0 / ni) * (Mod(w, ni) + cprlonOdd)
	}

	if lon > 180 {
		lon = lon - 360
	}

	lat = Round(lat, .5, 5)
	lon = Round(lon, .5, 5)

	data[AIRBORN_LAT] = lat
	data[AIRBORN_LON] = lon

	return data, nil
}

func (BDS05) AirbornePositionWithRef(ctx *MessageContext) (FlightData, error) {
	data := make(FlightData)

	var lat, lon float64

	if ctx.LastAirPosition == nil {
		return data, errors.New("Not found air ref position")
	}

	i := 0
	if ctx.GetOEFlag() {
		i = 1
	}

	dLat := 360.0 / 60
	if i == 1 {
		dLat = 360.0 / 59
	}

	bin := strings.Join(ctx.GetBin(), "")

	cprLat := float64(MustBinToInt(bin[54:71])) / 131072.0
	cprLon := float64(MustBinToInt(bin[71:88])) / 131072.0

	j := math.Floor(ctx.LastAirPosition.Lat/dLat) + math.Floor(0.5+(Mod(ctx.LastAirPosition.Lat, dLat)/dLat)-cprLat)

	lat = dLat * (j + cprLat)

	ni := CprNL(lat) - float64(i)

	dLon := 360.0
	if ni > 0 {
		dLon = 360.0 / ni
	}

	w := math.Floor(ctx.LastAirPosition.Lon/dLon) + math.Floor(0.5+(Mod(ctx.LastAirPosition.Lon, dLon)/dLon)-cprLon)

	lon = dLon * (w + cprLon)

	lat = Round(lat, .5, 5)
	lon = Round(lon, .5, 5)

	data[AIRBORN_LAT] = lat
	data[AIRBORN_LON] = lon

	return data, nil
}

func (BDS05) Altitude(ctx *MessageContext) (FlightData, error) {
	data := make(FlightData)

	if ctx.GetTypeCode() >= 5 && ctx.GetTypeCode() <= 8 {
		return data, errors.New("")
	}

	if ctx.GetBin()[47] != "1" {
		return data, errors.New("")
	}

	bin := strings.Join(ctx.GetBin(), "")

	n := int(MustBinToInt(bin[40:47] + bin[48:52]))

	data[ALTITUDE] = n*25 - 1000

	return data, nil
}
