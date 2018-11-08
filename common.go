package adsbdecoder

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	hexToBinLookup = map[string]string{
		"0": "0000",
		"1": "0001",
		"2": "0010",
		"3": "0011",
		"4": "0100",
		"5": "0101",
		"6": "0110",
		"7": "0111",
		"8": "1000",
		"9": "1001",
		"a": "1010",
		"b": "1011",
		"c": "1100",
		"d": "1101",
		"e": "1110",
		"f": "1111",
	}
	crcGenerator = [25]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1}
	chars        = "#ABCDEFGHIJKLMNOPQRSTUVWXYZ#####_###############0123456789######"
)

type Common struct{}

func MustHex2Bin(hex string) []string {
	bin := ""

	for _, r := range strings.ToLower(hex) {
		c, ok := hexToBinLookup[string(r)]
		if !ok {
			return []string{}
		}

		bin += c
	}

	return strings.Split(bin, "")
}

func MustBinToInt(bin string) int64 {
	var (
		i   int64
		err error
	)

	if i, err = strconv.ParseInt(bin, 2, 64); err != nil {
		panic(err)
	}

	return i
}

func MustHexToInt(hex string) int64 {
	var (
		i   int64
		err error
	)

	if i, err = strconv.ParseInt(hex, 16, 64); err != nil {
		panic(err)
	}

	return i
}

func DF(msgbin []string) int {
	df := MustBinToInt(strings.Join(msgbin[0:5], ""))

	return int(df)
}

func CRC(ctx *MessageContext, encode bool) (string, error) {
	bin := ctx.GetBin()

	if encode {
		bin = bin[:len(bin)-24]

		for i := 0; i < 24; i++ {
			bin = append(bin, "0")
		}
	}

	for i := 0; i < len(bin)-24; i++ {
		if bin[i] != "1" {
			continue
		}

		for ci, cv := range crcGenerator {
			vi, err := strconv.Atoi(bin[i+ci])
			if err != nil {
				return "", err
			}

			bin[i+ci] = strconv.Itoa(vi ^ cv)
		}
	}

	return strings.Join(bin[len(bin)-24:], ""), nil
}

func Gray2Int(graystr string) int64 {
	num := MustBinToInt(graystr)

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

func ICAO(ctx *MessageContext) (string, error) {
	var addr string

	switch ctx.GetDF() {
	case 11, 17, 18:
		addr = ctx.GetHex()[2:8]
		break
	case 0, 4, 5, 16, 20, 21:
		coCrx, err := CRC(ctx, true)
		if err != nil {
			return "", err
		}

		c0 := MustBinToInt(coCrx)
		if err != nil {
			return "", err
		}

		c1 := MustHexToInt(ctx.GetHex()[len(ctx.GetHex())-6:])
		if err != nil {
			return "", err
		}

		addr = fmt.Sprintf("%06X", c0^c1)
		break
	}

	return addr, nil

}

func IDCODE(ctx *MessageContext) string {
	bin := ctx.GetBin()

	C1 := bin[19]
	A1 := bin[20]
	C2 := bin[21]
	A2 := bin[22]
	C4 := bin[23]
	A4 := bin[24]
	// _ = mbin[25]
	B1 := bin[26]
	D1 := bin[27]
	B2 := bin[28]
	D2 := bin[29]
	B4 := bin[30]
	D4 := bin[31]

	byte1, _ := strconv.ParseInt(A4+A2+A1, 2, 10)
	byte2, _ := strconv.ParseInt(B4+B2+B1, 2, 10)
	byte3, _ := strconv.ParseInt(C4+C2+C1, 2, 10)
	byte4, _ := strconv.ParseInt(D4+D2+D1, 2, 10)

	return fmt.Sprintf("%d%d%d%d", byte1, byte2, byte3, byte4)
}

// AltCode Computes the altitude from DF4 or DF20 message, bit 20-32.
func AltCode(ctx *MessageContext) (int, error) {
	mbin := ctx.GetBin()

	mBit, qBit := mbin[25], mbin[27]

	var alt int64

	if mBit == "0" {
		if qBit == "1" {
			vbin := strings.Join(mbin[19:25], "") + mbin[26] + strings.Join(mbin[28:32], "")
			alt = MustBinToInt(vbin)
			alt = (alt * 25) - 1000
		} else {
			C1 := mbin[19]
			A1 := mbin[20]
			C2 := mbin[21]
			A2 := mbin[22]
			C4 := mbin[23]
			A4 := mbin[24]
			//# _ = mbin[25]
			B1 := mbin[26]
			//# D1 = mbin[27]     # always zero
			B2 := mbin[28]
			D2 := mbin[29]
			B4 := mbin[30]
			D4 := mbin[31]

			graystr := D2 + D4 + A1 + A2 + A4 + B1 + B2 + B4 + C1 + C2 + C4
			alt = int64(Gray2Alt(graystr))
		}
	} else {
		vbin := strings.Join(mbin[19:25], "") + mbin[26] + strings.Join(mbin[26:31], "")
		alt = MustBinToInt(vbin)
		alt = int64(float32(alt) * 3.28084)
	}

	return int(alt), nil
}

func Gray2Alt(codestr string) int {
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

func OEFlag(bin []string) (flag bool) {
	flag = false

	if bin[53] == "1" {
		flag = true
	}

	return
}

func TypeCode(bin []string) uint {
	return uint(MustBinToInt(strings.Join(bin[32:37], "")))
}

func FlightDataAppend(map1, map2 FlightData) FlightData {
	for key := range map2 {
		map1[key] = map2[key]
	}

	return map1
}

func ShuffleFlagMessage(base *MessageContext, list ...*MessageContext) (even, odd *MessageContext, err error) {
	if base.GetOEFlag() {
		even = base

		for i := range list {
			if list[i] == nil {
				continue
			}

			if !list[i].GetOEFlag() {
				odd = list[i]

				return
			}
		}
	}

	odd = base

	for i := range list {
		if list[i] == nil {
			continue
		}

		if list[i].GetOEFlag() {
			even = list[i]

			return
		}
	}

	return nil, nil, errors.New("Not found")
}
