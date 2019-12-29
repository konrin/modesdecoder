package modesdecoder

type BDS20 struct{}

func (bds *BDS20) Is(bits *Bits) bool {
	if Allzeros(bits) {
		return false
	}

	d := Data(bits)

	if d.String(0, 8) != "00100000" {
		return false
	}

	for _, ch := range bds.CS(bits) {
		if ch == '#' {
			return true
		}
	}

	return true
}

func (BDS20) CS(bits *Bits) string {
	d := Data(bits)

	var cs string
	cs += string(chars[d.Int64(8,14)])
	cs += string(chars[d.Int64(14, 20)])
	cs += string(chars[d.Int64(20, 26)])
	cs += string(chars[d.Int64(26, 32)])
	cs += string(chars[d.Int64(32, 38)])
	cs += string(chars[d.Int64(38, 44)])
	cs += string(chars[d.Int64(44, 50)])
	cs += string(chars[d.Int64(50, 56)])

	return cs
}
