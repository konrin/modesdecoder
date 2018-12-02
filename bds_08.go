package modesdecoder

import "strings"

// // BDS08 BDS 0,8
// // ADS-B TC=1-4
// // Aircraft identitification and category
type BDS08 struct{}

func (BDS08) Callsign(bin []uint8) string {
	csbin := bin[40:96]

	ct := [][]uint8{
		csbin[0:6], csbin[6:12], csbin[12:18], csbin[18:24],
		csbin[24:30], csbin[30:36], csbin[36:42], csbin[42:48],
	}

	var cs string

	for i := range ct {
		cs += string(chars[BinToInt(ct[i])])
	}

	cs = strings.Replace(cs, "#", "", -1)
	cs = strings.Replace(cs, "_", "", -1)

	return cs
}

func (BDS08) Category(bin []uint8) uint {
	return uint(BinToInt(bin[5:8]))
}
