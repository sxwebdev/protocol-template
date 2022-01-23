package adm

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/sxwebdev/protocol-template/internal/model"
	"github.com/sxwebdev/protocol-template/internal/protocol/base"
	"github.com/sxwebdev/protocol-template/utils"
)

func (s *ADM) Decode(conn *base.Conn) (model.Locations, error) {

	// Max packet length 1024
	buf := make([]byte, 1024)
	read_bytes, err := bufio.NewReader(conn.Reader).Read(buf)
	if err != nil {
		return nil, err
	}
	bs := buf[:read_bytes]

	r := bytes.NewReader(bs)

	locations := model.NewLocations()

	for r.Len() > 0 {

		var deviceID uint16
		if err := binary.Read(r, binary.LittleEndian, &deviceID); err != nil {
			return nil, err
		}

		var packetLength uint8
		if err := binary.Read(r, binary.LittleEndian, &packetLength); err != nil {
			return nil, err
		}

		if packetLength == 0 {
			return nil, fmt.Errorf("empty packet length")
		}

		if packetLength == 0x84 {
			return nil, fmt.Errorf("command packet is not supported")
		}

		var packetType uint8
		if err := binary.Read(r, binary.LittleEndian, &packetType); err != nil {
			return nil, err
		}

		buf := make([]byte, packetLength-4)
		_, err := r.Read(buf)
		if err != nil {
			return nil, err
		}

		reader := bytes.NewBuffer(buf)

		if packetType == typePhoto {
			return nil, fmt.Errorf("photo packet type is not supported")
		}

		// TODO найти документацию на этот пакет
		if packetType == typeADM5 {
			return nil, fmt.Errorf("adm5 packet type is not supported")
		}

		location := model.NewLocation()

		// First packet
		if packetType == typeFirstPacket {
			var packet packetFirst
			if err := binary.Read(reader, binary.LittleEndian, &packet); err != nil {
				return nil, fmt.Errorf("parse first packet error: %v", err)
			}

			if conn.IMEI == "" {
				if err := conn.SetIMEI(packet.GetIMEI()); err != nil {
					return nil, err
				}
			}
			conn.SetHardware(strconv.Itoa(int(packet.HW)))
			conn.SetParam("reply_enabled", packet.ReplyEnabled)
			conn.SetDeviceId(strconv.Itoa(int(deviceID)))

		}

		// ADM6 packet
		if packetType&(1<<0) == 0 && packetType&(1<<1) == 0 {
			var packet packetADM6
			if err := binary.Read(reader, binary.LittleEndian, &packet); err != nil {
				return nil, fmt.Errorf("parse adm6 packet error: %v", err)
			}

			conn.SetFirmware(strconv.Itoa(int(packet.FW)))

			location.Set(model.LPKey_Acceleration, packet.Acceleration/10)
			location.Set(model.LPKey_Sattelite_Count, packet.GetSatteliteCount())
			location.Set(model.LPKey_Direction, packet.Course/10)
			location.Set(model.LPKey_HDOP, float64(packet.HDOP)/10)
			location.Set(model.LPKey_Altitude, packet.Height)

			for status_name, value := range packet.GetStatus() {
				location.Set(status_name, value)
			}

			location.SetLatLng(float64(packet.Latitude), float64(packet.Longitude))
			location.SetSpeed(float64(packet.Speed) / 10)
			location.SetTimestamp(packet.GetTime())
		}

		// ACC packet
		if utils.CheckBit8(packetType, 2) {
			var packet packetAcc
			if err := binary.Read(reader, binary.LittleEndian, &packet); err != nil {
				return nil, fmt.Errorf("parse aac packet error: %v", err)
			}

			location.Set("vib", packet.Vib)
			location.Set("vib_count", packet.VibCount)

			for i := 0; i < 8; i++ {
				location.Set(fmt.Sprintf("out_d_%d", i), utils.CheckBit8(packet.OutDiscrete, uint8(i)))
				location.Set(fmt.Sprintf("in_alarm_%d", i), utils.CheckBit8(packet.InputAnalogAlarm, uint8(i)))
			}

		}

		// Input analog packet
		if utils.CheckBit8(packetType, 3) {
			for i := 0; i < 6; i++ {
				location.Set(fmt.Sprintf("%s%d", model.LPKey_InputAnalogPrefix, i), binary.LittleEndian.Uint16(reader.Next(2)))
			}
		}

		// Input discrete packet
		if utils.CheckBit8(packetType, 4) {
			for i := 0; i < 2; i++ {
				location.Set(fmt.Sprintf("%s%d", model.LPKey_InputDiscretePrefix, i), binary.LittleEndian.Uint32(reader.Next(4)))
			}
		}

		// Fuel packet
		if utils.CheckBit8(packetType, 5) {
			for i := 0; i < 3; i++ {
				location.Set(fmt.Sprintf("%s%d", model.LPKey_FuelLevelPrefix, i), binary.LittleEndian.Uint16(reader.Next(2)))
			}
			for i := 0; i < 3; i++ {
				location.Set(fmt.Sprintf("%s%d", model.LPKey_FuelTempPrefix, i), int8(reader.Next(1)[0]))
			}
		}

		// Can data packet
		if utils.CheckBit8(packetType, 6) {
			var canDataLength uint8
			if err := binary.Read(reader, binary.LittleEndian, &canDataLength); err != nil {
				return nil, fmt.Errorf("parse can packet error: %v", err)
			}
			if canDataLength > 1 {
				canData := make([]byte, canDataLength)
				reader.Read(canData)
				return nil, fmt.Errorf("can packet is not supported, length: %v; canData: %s", canDataLength, hex.EncodeToString(canData))
			}
		}

		// Odometer packet
		if utils.CheckBit8(packetType, 7) {
			location.Set(model.LPKey_Odometer, binary.LittleEndian.Uint32(reader.Next(4)))
		}

		locations.Add(location)
	}

	return locations, nil
}
