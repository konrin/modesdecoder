package adsbdecoder

type BDS40 struct{}

func (BDS40) Is(bin []uint8) bool {
	if Allzeros(bin) {
		return false
	}

	d := Data(bin)

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

	if BinToInt(d[39:47]) != 0 {
		return false
	}

	if BinToInt(d[51:53]) != 0 {
		return false
	}

	return true
}

// Selected altitude, MCP/FCU
func (BDS40) Alt(bin []uint8) (mcp, fms int) {
	d := Data(bin)

	if d[0] != 0 {
		mcp = int(BinToInt(d[1:13]) * 16) // ft
	}

	if d[13] != 0 {
		fms = int(BinToInt(d[14:26]) * 16) // ft
	}

	return
}

// Barometric pressure setting
// pressure in millibar
func (BDS40) Baro(bin []uint8) float32 {
	d := Data(bin)

	if d[26] == 0 {
		return 0
	}

	return float32(BinToInt(d[27:39]))*0.1 + 800
}
