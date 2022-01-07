package model

import (
	"time"
)

// DeviceCommand ...
type DeviceCommand struct {
	ID        uint64     `json:"id" db:"id"`
	Number    uint64     `json:"number" db:"number"`
	DeviceID  uint64     `json:"device_id" db:"device_id"`
	Command   string     `json:"command" db:"command"`
	Response  *string    `json:"response" db:"response"`
	SendedAt  *time.Time `json:"sended_at" db:"sended_at"`
	RepliedAt *time.Time `json:"replied_at" db:"replied_at"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}
