package adm

import (
	"time"

	"github.com/sxwebdev/protocol-template/internal/protocol/base"
	"github.com/sxwebdev/protocol-template/utils"
)

type ADM struct {
	*base.Base
}

const (
	typeFirstPacket byte = 0x03
	typePhoto       byte = 0x0A
	typeADM5        byte = 0x01
)

type packetFirst struct {
	IMEI         [15]byte
	HW           uint8
	ReplyEnabled uint8
	Unused       [44]byte
	CRC          uint8
}

func (d *packetFirst) GetIMEI() string {
	return string(d.IMEI[:])
}

type packetADM6 struct {
	FW                uint8
	RecordNumber      uint16
	Status            uint16
	Latitude          float32
	Longitude         float32
	Course            uint16
	Speed             uint16
	Acceleration      uint8
	Height            uint16
	HDOP              uint8
	SatCount          uint8
	Timestamp         uint32
	SupplimentVoltage uint16
	BatteryVoltage    uint16
}

type packetAcc struct {
	Vib             uint8
	VibCount        uint8
	OutDiscrete     uint8
	InputAnaloglarm uint8
}

type packetInputAnalog struct {
	InputAnalog_0 uint16
	InputAnalog_1 uint16
	InputAnalog_2 uint16
	InputAnalog_3 uint16
	InputAnalog_4 uint16
	InputAnalog_5 uint16
}

type packetInputDiscrete struct {
	InputDiscrete_0 uint32
	InputDiscrete_1 uint32
}

type packetFuel struct {
	FuelLevel_0 uint16
	FuelLevel_1 uint16
	FuelLevel_2 uint16
	FuelTemp_0  int8
	FuelTemp_1  int8
	FuelTemp_2  int8
}

type packetOdometer struct {
	Odometer uint32
}

type ADMPacketData struct {
	deviceID     uint16
	packetLength uint8
	packetType   uint8
	packetFirst
	packetADM6
	packetAcc
	packetInputAnalog
	packetInputDiscrete
	packetFuel
	packetOdometer
}

func (d *packetADM6) GetTime() time.Time {
	if time.Now().UnixMilli() < int64(d.Timestamp)*1000 {
		d.Timestamp -= 2678400
	}
	return time.UnixMilli(int64(d.Timestamp) * 1000)
}

func (d *packetADM6) GetSatteliteCount() uint8 {
	return d.SatCount>>4 + d.SatCount&0x0f
}

func (d *packetADM6) GetStatus() map[string]bool {
	status := make(map[string]bool)

	status["is_device_reload"] = utils.CheckBit16(d.Status, 0)
	status["sim_slot"] = utils.CheckBit16(d.Status, 1)
	status["is_server_unavailable"] = utils.CheckBit16(d.Status, 2)
	status["is_security_mode_enabled"] = utils.CheckBit16(d.Status, 3)
	status["is_bat_low_sup"] = utils.CheckBit16(d.Status, 4)
	status["is_bad_coord"] = utils.CheckBit16(d.Status, 5)
	status["is_coord_without_move"] = utils.CheckBit16(d.Status, 6)
	status["is_sup_err"] = utils.CheckBit16(d.Status, 7)
	status["is_security_alarm"] = utils.CheckBit16(d.Status, 8)
	status["is_antenna_unavailable"] = utils.CheckBit16(d.Status, 9)
	status["is_antenna_short_circuit"] = utils.CheckBit16(d.Status, 10)
	status["is_sup_voltage_high"] = utils.CheckBit16(d.Status, 11)
	status["is_sd"] = utils.CheckBit16(d.Status, 12)
	status["is_box_open"] = utils.CheckBit16(d.Status, 13)
	status["is_coord_by_gsm"] = utils.CheckBit16(d.Status, 14)
	status["is_alarm"] = utils.CheckBit16(d.Status, 15)

	return status
}
