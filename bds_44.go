package modesdecoder

type BDS44 struct{}

func (bds BDS44) Is(bits *Bits) bool {

	if Allzeros(bits) {
		return false
	}

	d := Data(bits)

	if Wrongstatus(d, 5, 6, 23) {
		return false
	}

	if Wrongstatus(d, 35, 36, 46) {
		return false
	}

	if Wrongstatus(d, 47, 48, 49) {
		return false
	}

	if Wrongstatus(d, 50, 51, 56) {
		return false
	}

	if d.Int64(0, 4) > 4 {
		return false
	}

	v, w := bds.Wind(bits)
	if (v+w) != 0 && v > 250 {
		return false
	}

	temp := bds.Temp(bits)
	if temp == 0 || temp > 60 || temp < -80 {
		return false
	}

	return true
}

func (BDS44) Wind(bits *Bits) (float32, float32) {
	d := Data(bits)

	if d.At(4) == 0 {
		return 0, 0
	}

	speed := d.Int64(5, 14)
	directions := d.Int64(14, 23) * 180.0 / 256.0

	return float32(Round(float64(speed), 0, 2)),
		float32(Round(float64(directions), 1, 2))
}

func (BDS44) Temp(bits *Bits) float32 {
	d := Data(bits)

	sign := d.At(23)
	val := d.Int64(24, 34)

	if sign > 0 {
		val = val - 1024
	}

	temp := float64(val) * 0.125
	temp = Round(temp, 1, 2)

	return float32(temp)
}

func (BDS44) Pressure(bits *Bits) int {
	d := Data(bits)

	if d.At(34) == 0 {
		return 0
	}

	return int(d.Int64(35, 46))
}

func (BDS44) Hum(bits *Bits) float32 {
	d := Data(bits)

	if d.At(49) == 0 {
		return 0
	}

	return float32(d.Int64(50, 56) * 100.0 / 64)
}
