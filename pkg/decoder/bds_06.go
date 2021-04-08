package decoder

import "github.com/konrin/modesdecoder/pkg/common"

// ADS-B TC=5-8
// Surface movment
type BDS06 struct{}

func (BDS06) SurfacePosition(msg *common.Message) {

}

func (BDS06) SurfacePositionWithRef(msg *common.Message) {

}

func (BDS06) SurfaceVelocity(msg *common.Message) {

}
