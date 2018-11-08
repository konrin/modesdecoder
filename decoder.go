package adsbdecoder

import (
	"sync"
	"time"
)

type Decoder struct {
	poolMessageContext sync.Pool
}

func NewDecoder() *Decoder {
	return &Decoder{
		poolMessageContext: sync.Pool{
			New: func() interface{} {
				return NewMessageContext("")
			},
		},
	}
}

func (d *Decoder) GetMessageContext(hex string) *MessageContext {
	ctx := (d.poolMessageContext.Get()).(*MessageContext)

	ctx.time = time.Now()

	return ctx.SetHex(hex)
}

func (d *Decoder) Decode(ctx *MessageContext) (FlightData, error) {
	data := make(FlightData)

	BDSDecoders := make([]BDS, 0)

	switch ctx.GetDF() {
	case 17, 18:
		break
	case 20, 21:

		break
	}

	for i := range BDSDecoders {
		newData, err := BDSDecoders[i].Decode(ctx)
		if err != nil {
			// oops, logging
		} else {
			data = FlightDataAppend(data, newData)
		}
	}

	d.poolMessageContext.Put(ctx)

	return data, nil
}
