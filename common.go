package adsbdecoder

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	hexToBinLookup = map[string][]uint8{
		"0": []uint8{0, 0, 0, 0},
		"1": []uint8{0, 0, 0, 1},
		"2": []uint8{0, 0, 1, 0},
		"3": []uint8{0, 0, 1, 1},
		"4": []uint8{0, 1, 0, 0},
		"5": []uint8{0, 1, 0, 1},
		"6": []uint8{0, 1, 1, 0},
		"7": []uint8{0, 1, 1, 1},
		"8": []uint8{1, 0, 0, 0},
		"9": []uint8{1, 0, 0, 1},
		"a": []uint8{1, 0, 1, 0},
		"b": []uint8{1, 0, 1, 1},
		"c": []uint8{1, 1, 0, 0},
		"d": []uint8{1, 1, 0, 1},
		"e": []uint8{1, 1, 1, 0},
		"f": []uint8{1, 1, 1, 1},
	}
	crcGenerator = [25]uint8{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1}
	chars        = "#ABCDEFGHIJKLMNOPQRSTUVWXYZ#####_###############0123456789######"
)

type Common struct{}

func Hex2Bin(hex string) []uint8 {
	bin := []uint8{}

	for _, r := range strings.ToLower(hex) {
		c, ok := hexToBinLookup[string(r)]
		if !ok {
			return []uint8{}
		}

		bin = append(bin, c...)
	}

	return bin
}

// []uint8{0,0,1,0,1} => "00101"
func BinToString(bin []uint8) string {
	str := ""

	for i := range bin {
		str += strconv.Itoa(int(bin[i]))
	}

	return str
}

// "00101" => []uint8{0,0,1,0,1}
func StringToBin(bin string) []uint8 {
	sbin := []uint8{}

	for _, ch := range bin {
		i, err := strconv.Atoi(string(ch))
		if err != nil {
			return []uint8{}
		}

		sbin = append(sbin, uint8(i))
	}

	return sbin
}

func BinToInt(bin []uint8) int64 {
	i, err := strconv.ParseInt(BinToString(bin), 2, 64)
	if err != nil {
		return 0
	}

	return i
}

func HexToInt(hex string) int64 {
	var (
		i   int64
		err error
	)

	if i, err = strconv.ParseInt(hex, 16, 64); err != nil {
		panic(err)
	}

	return i
}

func CRC(bin []uint8, encode bool) []uint8 {
	if encode {
		bin = bin[:len(bin)-24]

		for i := 0; i < 24; i++ {
			bin = append(bin, 0)
		}
	}

	for i := 0; i < len(bin)-24; i++ {
		if bin[i] != 1 {
			continue
		}

		for ci, cv := range crcGenerator {
			vi := bin[i+ci]

			bin[i+ci] = vi ^ cv
		}
	}

	return bin[len(bin)-24:]
}

func Gray2Int(graystr []uint8) int64 {
	num := BinToInt(graystr)

	num ^= (num >> 8)
	num ^= (num >> 4)
	num ^= (num >> 2)
	num ^= (num >> 1)

	return num
}

func Round(val float64, roundOn float64, places int) float64 {
	var round float64

	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)

	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}

	return round / pow
}

func Mod(a, b float64) float64 {
	mode := math.Mod(a, b)
	if mode < 0 {
		mode = mode + b
	}

	return mode
}

func CprNL(lat float64) float64 {
	if lat == 0 {
		return 59
	}

	if lat == 87 || lat < -87 {
		return 2
	}

	if lat > 87 || lat < -87 {
		return 1
	}

	a := 1 - math.Cos(math.Pi/(2*15))
	b := math.Pow(math.Cos(math.Pi/180.0*math.Abs(lat)), 2)
	nl := 2 * math.Pi / math.Acos(1-a/b)

	return math.Floor(nl)
}

