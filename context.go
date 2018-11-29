package adsbdecoder

// const (
// 	OEFlagEven = true
// 	OEFlagOdd  = false
// )

// type MessageContext struct {
// 	hex                 string
// 	bin                 []string
// 	df                  int
// 	time                time.Time
// 	typeCode            uint
// 	LastAirEvenMsg      *MessageContext
// 	LastAirOddMsg       *MessageContext
// 	LastAirPosition     *GeoPosition
// 	LastSurfacePosition *GeoPosition
// }

// type GeoPosition struct {
// 	Lat float64
// 	Lon float64
// }

// func NewMessageContext(hex string) *MessageContext {
// 	ctx := MessageContext{
// 		hex:  hex,
// 		time: time.Now(),
// 	}

// 	return &ctx
// }

// func (m *MessageContext) SetHex(hex string) *MessageContext {
// 	m.hex = hex

// 	return m
// }

// func (m *MessageContext) GetHex() string {
// 	return m.hex
// }

// func (m *MessageContext) GetBin() []string {
// 	if len(m.bin) == 0 {
// 		m.bin = MustHex2Bin(m.hex)
// 	}

// 	return m.bin
// }

// func (m *MessageContext) GetDF() int {
// 	if m.df == 0 {
// 		m.df = DF(MustHex2Bin(m.hex))
// 	}

// 	return m.df
// }

// func (m *MessageContext) GetOEFlag() bool {
// 	return OEFlag(m.GetBin())
// }

// func (m *MessageContext) SetLastAirPositionMessage(ctx *MessageContext, flag bool) *MessageContext {
// 	if flag {
// 		m.LastAirEvenMsg = ctx
// 	} else {
// 		m.LastAirOddMsg = ctx
// 	}

// 	ctx.LastAirEvenMsg = nil
// 	ctx.LastAirOddMsg = nil

// 	return m
// }

// func (m *MessageContext) GetTime() time.Time {
// 	return m.time
// }

// func (m *MessageContext) SetTime(time time.Time) *MessageContext {
// 	m.time = time

// 	return m
// }

// func (m *MessageContext) SetLastAirPosition(geo *GeoPosition) *MessageContext {
// 	m.LastAirPosition = geo

// 	return m
// }

// func (m *MessageContext) GetTypeCode() uint {
// 	if m.typeCode == 0 {
// 		m.typeCode = TypeCode(m.GetBin())
// 	}

// 	return m.typeCode
// }

// func (m *MessageContext) Clear() {
// 	m.hex = ""
// 	m.df = 0
// 	m.LastAirEvenMsg = nil
// 	m.LastAirOddMsg = nil
// 	m.LastAirPosition = nil
// 	m.LastSurfacePosition = nil
// 	m.typeCode = 0
// }
