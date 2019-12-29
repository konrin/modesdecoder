package modesdecoder

type BDS10 struct{}

func (BDS10) Is(bits *Bits) bool {
	if Allzeros(bits) {
		return false
	}

	d := Data(bits)

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
func (BDS10) OVC(bits *Bits) int {
	return int(Data(bits).At(14))
}
