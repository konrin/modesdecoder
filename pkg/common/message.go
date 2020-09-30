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

	ICAO     string
	Altitude int
	Callsign string
	Category uint
	Speed    float64
	Track    float64
	Tag      string
	Rocd     int

	Squawk string

	Lat,
	Lon float64

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
	var str = "ICAO: %s\nDF: %d\n"
	var data = []interface{}{m.ICAO, m.DF}

	// var addData = func(d ...interface{}) []interface{} {

	// }

	if m.DF == 4 || m.DF == 20 {
		return fmt.Sprintf(
			str+"Altitude: %d\n",
			append(data, m.Altitude)...,
		)
	}

	return ""
}

func (m *Message) JSON() string {
	return "{}"
}
