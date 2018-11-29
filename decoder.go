package adsbdecoder

import (
	"time"
)

type Decoder struct {
	BDS06 BDS06
	BDS08 BDS08
}

func NewDecoer() *Decoder {
	return &Decoder{
		BDS06: BDS06{},
		BDS08: BDS08{},
	}
}

func (d *Decoder) Decode(msg, lastOdd, lastEven string, latRef, lonRef float64) {
	m := NewMessage(msg, time.Now())

	if m.DF == 11 {

	} else if m.DF == 17 || m.DF == 18 {
		// Automatic Dependent Surveillance - Broadcast (ADS-B)

		if m.TC >= 1 && m.TC <= 4 {
			// BDS 0,8: Aircraft identification and category

			m.Callsign = d.BDS08.Callsign(m)
			m.Category = d.BDS08.Category(m)
		} else if m.TC >= 5 && m.TC <= 8 {
			// BDS 0,6: Surface position

		} else if m.TC >= 9 && m.TC <= 18 {
			// BDS 0,5: Airborne position

		} else if m.TC == 19 {
			// BDS 0,9: Airborne velocity
		} else if m.TC == 28 {
			// BDS 6,1: Airborne status [to be implemented]
		} else if m.TC == 29 {
			// BDS 6,2: Target state and status information [to be implemented]
		} else if m.TC == 31 {
			// BDS 6,5: Aircraft operational status [to be implemented]
		}

	} else if m.DF == 20 || m.DF == 21 {
		// Mode-S Comm-B replies
	}

	if m.DF == 4 || m.DF == 20 {
		// Altitude code
	}

	if m.DF == 5 || m.DF == 21 {
		//  Identity code (squawk code)
	}
}
