package decoder

import (
	"math"

	"github.com/konrin/modesdecoder/pkg/common"
)

type BDS09 struct{}

func (m *BDS09) AirborneVelocity(bits *common.Bits) (spd float64, track float64, rocd int, tag string, err error) {
	subType := bits.Int64(37, 40)

	if bits.IsZero(46, 56) || bits.IsZero(57, 67) {
		return
	}

	if subType == 1 || subType == 2 {
		vEwSign := 1
		if bits.At(45) == 1 {
			vEwSign = -1
		}
		vEw := int(bits.Int64(46, 56)) - 1

		vNsSign := 1
		if bits.At(56) == 1 {
			vNsSign = -1
		}
		vNs := int(bits.Int64(57, 67)) - 1

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
		hdg := float64(bits.Int64(46, 56)) / 1024.0 * 360.0
		spd = float64(bits.Int64(57, 67))

		tag = "AS"
		track = hdg
	}

	vrSign := 1
	if bits.At(68) == 1 {
		vrSign = -1
	}

	vr := int((bits.Int64(69, 78) - 1) * 64)
	rocd = vrSign * vr

	spd = common.Round(spd, .5, 1)
	track = common.Round(track, .5, 1)

	return
}
