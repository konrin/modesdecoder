package decoder

import (
	"math"

	"github.com/konrin/modesdecoder/pkg/common"
)

// ADS-B TC=19
// Aircraft Airborn velocity
type BDS09 struct{}

// Decode airborne velocity
func (m *BDS09) AirborneVelocity(bits *common.Bits) (
	// Speed (kt)
	speed float64,
	// Angle (degree), either ground track or heading
	track float64,
	// Vertical rate (ft/min)
	verticalRate int,
	// Speed type ('GS' for ground speed, 'AS' for airspeed)
	speedType string,
	// Direction source ('TRUE_NORTH' or 'MAGNETIC_NORTH')
	dirType string,
	err error,
) {
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

		speed = math.Sqrt(vSn*vSn + vWe*vWe)

		trk := math.Atan2(vWe, vSn)
		trk = trk * 180 / math.Pi
		if trk < 0 {
			trk = trk + 360
		}

		speedType = "GS"
		dirType = "TRUE_NORTH"
		track = trk
	} else {
		hdg := float64(bits.Int64(46, 56)) / 1024.0 * 360.0
		speed = float64(bits.Int64(57, 67))

		speedType = "AS"
		dirType = "MAGNETIC_NORTH"
		track = hdg
	}

	vrSign := 1
	if bits.At(68) == 1 {
		vrSign = -1
	}

	vr := int((bits.Int64(69, 78) - 1) * 64)
	verticalRate = vrSign * vr

	speed = common.Round(speed, .5, 1)
	track = common.Round(track, .5, 1)

	return
}
