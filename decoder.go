package modesdecoder

import (
	"sync"
	"time"
)

const (
	// в сек
	CACHE_TTL float64 = 60
)

type Decoder struct {
	BDS05 BDS05
	BDS06 BDS06
	BDS08 BDS08
	BDS09 BDS09

	cacheTTL float64

	sync.Mutex
	cachePosition map[string]*AircraftPositionInfo
}

type AircraftPositionInfo struct {
	modd  *Message
	meven *Message
	lat,
	lon float64
	updateAt time.Time
}

func (i *AircraftPositionInfo) HasPosition() bool {
	return (i.lat + i.lon) > 0
}

func (i *AircraftPositionInfo) HasEvenMsg() bool {
	return i.meven != nil
}

func (i *AircraftPositionInfo) HasOddMsg() bool {
	return i.modd != nil
}

func (i *AircraftPositionInfo) GetLat() float64 {
	return i.lat
}

func (i *AircraftPositionInfo) GetLon() float64 {
	return i.lon
}

func (i *AircraftPositionInfo) GetEvenMsg() *Message {
	return i.meven
}

func (i *AircraftPositionInfo) GetOddMsg() *Message {
	return i.modd
}

func (i *AircraftPositionInfo) SetEvenMsg(msg *Message) {
	i.meven = msg
}

func (i *AircraftPositionInfo) SetOddMsg(msg *Message) {
	i.modd = msg
}

func NewDecoder(cacheTTL float64) *Decoder {
	return &Decoder{
		BDS05: BDS05{},
		BDS06: BDS06{},
		BDS08: BDS08{},
		BDS09: BDS09{},

		cacheTTL:      cacheTTL,
		cachePosition: make(map[string]*AircraftPositionInfo, 0),
	}
}

func (d *Decoder) Decode(msg *Message) error {
	var err error

	posInfo := d.GetAircraftPositionFromCache(msg.ICAO)

	if msg.DF == 11 {

	} else if msg.DF == 17 || msg.DF == 18 {
		// Automatic Dependent Surveillance - Broadcast (ADS-B)

		if msg.TC >= 1 && msg.TC <= 4 {
			// BDS 0,8: Aircraft identification and category

			msg.Callsign = d.BDS08.Callsign(msg.GetBin())
			msg.Category = d.BDS08.Category(msg.GetBin())
		} else if msg.TC >= 5 && msg.TC <= 8 {
			// BDS 0,6: Surface position

		} else if msg.TC >= 9 && msg.TC <= 18 {
			// BDS 0,5: Airborne position

			err = d.decodeAirbonPosition(msg, *posInfo)
			if err != nil {
				return err
			}

			if msg.OE {
				posInfo.SetEvenMsg(msg)
			} else {
				posInfo.SetOddMsg(msg)
			}

			msg.Altitude, err = d.BDS05.Altitude(msg.GetBin(), msg.TC)
			if err != nil {
				return err
			}

			d.AircraftPositionMarkUpdated(posInfo)
		} else if msg.TC == 19 {
			// BDS 0,9: Airborne velocity
			msg.Speed, msg.Track, msg.Rocd, msg.Tag, err = d.BDS09.AirborneVelocity(msg.GetBin())
			if err != nil {
				return err
			}
		} else if msg.TC == 28 {
			// BDS 6,1: Airborne status
		} else if msg.TC == 29 {
			// BDS 6,2: Target state and status information
		} else if msg.TC == 31 {
			// BDS 6,5: Aircraft operational status
		}

	} else if msg.DF == 20 || msg.DF == 21 {
		// Mode-S Comm-B replies
	}

	if msg.DF == 4 || msg.DF == 20 {
		// Altitude code
		msg.Altitude, err = AltCode(msg.GetBin())
		if err != nil {
			return err
		}
	}

	if msg.DF == 5 || msg.DF == 21 {
		//  Identity code (squawk code)

		msg.Squawk = IDCODE(msg.GetBin())
	}

	return nil
}

func (d *Decoder) decodeAirbonPosition(msg *Message, info AircraftPositionInfo) error {
	var err error

	if info.HasPosition() {
		msg.Lat, msg.Lon, err = d.BDS05.AirbornePositionWithRef(
			msg.GetBin(),
			msg.OE,
			info.GetLat(),
			info.GetLon(),
		)
		if err != nil {
			return err
		}
	} else {
		var (
			bin1,
			bin2 []uint8

			at1,
			at2 time.Time
		)

		if msg.OE && info.HasOddMsg() {
			bin1 = msg.GetBin()
			bin2 = info.GetOddMsg().GetBin()
			at1 = msg.ReceiptAt
			at2 = info.GetOddMsg().ReceiptAt
		} else if !msg.OE && info.HasEvenMsg() {
			bin2 = msg.GetBin()
			bin1 = info.GetEvenMsg().GetBin()
			at2 = msg.ReceiptAt
			at1 = info.GetEvenMsg().ReceiptAt
		} else {
			return nil
		}

		msg.Lat, msg.Lon, err = d.BDS05.AirbornePosition(bin1, at1, bin2, at2)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) GetAircraftPositionFromCache(icao string) *AircraftPositionInfo {
	d.Lock()
	defer d.Unlock()

	info, ok := d.cachePosition[icao]
	if !ok {
		info = &AircraftPositionInfo{
			updateAt: time.Now(),
		}
	}

	if info.updateAt.Sub(time.Now()).Seconds() > d.cacheTTL {
		info = &AircraftPositionInfo{
			updateAt: time.Now(),
		}
	}

	d.cachePosition[icao] = info

	return info
}

func (d *Decoder) AircraftPositionMarkUpdated(info *AircraftPositionInfo) {
	info.updateAt = time.Now()
}
