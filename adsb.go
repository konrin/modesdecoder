package adsbdecoder

type ADSB struct {
	bds05 BDS05
	bds06 BDS06
	bds08 BDS08
}

func (ADSB) Decode(ctx *MessageContext) {

}

func (a *ADSB) Callsign(ctx *MessageContext) (FlightData, error) {
	return a.bds08.Callsign(ctx)
}

func (a *ADSB) Category(ctx *MessageContext) (FlightData, error) {
	return a.bds08.Category(ctx)
}

func (a *ADSB) Icao(ctx *MessageContext) (FlightData, error) {
	data := make(FlightData)

	icao, err := ICAO(ctx)
	if err != nil {
		return data, err
	}

	data[ICAO_FD] = icao

	return data, nil
}

func (a *ADSB) Position(ctx *MessageContext) {

	if ctx.GetTypeCode() >= 5 && ctx.GetTypeCode() <= 8 {
		a.bds06.SurfacePosition(ctx)
	} else if ctx.GetTypeCode() >= 9 && ctx.GetTypeCode() <= 18 {
		// Airborne position with barometric height
		a.bds05.AirbornePosition(ctx)
	} else if ctx.GetTypeCode() >= 20 && ctx.GetTypeCode() <= 22 {
		// Airborne position with GNSS height
		a.bds05.AirbornePosition(ctx)
	}

}
