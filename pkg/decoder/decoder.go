package decoder

import (
	"sync"
	"time"

	"github.com/konrin/modesdecoder/pkg/common"
)

const (
	CacheTTL = time.Minute
)

type Decoder struct {
	BDS05 BDS05
	BDS06 BDS06
	BDS08 BDS08
	BDS09 BDS09
	BDS10 BDS10
	BDS17 BDS17
	BDS20 BDS20
	BDS40 BDS40
	BDS44 BDS44
	BDS50 BDS50

	cacheTTL time.Duration

	mu            sync.Mutex
	cachePosition map[string]*AircraftPositionInfo
}

type AircraftPositionInfo struct {
	modd  *common.Message
	meven *common.Message
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

func (i *AircraftPositionInfo) GetEvenMsg() *common.Message {
	return i.meven
}

func (i *AircraftPositionInfo) GetOddMsg() *common.Message {
	return i.modd
}

func (i *AircraftPositionInfo) SetEvenMsg(msg *common.Message) {
	i.meven = msg
	i.updateAt = time.Now()
}

func (i *AircraftPositionInfo) SetOddMsg(msg *common.Message) {
	i.modd = msg
	i.updateAt = time.Now()
}

func (i *AircraftPositionInfo) SetPosition(lat, lon float64) {
	i.lat = lat
	i.lon = lon
	i.updateAt = time.Now()
}

func NewDecoder(cacheTTL time.Duration) *Decoder {
	return &Decoder{
		BDS05: BDS05{},
		BDS06: BDS06{},
		BDS08: BDS08{},
		BDS09: BDS09{},
		BDS10: BDS10{},
		BDS17: BDS17{},
		BDS20: BDS20{},
		BDS40: BDS40{},
		BDS44: BDS44{},
		BDS50: BDS50{},

		cacheTTL:      cacheTTL,
		cachePosition: make(map[string]*AircraftPositionInfo),
	}
}

func (d *Decoder) Decode(msg *common.Message) error {
	var err error

	if msg.DF == 17 || msg.DF == 18 {
		// Automatic Dependent Surveillance - Broadcast (ADS-B)
		return d.decodeAdsB(msg)
	}

	if msg.DF == 20 || msg.DF == 21 {
		// Mode-S Comm-B replies
		return d.decodeCommB(msg)
	}

	if msg.DF == 4 || msg.DF == 20 {
		// Altitude code
		msg.Alt, err = common.AltCode(msg.GetBin())
		if err != nil {
			return err
		}
	}

	if msg.DF == 5 || msg.DF == 21 {
		//  Identity code (squawk code)
		msg.Squawk = common.IDCODE(msg.GetBin())
	}

	return nil
}

func (d *Decoder) decodeAdsB(msg *common.Message) error {
	var err error

	posInfo := d.GetAircraftPositionFromCache(msg.ICAO)

	msg.BdsType = common.BdsTypeAdsB

	if msg.TC >= 1 && msg.TC <= 4 {
		// BDS 0,8 Aircraft identification and category
		msg.BdsCode = common.BdsCode0_8

		msg.Callsign = d.BDS08.Callsign(msg.GetBin())
		msg.Category = d.BDS08.Category(msg.GetBin())
	}

	if msg.TC >= 9 && msg.TC <= 18 {
		// BDS 0,5 Airborne position
		msg.BdsCode = common.BdsCode0_5
		msg.IsAirborn = true

		defer func() {
			if msg.OE {
				posInfo.SetEvenMsg(msg)
			} else {
				posInfo.SetOddMsg(msg)
			}
		}()

		if posInfo.HasPosition() {
			msg.Lat, msg.Lon, err = d.BDS05.AirbornePositionWithRef(
				msg.GetBinRaw(),
				msg.OE,
				posInfo.GetLat(),
				posInfo.GetLon(),
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

			if msg.OE && posInfo.HasOddMsg() {
				bin1 = msg.GetBinRaw()
				bin2 = posInfo.GetOddMsg().GetBinRaw()
				at1 = msg.ReceiptAt
				at2 = posInfo.GetOddMsg().ReceiptAt
			} else if !msg.OE && posInfo.HasEvenMsg() {
				bin2 = msg.GetBinRaw()
				bin1 = posInfo.GetEvenMsg().GetBinRaw()
				at2 = msg.ReceiptAt
				at1 = posInfo.GetEvenMsg().ReceiptAt
			} else {
				return nil
			}

			msg.Lat, msg.Lon, err = d.BDS05.AirbornePosition(bin1, at1, bin2, at2)
			if err != nil {
				return err
			}
		}

		posInfo.SetPosition(msg.Lat, msg.Lon)

		msg.Alt, err = d.BDS05.Altitude(msg.GetBinRaw(), msg.TC)
		if err != nil {
			return err
		}
	}

	if msg.TC >= 5 && msg.TC <= 8 {
		// BDS 0,6 Surface position
		msg.BdsCode = common.BdsCode0_8
		msg.IsAirborn = false

		// TODO
	}

	if msg.TC == 19 {
		// BDS 0,9 Airborne velocity
		msg.BdsCode = common.BdsCode0_9

		msg.Speed,
			msg.Track,
			msg.VerticalRate,
			msg.SpeedType,
			msg.DirType,
			err = d.BDS09.AirborneVelocity(msg.GetBin())
		if err != nil {
			return err
		}
	}

	if msg.TC == 28 {
		// BDS 6,1: Airborne status
		msg.BdsCode = common.BdsCode6_1

		// TODO
	}

	if msg.TC == 29 {
		// BDS 6,2: Target state and status information
		msg.BdsCode = common.BdsCode6_2

		// TODO
	}

	if msg.TC == 31 {
		// BDS 6,5: Aircraft operational status
		msg.BdsCode = common.BdsCode6_5

		// TODO
	}

	return nil
}

func (d *Decoder) decodeCommB(msg *common.Message) error {
	msg.BdsType = common.BdsTypeCommB

	if d.BDS10.Is(msg.GetBin()) {
		msg.BdsCode = common.BdsCode1_0

		msg.OVC = d.BDS10.OVC(msg.GetBin())
	}

	if d.BDS17.Is(msg.GetBin()) {
		msg.BdsCode = common.BdsCode1_7

		msg.Capacity = d.BDS17.Cap(msg.GetBin())
	}

	if d.BDS20.Is(msg.GetBin()) {
		msg.BdsCode = common.BdsCode2_0

		msg.Callsign = d.BDS20.CS(msg.GetBin())
	}

	if d.BDS40.Is(msg.GetBin()) {
		msg.BdsCode = common.BdsCode4_0
		msg.SelectedAltMcp = d.BDS40.AltMcp(msg.GetBin())
		msg.SelectedAltFms = d.BDS40.AltFms(msg.GetBin())
		msg.BaroSetting = d.BDS40.Baro(msg.GetBin())
	}

	if d.BDS44.Is(msg.GetBin()) {
		msg.BdsCode = common.BdsCode4_4

		msg.WindSpeed, msg.WindDiraction = d.BDS44.Wind(msg.GetBin())
		msg.Temp = d.BDS44.Temp(msg.GetBin())
		msg.Hum = d.BDS44.Hum(msg.GetBin())
		msg.Pressure = d.BDS44.Pressure(msg.GetBin())
	}

	if d.BDS50.Is(msg.GetBin()) {
		msg.BdsCode = common.BdsCode5_0
		msg.Roll = d.BDS50.Roll(msg.GetBin())
		msg.Trk = d.BDS50.TRK(msg.GetBin())
		msg.GS = d.BDS50.GS(msg.GetBin())
		msg.Rtrk = d.BDS50.RTRK(msg.GetBin())
		msg.Tas = d.BDS50.TAS(msg.GetBin())
	}

	return nil
}

func (d *Decoder) GetAircraftPositionFromCache(icao string) *AircraftPositionInfo {
	d.mu.Lock()
	defer d.mu.Unlock()

	info, ok := d.cachePosition[icao]
	if !ok {
		info = &AircraftPositionInfo{
			updateAt: time.Now(),
		}
	}

	if time.Until(info.updateAt) > d.cacheTTL {
		info = &AircraftPositionInfo{
			updateAt: time.Now(),
		}
	}

	d.cachePosition[icao] = info

	return info
}
