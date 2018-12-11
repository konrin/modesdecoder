package modesdecoder

import (
	"time"
)

type Message struct {
	Hex       string
	Bin       []uint8
	ReceiptAt time.Time

	DF   uint
	TC   uint
	OE   bool
	ICAO string
	// высота
	Altitude int
	// позывной
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
}

func NewMessage(msg string, receiptAt time.Time) *Message {
	bin := Hex2Bin(msg)
	df := uint(BinToInt(bin[0:5]))

	m := &Message{
		Hex:       msg,
		Bin:       bin,
		ReceiptAt: receiptAt,
		DF:        df,
	}

	m.ICAO = ICAO(m)

	if df == 17 || df == 18 {
		m.TC = TypeCode(bin)
		m.OE = OEFlag(bin)
	}

	return m
}

func (m *Message) GetBin() []uint8 {
	newBin := make([]uint8, len(m.Bin))
	copy(newBin, m.Bin)

	return newBin
}
