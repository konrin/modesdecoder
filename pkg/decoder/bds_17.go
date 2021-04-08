package decoder

import "github.com/konrin/modesdecoder/pkg/common"

var AllBDS = []string{
	"05", "06", "07", "08", "09", "0A", "20", "21", "40", "41",
	"42", "43", "44", "45", "48", "50", "51", "52", "53", "54",
	"55", "56", "5F", "60", "NA", "NA", "E1", "E2",
}

// Common usage GICB capability report
type BDS17 struct{}

func (BDS17) Is(bits *common.Bits) bool {
	if common.Allzeros(bits) {
		return false
	}

	d := common.Data(bits)

	if !d.IsZero(28, 56) {
		return false
	}

	return true
}

// Extract capacities from BDS 1,7 message
// returns list of support BDS codes
func (BDS17) Cap(bits *common.Bits) []string {
	d := common.Data(bits)

	capacity := []string{}

	idx := []int{}
	for i, v := range d.String(0, 28) {
		if v == 1 {
			idx = append(idx, i)
		}
	}

	for _, v := range idx {
		if AllBDS[v] == "NA" {
			continue
		}

		capacity = append(capacity, "BDS"+AllBDS[v])
	}

	return capacity
}
