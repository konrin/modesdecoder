package decoder

import "github.com/konrin/modesdecoder/pkg/common"

type BDS10 struct{}

func (BDS10) Is(bits *common.Bits) bool {
	if common.Allzeros(bits) {
		return false
	}

	d := common.Data(bits)

	if bits.String(0, 8) != "00010000" {
		return false
	}

	if !bits.IsZero(9, 14) {
		return false
	}

	if d.At(14) == 1 && bits.Int64(16, 23) < 5 {
		return false
	}

	if d.At(14) == 0 && bits.Int64(16, 23) > 4 {
		return false
	}

	return true
}

// OVC returning whether the transponder is OVC capable
func (BDS10) OVC(bits *common.Bits) int {
	return int(common.Data(bits).At(14))
}
