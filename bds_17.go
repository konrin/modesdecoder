package adsbdecoder

var AllBDS = []string{
	"05", "06", "07", "08", "09", "0A", "20", "21", "40", "41",
	"42", "43", "44", "45", "48", "50", "51", "52", "53", "54",
	"55", "56", "5F", "60", "NA", "NA", "E1", "E2",
}

type BDS17 struct{}

func (BDS17) Is(bin []uint8) bool {
	if Allzeros(bin) {
		return false
	}

	d := Data(bin)

	if BinToInt(d[28:56]) != 0 {
		return false
	}

	return true
}

func (BDS17) Cap(bin []uint8) []string {
	d := Data(bin)

	capacity := []string{}

	idx := []int{}
	for i, v := range d[0:28] {
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
