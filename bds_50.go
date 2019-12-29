package modesdecoder

type BDS50 struct{}

func (BDS50) Is(bits *Bits) bool {
	if Allzeros(bits) {
		return false
	}

	d := Data(bits)

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

func (BDS50) Roll(bits *Bits) float32 {
	d := Data(bits)

	if d.At(0) == 0 {
		return 0
	}

	val := d.Int64(2, 11)

	if d.At(1) > 0 {
		val = val - 512
	}

	angle := float64(val) * 45.0 / 256.0

	return float32(Round(angle, .5, 1))
}

func (BDS50) TRK(bits *Bits) float32 {
	d := Data(bits)

	if d.At(11) == 0 {
		return 0
	}

	val := d.Int64(13, 23)

	if d.At(12) > 0 {
		val = val - 1024
	}

	trk := float64(val) * 90.0 / 512.0
	if trk < 0 {
		trk = trk + 360
	}

	return float32(Round(trk, .5, 3))
}

func (BDS50) GS(bits *Bits) int {
	d := Data(bits)

	if d.At(23) == 0 {
		return 0
	}

	return int(d.Int64(24,34) * 2)
}

func (BDS50) RTRK(bits *Bits) float32 {
	d := Data(bits)

	if d.At(34) == 0 {
		return 0
	}

	if d.String(36, 45) == "111111111" {
		return 0
	}

	val := float64(d.Int64(36, 45))

	if d.At(35) > 0 {
		val = val - 512
	}

	angle := val * 8.0 / 256.0

	return float32(Round(angle, .5, 3))
}

func (BDS50) TAS(bits *Bits) int {
	d := Data(bits)
	if d.At(45) == 0 {
		return 0
	}

	return int(d.Int64(46, 56) * 2)
}
