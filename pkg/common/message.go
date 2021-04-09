package common

import (
	"fmt"
	"time"
)

const (
	BdsTypeAdsB  = "ADSB"
	BdsTypeCommB = "Comm-B"

	// ADSB
	//
	// Airborne position
	BdsCode0_5 = "0,5"
	// Surface position
	BdsCode0_6 = "0,6"
	// Aircraft identification and category
	BdsCode0_8 = "0,8"
	// Airborne velocity
	BdsCode0_9 = "0,9"
	// Airborne status
	BdsCode6_1 = "6,1"
	// Target state and status information
	BdsCode6_2 = "6,2"
	// Aircraft operational status
	BdsCode6_5 = "6,5"

	// Comm B
	//
	// Data link capability report
	BdsCode1_0 = "1,0"
	// Common usage GICB capability report
	BdsCode1_7 = "1,7"
	// Aircraft identification
	BdsCode2_0 = "2,0"
	// ACAS active resolution advisory
	BdsCode3_0 = "3,0"
	// Selected vertical intention
	BdsCode4_0 = "4,0"
	// Meteorological routine air report
	BdsCode4_4 = "4,4"
	// Meteorological hazard report
	BdsCode4_5 = "4,5"
	// Track and turn report
	BdsCode5_0 = "5,0"
	//Heading and speed report
	BdsCode6_0 = "6,0"
)

type Message struct {
	ReceiptAt time.Time
	BdsType   string
	BdsCode   string

	DF uint
	TC uint
	OE bool

	ICAO           string
	Alt            int
	SelectedAltMcp int
	SelectedAltFms int
	BaroSetting    float32
	Callsign       string
	Category       uint
	Speed          float64
	Track          float64
	SpeedType      string
	DirType        string
	VerticalRate   int

	Squawk string

	Lat,
	Lon float64

	OVC      int
	Capacity []string

	WindSpeed,
	WindDiraction,
	Temp,
	Hum float32
	Pressure  int
	Turblence int

	Roll float32
	Trk  float32
	GS   int
	Rtrk float32
	Tas  int

	IsAirborn bool

	hex string
	bin *Bits
}

func NewMessage(msg string, receiptAt time.Time) (*Message, error) {
	msgLen := len(msg)
	if msgLen != 14 && msgLen != 28 {
		return nil, fmt.Errorf("Incorrect message length: %d", msgLen)
	}

	bits, err := ParseHex(msg)
	if err != nil {
		return nil, err
	}

	df := bits.Uint(0, 5)

	m := &Message{
		hex:       msg,
		bin:       bits,
		ReceiptAt: receiptAt,
		DF:        df,
	}

	m.ICAO = ICAO(m)

	if df == 17 || df == 18 {
		m.TC = TypeCode(bits)
		m.OE = OEFlag(bits)
	}

	return m, nil
}

func (m *Message) GetBin() *Bits {
	return m.bin.Copy()
}

func (m *Message) GetBinRaw() []uint8 {
	return m.bin.Raw()
}

func (m *Message) GetHex() string {
	return m.hex
}

func (m *Message) String() string {
	type DataMap map[string]interface{}

	var render = func(metaData DataMap, flightData DataMap) string {
		var format = "ICAO: %s\nTime: %s\nMessage: %s"
		var data = []interface{}{
			m.ICAO,
			m.ReceiptAt.Format(time.RFC3339),
			m.GetHex(),
		}

		if len(m.BdsType) != 0 {
			format = format + "\n" + "BDS: %s(%s)"
			data = append(data, m.BdsType, m.BdsCode)
		}

		format = format + "\n"

		for key, value := range metaData {
			format += key + ":%v;"
			data = append(data, value)
		}

		format = format + "\n"

		for key, value := range flightData {
			format += "	" + key + ": %v\n"
			data = append(data, value)
		}

		return fmt.Sprintf(format, data...)
	}

	metaData := DataMap{}
	flightData := DataMap{}

	if m.DF == 4 || m.DF == 20 {
		metaData["DF"] = m.DF
		flightData["Altitude"] = m.SelectedAltMcp
	}

	if m.DF == 17 || m.DF == 18 {
		metaData["DF"] = m.DF
		metaData["TC"] = m.TC

		if m.TC >= 1 && m.TC <= 4 {
			flightData["Callsign"] = m.Callsign
			flightData["Category"] = m.Category
		}

		if m.TC >= 9 && m.TC <= 18 {
			var oe = 0
			if m.OE {
				oe = 1
			}

			metaData["OE"] = oe

			flightData["IsAirborn"] = m.IsAirborn
			flightData["Lat"] = m.Lat
			flightData["Lon"] = m.Lon
			flightData["Altitude"] = m.SelectedAltMcp
		}

		if m.TC == 19 {
			flightData["Speed"] = m.Speed
			flightData["Track"] = m.Track
			flightData["VerticalRate"] = m.VerticalRate
			flightData["SpeedType"] = m.SpeedType
		}
	}

	if m.DF == 20 || m.DF == 21 {
		if m.BdsCode == BdsCode1_0 {
			flightData["OVC"] = m.OVC
		}

		if m.BdsCode == BdsCode1_7 {
			flightData["Capacity"] = m.Capacity
		}

		if m.BdsCode == BdsCode2_0 {
			flightData["Callsign"] = m.Callsign
		}

		if m.BdsCode == BdsCode4_0 {
			flightData["SelectedAltMcp"] = m.SelectedAltMcp
			flightData["SelectedAltFms"] = m.SelectedAltFms
			flightData["Baro"] = m.BaroSetting
		}

		if m.BdsCode == BdsCode4_4 {
			flightData["WindSpeed"] = m.WindSpeed
			flightData["WindDiraction"] = m.WindDiraction
			flightData["Temp"] = m.Temp
			flightData["Hum"] = m.Hum
			flightData["Pressure"] = m.Pressure
			flightData["Turblence"] = m.Turblence
		}

		if m.BdsCode == BdsCode5_0 {
			flightData["Roll"] = m.Roll
			flightData["Trk"] = m.Trk
			flightData["GS"] = m.GS
			flightData["Rtrk"] = m.Rtrk
			flightData["Tas"] = m.Tas
		}
	}

	return render(metaData, flightData)
}

func (m *Message) JSON() string {
	return "{}"
}
