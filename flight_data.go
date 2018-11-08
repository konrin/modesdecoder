package adsbdecoder

type (
	FieldData  string
	FlightData map[FieldData]interface{}
)

const (
	AIRBORN_LAT FieldData = "airborn_lat"
	AIRBORN_LON FieldData = "airborn_lon"
	ALTITUDE    FieldData = "altitude"
	CALLSING    FieldData = "callsign"
	CATEGORY    FieldData = "category"
	ICAO_FD     FieldData = "icao"
)
