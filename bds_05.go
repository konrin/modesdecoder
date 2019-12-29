package modesdecoder

import (
	"errors"
	"math"
	"time"
)

type BDS05 struct{}

func (BDS05) AirbornePosition(binEven []uint8, timeEven time.Time, binOdd []uint8, timeOdd time.Time) (float64, float64, error) {
	cprlatEven := float64(BinToInt(binOdd[54:71])) / 131072.0
	cprlonEven := float64(BinToInt(binOdd[71:88])) / 131072.0
	cprlatOdd := float64(BinToInt(binEven[54:71])) / 131072.0
	cprlonOdd := float64(BinToInt(binEven[71:88])) / 131072.0

	airDLatEven := 360.0 / 60
	airDLatOdd := 360.0 / 59

	j := math.Floor(59*cprlatEven - 60*cprlatOdd + 0.5)

	latEven := airDLatEven * (Mod(j, 60) + cprlatEven)
	latOdd := airDLatOdd * (Mod(j, 59) + cprlatOdd)

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

	if timeEven.UnixNano() > timeOdd.UnixNano() {
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

func (BDS05) AirbornePositionWithRef(bin []uint8, oeFlag bool, latRef, lonRef float64) (float64, float64, error) {
	var lat, lon float64

	i := 0
	if oeFlag {
		i = 1
	}

	dLat := 360.0 / 60
	if i == 1 {
		dLat = 360.0 / 59
	}

	cprLat := float64(BinToInt(bin[54:71])) / 131072.0
	cprLon := float64(BinToInt(bin[71:88])) / 131072.0

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

func (BDS05) Altitude(bin []uint8, ts uint) (alt int, err error) {
	mb := bin[32:]

	if ts < 19 {
		// barometric altitude
		if mb[15] == 1 {
			n := int(BinToInt(append(mb[8:15], mb[16:20]...)))
			alt = n*25 - 1000
		}
	} else {
		// GNSS altitude, meters
		alt = int(BinToInt(mb[8:20]))
	}

	return
}
