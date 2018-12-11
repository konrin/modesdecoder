package modesdecoder

type BDS44 struct{}

func (bds BDS44) Is(bin []uint8) bool {

	if Allzeros(bin) {
		return false
	}

	d := Data(bin)

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

	if BinToInt(d[0:4]) > 4 {
		return false
	}

	v, w := bds.Wind(bin)
	if (v+w) != 0 && v > 250 {
		return false
	}

	temp := bds.Temp(bin)
	if temp == 0 || temp > 60 || temp < -80 {
		return false
	}

	return true
}

func (BDS44) Wind(bin []uint8) (float32, float32) {
	d := Data(bin)

	if d[4] == 0 {
		return 0, 0
	}

	speed := BinToInt(d[5:14])
	directions := BinToInt(d[14:23]) * 180.0 / 256.0

	return float32(Round(float64(speed), 0, 2)), float32(Round(float64(directions), 1, 2))
}

func (BDS44) Temp(bin []uint8) float32 {
	d := Data(bin)

	sign := d[23]
	val := BinToInt(d[24:34])

	if sign > 0 {
		val = val - 1024
	}

	temp := float64(val) * 0.125
	temp = Round(temp, 1, 2)

	return float32(temp)
}

func (BDS44) Pressure(bin []uint8) int {
	d := Data(bin)

	if d[34] == 0 {
		return 0
	}

	return int(BinToInt(d[35:46]))
}

func (BDS44) Hum(bin []uint8) float32 {
	d := Data(bin)

	if d[49] == 0 {
		return 0
	}

	return float32(BinToInt(d[50:56]) * 100.0 / 64)
}
