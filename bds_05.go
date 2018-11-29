package adsbdecoder

import (
	"errors"
	"math"
)

type BDS05 struct{}

func (BDS05) AirbornePosition(msgEven *Message, msgOdd *Message) (float64, float64, error) {
	cprlatEven := float64(BinToInt(msgOdd.Bin[54:71])) / 131072.0
	cprlonEven := float64(BinToInt(msgOdd.Bin[71:88])) / 131072.0
	cprlatOdd := float64(BinToInt(msgEven.Bin[54:71])) / 131072.0
	cprlonOdd := float64(BinToInt(msgEven.Bin[71:88])) / 131072.0

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
		return 0, 0, errors.New("cprNLEven != cprNLOdd")
	}

	var lat, lon float64

	if msgEven.ReceiptAt.Unix() > msgOdd.ReceiptAt.Unix() {
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

	return lat, lon, nil
}

func (BDS05) AirbornePositionWithRef(msg *Message, latRef, lonRef float64) (float64, float64, error) {
	var lat, lon float64

	i := 0
	if msg.OE {
		i = 1
	}

	dLat := 360.0 / 60
	if i == 1 {
		dLat = 360.0 / 59
	}

	cprLat := float64(BinToInt(msg.Bin[54:71])) / 131072.0
	cprLon := float64(BinToInt(msg.Bin[71:88])) / 131072.0

	j := math.Floor(latRef/dLat) + math.Floor(0.5+(Mod(latRef, dLat)/dLat)-cprLat)

	lat = dLat * (j + cprLat)

	ni := CprNL(lat) - float64(i)

	dLon := 360.0
	if ni > 0 {
		dLon = 360.0 / ni
	}

	w := math.Floor(lonRef/dLon) + math.Floor(0.5+(Mod(lonRef, dLon)/dLon)-cprLon)

	lon = dLon * (w + cprLon)

	lat = Round(lat, .5, 5)
	lon = Round(lon, .5, 5)

	return lat, lon, nil
}

func (BDS05) Altitude(msg *Message) (int, error) {
	if msg.Bin[47] != 1 {
		return 0, errors.New("")
	}

	n := int(BinToInt(append(msg.Bin[40:47], msg.Bin[48:52]...)))

	return n*25 - 1000, nil
}
