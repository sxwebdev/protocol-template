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
	LPKey_Speed               = "speed"
	LPKey_Accurency           = "accurency"
	LPKey_Acceleration        = "acc"
	LPKey_Altitude            = "altitude"
	LPKey_HDOP                = "hdop"
	LPKey_PDOP                = "pdop"
	LPKey_Direction           = "direction"
	LPKey_Sattelite_Count     = "sat_c"
	LPKey_Supplement_Voltage  = "sup_v"
	LPKey_Battery_Voltage     = "bat_v"
	LPKey_Battery_Amperage    = "bat_a"
	LPKey_Ignition            = "ign"
	LPKey_Motion              = "motion"
	LPKey_GSM_Signal          = "gsm_signal"
	LPKey_InputAnalogPrefix   = "in_a_"
	LPKey_InputDiscretePrefix = "in_d_"
	LPKey_OutAnalogPrefix     = "out_a_"
	LPKey_OutDiscretePrefix   = "out_d_"
	LPKey_FuelLevel           = "fuel_level"
	LPKey_FuelTemp            = "fuel_temp"
	LPKey_FuelLevelPrefix     = "fuel_level_"
	LPKey_FuelTempPrefix      = "fuel_temp_"
	LPKey_Odometer            = "odometer"
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
