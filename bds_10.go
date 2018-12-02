package modesdecoder

type BDS10 struct{}

func (BDS10) Is(bin []uint8) bool {
	if Allzeros(bin) {
		return false
	}

	d := Data(bin)

	if BinToString(d[0:8]) != "00010000" {
		return false
	}

	if BinToInt(d[9:14]) != 0 {
		return false
	}

	if d[14] == 1 && BinToInt(d[16:23]) < 5 {
		return false
	}

	if d[14] == 0 && BinToInt(d[16:23]) > 4 {
		return false
	}

	return true
}

// OVC returning whether the transponder is OVC capable
func (BDS10) OVC(bin []uint8) int {
	return int(Data(bin)[14])
}
