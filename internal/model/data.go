package model

import (
	"reflect"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/paulmach/orb"
)

// Data ...
type Data struct {
	ID         uint64     `json:"id" db:"id"`
	LatLng     orb.Point  `json:"lat_lng" db:"lat_lng"`
	Speed      float64    `json:"speed" db:"speed"`
	Timestamp  time.Time  `json:"timestamp" db:"timestamp"`
	DataParams DataParams `json:"d" db:"d"`
	DeviceID   uint64     `json:"device_id" db:"device_id"`
}

type DataParams map[string]interface{}

func (s *DataParams) Equal(params DataParams) bool {
	return reflect.DeepEqual(s, params)
}

// Validate ...
func (s *Data) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Timestamp, validation.Required),
		validation.Field(&s.DeviceID, validation.Required),
	)
}

// Datas
type Datas []*Data
