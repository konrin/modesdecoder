package modesdecoder

type BDS50 struct{}

func (BDS50) Is(bin []uint8) bool {
	if Allzeros(bin) {
		return false
	}

	d := Data(bin)

	if Wrongstatus(d, 1, 3, 11) {
		return false
	}

	if Wrongstatus(d, 12, 13, 23) {
		return false
	}

	if Wrongstatus(d, 24, 25, 34) {
		return false
	}

	if Wrongstatus(d, 35, 36, 45) {
		return false
	}

	if Wrongstatus(d, 46, 47, 56) {
		return false
	}

	return true
}

func (BDS50) Roll(bin []uint8) float32 {
	d := Data(bin)

	if d[0] == 0 {
		return 0
	}

	val := BinToInt(d[2:11])

	if d[1] > 0 {
		val = val - 512
	}

	angle := float64(val) * 45.0 / 256.0

	return float32(Round(angle, .5, 1))
}

func (BDS50) TRK(bin []uint8) float32 {
	d := Data(bin)

	if d[11] == 0 {
		return 0
	}

	val := BinToInt(d[13:23])

	if d[12] > 0 {
		val = val - 1024
	}

	trk := float64(val) * 90.0 / 512.0
	if trk < 0 {
		trk = trk + 360
	}

	return float32(Round(trk, .5, 3))
}

func (BDS50) GS(bin []uint8) int {
	d := Data(bin)

	if d[23] == 0 {
		return 0
	}

	return int(BinToInt(d[24:34]) * 2)
}

func (BDS50) RTRK(bin []uint8) float32 {
	d := Data(bin)

	if d[34] == 0 {
		return 0
	}

	if BinToString(d[36:45]) == "111111111" {
		return 0
	}

	val := float64(BinToInt(d[36:45]))

	if d[35] > 0 {
		val = val - 512
	}

	angle := val * 8.0 / 256.0

	return float32(Round(angle, .5, 3))
}

func (BDS50) TAS(bin []uint8) int {
	d := Data(bin)
	if d[45] == 0 {
		return 0
	}

	return int(BinToInt(d[46:56]) * 2)
}
