package decoder

import "github.com/konrin/modesdecoder/pkg/common"

// Selected vertical intention
type BDS40 struct{}

func (BDS40) Is(bits *common.Bits) bool {
	if common.Allzeros(bits) {
		return false
	}

	d := common.Data(bits)

	if common.Wrongstatus(d, 1, 2, 13) {
		return false
	}

	if common.Wrongstatus(d, 4, 15, 26) {
		return false
	}

	if common.Wrongstatus(d, 27, 28, 39) {
		return false
	}

	if common.Wrongstatus(d, 48, 49, 51) {
		return false
	}

	if common.Wrongstatus(d, 54, 55, 56) {
		return false
	}

	if !d.IsZero(39, 47) {
		return false
	}

	if !d.IsZero(51, 53) {
		return false
	}

	return true
}

// Selected altitude, MCP
// https://www.skybrary.aero/index.php/Mode_Control_Panel_(MCP)
// altitude in feet
func (BDS40) AltMcp(bits *common.Bits) int {
	d := common.Data(bits)

	return int(d.Int64(1, 13) * 16) // ft
}

// Selected altitude, FMS
// https://www.skybrary.aero/index.php/Flight_Management_System
// altitude in feet
func (BDS40) AltFms(bits *common.Bits) int {
	d := common.Data(bits)

	return int(d.Int64(14, 26) * 16) // ft
}

// Barometric pressure setting
// pressure in millibar
func (BDS40) Baro(bits *common.Bits) float32 {
	d := common.Data(bits)

	if d.At(26) == 0 {
		return 0
	}

	return float32(d.Int64(27, 39))*0.1 + 800
}
