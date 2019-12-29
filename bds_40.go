package modesdecoder

type BDS40 struct{}

func (BDS40) Is(bits *Bits) bool {
	if Allzeros(bits) {
		return false
	}

	d := Data(bits)

	if Wrongstatus(d, 1, 2, 13) {
		return false
	}

	if Wrongstatus(d, 4, 15, 26) {
		return false
	}

	if Wrongstatus(d, 27, 28, 39) {
		return false
	}

	if Wrongstatus(d, 48, 49, 51) {
		return false
	}

	if Wrongstatus(d, 54, 55, 56) {
		return false
	}

	if !d.IsZero(39,47) {
		return false
	}

	if !d.IsZero(51,53) {
		return false
	}

	return true
}

// Selected altitude, MCP/FCU
func (BDS40) Alt(bits *Bits) (mcp, fms int) {
	d := Data(bits)

	if d.At(0) != 0 {
		mcp = int(d.Int64(1,13) * 16) // ft
	}

	if d.At(13) != 0 {
		fms = int(d.Int64(14,26) * 16) // ft
	}

	return
}

// Barometric pressure setting
// pressure in millibar
func (BDS40) Baro(bits *Bits) float32 {
	d := Data(bits)

	if d.At(26) == 0 {
		return 0
	}

	return float32(d.Int64(27,39))*0.1 + 800
}
