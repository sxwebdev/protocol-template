package model

import (
	"reflect"
	"strings"
	"time"

	"github.com/goccy/go-json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/paulmach/orb"
)

const (
	LPKey_Hardware             = "hw"
	LPKey_Firmware             = "fw"
	LPKey_IMEI                 = "imei"
	LPKey_DeviceId             = "device_id"
	LPKey_RecordId             = "record_id"
	LPKey_DeviceStatus         = "d_status"
	LPKey_DeviceTemp           = "d_temp"
	LPKey_Timestamp            = "timestamp"
	LPKey_Latitude             = "lat"
	LPKey_Longitude            = "lng"
	LPKey_Speed                = "speed"
	LPKey_Accurency            = "accurency"
	LPKey_Acceleration         = "acc"
	LPKey_Accelerometer        = "accl"
	LPKey_Altitude             = "altitude"
	LPKey_HDOP                 = "hdop"
	LPKey_PDOP                 = "pdop"
	LPKey_Direction            = "direction"
	LPKey_EcoScore             = "eco_score"
	LPKey_Accel_Course         = "course_accel"
	LPKey_Accel_Braking        = "braking_accel"
	LPKey_Accel_Turn           = "turn_accel"
	LPKey_Accel_Vertical       = "vertical_accel"
	LPKey_Sattelite_Count      = "sat_c"
	LPKey_External_Voltage     = "ext_v"
	LPKey_Battery_Voltage      = "bat_v"
	LPKey_Battery_Amperage     = "bat_a"
	LPKey_Battery_Level        = "bat_l"
	LPKey_Ignition             = "ign"
	LPKey_Motion               = "motion"
	LPKey_GSM_Signal           = "gsm_signal"
	LPKey_ICCIDPreffix         = "iccid_"
	LPKey_InputAnalogPrefix    = "in_a_"
	LPKey_InputDiscretePrefix  = "in_d_"
	LPKey_OutAnalogPrefix      = "out_a_"
	LPKey_OutDiscretePrefix    = "out_d_"
	LPKey_RS232Prefix          = "rs232_"
	LPKey_RS485Prefix          = "rs485_"
	LPKey_Digital_SensorPrefix = "d_sensor_"
	LPKey_FuelLevel            = "fuel_level"
	LPKey_FuelLevelPercent     = "fuel_level_p"
	LPKey_FuelTemp             = "fuel_temp"
	LPKey_FuelLevelPrefix      = "fuel_level_"
	LPKey_FuelTempPrefix       = "fuel_temp_"
	LPKey_Engine_RPM           = "rpm"
	LPKey_Coolant_Temp         = "coolant_temp"
	LPKey_Mileage              = "mileage"
	LPKey_Odometer             = "odometer"
	LPKey_Trip_Odometer        = "odometer_trip"
	LPKey_BT_Status            = "bt_status"
	LPKey_BLE_TempPrefix       = "ble_temp_"
	LPKey_BLE_HumidityPrefix   = "ble_hum_"
	LPKey_BLE_BatPrefix        = "ble_temp_"
	LPKey_UserDataArary        = "user_data_arr"
)

type Location struct {
	ID        uint64         `json:"id" db:"id"`
	LatLng    orb.Point      `json:"lat_lng" db:"lat_lng"`
	Speed     float64        `json:"speed" db:"speed"`
	Timestamp time.Time      `json:"timestamp" db:"timestamp"`
	Params    LocationParams `json:"d" db:"d"`
	DeviceID  uint64         `json:"device_id" db:"device_id"`
}

func NewLocation() *Location {
	return &Location{
		Params: make(LocationParams),
	}
}

func (l *Location) SetLatLng(latitude float64, longitude float64) {
	l.LatLng = orb.Point{longitude, latitude}
}

func (l *Location) SetSpeed(speed float64) {
	l.Speed = speed
}

func (l *Location) SetTimestamp(t time.Time) {
	l.Timestamp = t
}

func (l *Location) Set(key string, value interface{}) {
	key = strings.ToLower(key)
	l.Params[key] = value
}

func (l *Location) SetParamsFromJSON(data interface{}) error {
	result, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var params LocationParams
	if err := json.Unmarshal(result, &params); err != nil {
		return err
	}

	for key, value := range params {
		l.Set(key, value)
	}

	return nil
}

type LocationParams map[string]interface{}

func (s *LocationParams) Equal(params LocationParams) bool {
	return reflect.DeepEqual(s, params)
}

func (s *Location) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Timestamp, validation.Required),
		validation.Field(&s.DeviceID, validation.Required),
	)
}

type Locations []*Location

func NewLocations() Locations {
	return make(Locations, 0)
}

func (l *Locations) Add(location *Location) {
	*l = append(*l, location)
}
