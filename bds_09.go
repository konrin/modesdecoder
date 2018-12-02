package modesdecoder

import (
	"math"
)

type BDS09 struct{}

func (m *BDS09) AirborneVelocity(bin []uint8) (spd float64, track float64, rocd int, tag string, err error) {
	subType := BinToInt(bin[37:40])

	if BinToInt(bin[46:56]) == 0 || BinToInt(bin[57:67]) == 0 {
		return
	}

	if subType == 1 || subType == 2 {
		vEwSign := 1
		if bin[45] == 1 {
			vEwSign = -1
		}
		vEw := int(BinToInt(bin[46:56])) - 1

		vNsSign := 1
		if bin[56] == 1 {
			vNsSign = -1
		}
		vNs := int(BinToInt(bin[57:67])) - 1

		vWe := float64(vEwSign * vEw)
		vSn := float64(vNsSign * vNs)

		spd = math.Sqrt(vSn*vSn + vWe*vWe)

		trk := math.Atan2(vWe, vSn)
		trk = trk * 180 / math.Pi
		if trk < 0 {
			trk = trk + 360
		}

		tag = "GS"
		track = trk
	} else {
		hdg := float64(BinToInt(bin[46:56])) / 1024.0 * 360.0
		spd = float64(BinToInt(bin[57:67]))

		tag = "AS"
		track = hdg
	}

	vrSign := 1
	if bin[68] == 1 {
		vrSign = -1
	}

	vr := int((BinToInt(bin[69:78]) - 1) * 64)
	rocd = (vrSign * vr)

	spd = Round(spd, .5, 1)
	track = Round(track, .5, 1)

	return
}