func ICAO(msg *Message) string {
	var addr string

	switch msg.DF {
	case 11, 17, 18:
		addr = msg.Hex[2:8]
		break
	case 0, 4, 5, 16, 20, 21, 24:
		coCrx := CRC(msg.GetBin(), true)

		c0 := BinToInt(coCrx)
		c1 := HexToInt(msg.Hex[len(msg.Hex)-6:])

		addr = fmt.Sprintf("%06X", c0^c1)
		break
	}

	return addr
}

func IDCODE(bin []uint8) string {
	C1 := bin[19]
	A1 := bin[20]
	C2 := bin[21]
	A2 := bin[22]
	C4 := bin[23]
	A4 := bin[24]
	// _ = bin[25]
	B1 := bin[26]
	D1 := bin[27]
	B2 := bin[28]
	D2 := bin[29]
	B4 := bin[30]
	D4 := bin[31]

	byte1, _ := strconv.ParseInt(BinToString([]uint8{A4, A2, A1}), 2, 10)
	byte2, _ := strconv.ParseInt(BinToString([]uint8{B4, B2, B1}), 2, 10)
	byte3, _ := strconv.ParseInt(BinToString([]uint8{C4, C2, C1}), 2, 10)
	byte4, _ := strconv.ParseInt(BinToString([]uint8{D4, D2, D1}), 2, 10)

	return fmt.Sprintf("%d%d%d%d", byte1, byte2, byte3, byte4)
}

// AltCode Computes the altitude from DF4 or DF20 message, bit 20-32.
func AltCode(bin []uint8) (int, error) {
	mBit, qBit := bin[25], bin[27]

	var alt int64

	if mBit == 0 {
		if qBit == 1 {
			vbin := append(bin[19:25], bin[26])
			vbin = append(vbin, bin[28:32]...)
			//vbin := strings.Join(bin[19:25], "") + bin[26] + strings.Join(bin[28:32], "")
			alt = BinToInt(vbin)
			alt = (alt * 25) - 1000
		} else {
			C1 := bin[19]
			A1 := bin[20]
			C2 := bin[21]
			A2 := bin[22]
			C4 := bin[23]
			A4 := bin[24]
			//# _ = bin[25]
			B1 := bin[26]
			//# D1 = bin[27]     # always zero
			B2 := bin[28]
			D2 := bin[29]
			B4 := bin[30]
			D4 := bin[31]

			graystr := []uint8{D2, D4, A1, A2, A4, B1, B2, B4, C1, C2, C4}
			alt = int64(Gray2Alt(graystr))
		}
	} else {
		vbin := append(bin[19:25], bin[26])
		vbin = append(vbin, bin[26:31]...)
		//vbin := strings.Join(bin[19:25], "") + bin[26] + strings.Join(bin[26:31], "")
		alt = BinToInt(vbin)
		alt = int64(float32(alt) * 3.28084)
	}

	return int(alt), nil
}

func Gray2Alt(codestr []uint8) int {
	gc500 := codestr[:8]
	n500 := Gray2Int(gc500)

	gc100 := codestr[8:]
	n100 := Gray2Int(gc100)

	if n100 == 0 || n100 == 5 || n100 == 6 {
		return 0
	}

	if n100 == 7 {
		n100 = 5
	}

	if (n500 % 2) > 0 {
		n100 = 6 - n100
	}

	alt := (n500*500 + n100*100) - 1300

	return int(alt)
}

func OEFlag(bin []uint8) bool {
	return bin[53] == 1
}

func TypeCode(bin []uint8) uint {
	return uint(BinToInt(bin[32:37]))
}

func Data(bin []uint8) []uint8 {
	return bin[32 : len(bin)-24]
}

func Allzeros(bin []uint8) bool {
	return BinToInt(Data(bin)) == 0
}

// Check if the status bit and field bits are consistency. This Function
// is used for checking BDS code versions.
func Wrongstatus(data []uint8, sb, msb, lsb int) bool {
	status := int(data[sb-1])
	value := BinToInt(data[msb-1 : lsb])

	if status == 0 && value != 0 {
		return true
	}

	return false
}
