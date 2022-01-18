package model

type Device struct {
	ID               uint64    `json:"id" db:"id"`
	IMEI             string    `json:"imei" db:"imei"`
	SerialNumber     string    `json:"serial_number" db:"serial_number"`
	ExternalDeviceID string    `json:"external_device_id" db:"external_device_id"`
	ManufacturerID   uint32    `json:"manufacturer_id" db:"manufacturer_id"`
	Hardware         string    `json:"hardware" db:"hardware"`
	Firmware         string    `json:"firmware" db:"firmware"`
	RemoteAddress    string    `json:"-" db:"remote_address"`
	LastLocations    *Location `json:"last_location" db:"last_location"`
}
