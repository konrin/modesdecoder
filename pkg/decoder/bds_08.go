package decoder

import (
	"strings"

	"github.com/konrin/modesdecoder/pkg/common"
)

// // BDS08 BDS 0,8
// // ADS-B TC=1-4
// // Aircraft identitification and category
type BDS08 struct{}

func (BDS08) Callsign(bits *common.Bits) string {
	csbin := bits.Slice(40, 96)

	ct := [][]uint8{
		csbin.Slice(0, 6).Raw(), csbin.Slice(6, 12).Raw(),
		csbin.Slice(12, 18).Raw(), csbin.Slice(18, 24).Raw(),
		csbin.Slice(24, 30).Raw(), csbin.Slice(30, 36).Raw(),
		csbin.Slice(36, 42).Raw(), csbin.Slice(42, 48).Raw(),
	}

	var cs string

	for i := range ct {
		cs += string(common.Chars[common.BinToInt(ct[i])])
	}

	cs = strings.Replace(cs, "#", "", -1)
	cs = strings.Replace(cs, "_", "", -1)

	return cs
}

func (BDS08) Category(bits *common.Bits) uint {
	return uint(bits.Int64(5, 8))
}
