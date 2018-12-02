package modesdecoder

type BDS20 struct{}

func (bds *BDS20) Is(bin []uint8) bool {
	if Allzeros(bin) {
		return false
	}

	d := Data(bin)

	if BinToString(d[0:8]) != "00100000" {
		return false
	}

	for _, ch := range bds.CS(bin) {
		if ch == '#' {
			return true
		}
	}

	return true
}

func (BDS20) CS(bin []uint8) string {
	d := Data(bin)

	var cs string
	cs += string(chars[BinToInt(d[8:14])])
	cs += string(chars[BinToInt(d[14:20])])
	cs += string(chars[BinToInt(d[20:26])])
	cs += string(chars[BinToInt(d[26:32])])
	cs += string(chars[BinToInt(d[32:38])])
	cs += string(chars[BinToInt(d[38:44])])
	cs += string(chars[BinToInt(d[44:50])])
	cs += string(chars[BinToInt(d[50:56])])

	return cs
}
