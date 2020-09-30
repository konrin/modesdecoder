package decoder

import "github.com/konrin/modesdecoder/pkg/common"

type BDS20 struct{}

func (bds *BDS20) Is(bits *common.Bits) bool {
	if common.Allzeros(bits) {
		return false
	}

	d := common.Data(bits)

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

func (BDS20) CS(bits *common.Bits) string {
	d := common.Data(bits)

	var cs string
	cs += string(common.Chars[d.Int64(8, 14)])
	cs += string(common.Chars[d.Int64(14, 20)])
	cs += string(common.Chars[d.Int64(20, 26)])
	cs += string(common.Chars[d.Int64(26, 32)])
	cs += string(common.Chars[d.Int64(32, 38)])
	cs += string(common.Chars[d.Int64(38, 44)])
	cs += string(common.Chars[d.Int64(44, 50)])
	cs += string(common.Chars[d.Int64(50, 56)])

	return cs
}
