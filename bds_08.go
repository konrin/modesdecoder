package adsbdecoder

import "strings"

// // BDS08 BDS 0,8
// // ADS-B TC=1-4
// // Aircraft identitification and category
type BDS08 struct{}

func (BDS08) Callsign(msg *Message) string {
	csbin := msg.Bin[40:96]

	ct := [][]uint8{
		csbin[0:6], csbin[6:12], csbin[12:18], csbin[18:24],
		csbin[24:30], csbin[30:36], csbin[36:42], csbin[42:48],
	}

	var cs string

	for i := range ct {
		cs += string(chars[BinToInt(ct[i])])
	}

	return strings.Replace(cs, "#", "", -1)
}

func (BDS08) Category(msg *Message) uint {
	return uint(BinToInt(msg.Bin[5:8]))
}

// func (BDS08) Callsign(ctx *MessageContext) (FlightData, error) {
// 	data := make(FlightData)

// 	if ctx.GetTypeCode() < 1 || ctx.GetTypeCode() > 4 {
// 		return data, errors.New("Not a identification message")
// 	}

// 	csbin := strings.Join(ctx.GetBin()[40:96], "")

// 	var cs string

// 	ct := [8]string{
// 		csbin[0:6], csbin[6:12], csbin[12:18], csbin[18:24],
// 		csbin[24:30], csbin[30:36], csbin[36:42], csbin[42:48],
// 	}

// 	for i := range ct {
// 		cs += string(chars[MustBinToInt(ct[i])])
// 	}

// 	data[CALLSING] = strings.Replace(cs, "#", "", -1)

// 	return data, nil
// }

// func (BDS08) Category(ctx *MessageContext) (FlightData, error) {
// 	data := make(FlightData)

// 	data[CATEGORY] = int(MustBinToInt(strings.Join(ctx.GetBin()[5:8], "")))

// 	return data, nil
// }
