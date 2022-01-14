package adm

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/sxwebdev/protocol-template/internal/protocol/base"
	"github.com/sxwebdev/protocol-template/utils"
)

func (s *ADM) Decode(conn *base.Conn, r *bytes.Reader) (interface{}, map[string]interface{}, error) {

	packets := make([]ADMPacketData, 0)

	for r.Len() > 0 {

		var deviceID uint16
		if err := binary.Read(r, binary.LittleEndian, &deviceID); err != nil {
			return nil, nil, err
		}

		var packetLength uint8
		if err := binary.Read(r, binary.LittleEndian, &packetLength); err != nil {
			return nil, nil, err
		}

		if packetLength == 0 {
			return nil, nil, fmt.Errorf("empty packet length")
		}

		if packetLength == 0x84 {
			return nil, nil, fmt.Errorf("command packet is not supported")
		}

		var packetType uint8
		if err := binary.Read(r, binary.LittleEndian, &packetType); err != nil {
			return nil, nil, err
		}
		packetData := ADMPacketData{
			deviceID:     deviceID,
			packetLength: packetLength,
			packetType:   packetType,
		}

		buf := make([]byte, packetLength-4)
		_, err := r.Read(buf)
		if err != nil {
			return nil, nil, err
		}

		reader := bytes.NewReader(buf)

		if packetType == typePhoto {
			return nil, nil, fmt.Errorf("photo packet type is not supported")
		}

		// TODO найти документацию на этот пакет
		if packetType == typeADM5 {
			return nil, nil, fmt.Errorf("adm5 packet type is not supported")
		}

		// First packet
		if packetType == typeFirstPacket {
			if err := binary.Read(reader, binary.LittleEndian, &packetData.packetFirst); err != nil {
				return nil, nil, fmt.Errorf("parse iemi packet error: %v", err)
			}
		}

		// ADM6 packet
		if packetType&(1<<0) == 0 && packetType&(1<<1) == 0 {
			if err := binary.Read(reader, binary.LittleEndian, &packetData.packetADM6); err != nil {
				return nil, nil, fmt.Errorf("parse adm6 packet error: %v", err)
			}
		}

		// ACC packet
		if utils.CheckBit8(packetType, 2) {
			if err := binary.Read(reader, binary.LittleEndian, &packetData.packetAcc); err != nil {
				return nil, nil, fmt.Errorf("parse aac packet error: %v", err)
			}
		}

		// Input analog packet
		if utils.CheckBit8(packetType, 3) {
			if err := binary.Read(reader, binary.LittleEndian, &packetData.packetInputAnalog); err != nil {
				return nil, nil, fmt.Errorf("parse input analog packet error: %v", err)
			}
		}

		// Input discrete packet
		if utils.CheckBit8(packetType, 4) {
			if err := binary.Read(reader, binary.LittleEndian, &packetData.packetInputDiscrete); err != nil {
				return nil, nil, fmt.Errorf("parse input discrete packet error: %v", err)
			}
		}

		// Fuel packet
		if utils.CheckBit8(packetType, 5) {
			if err := binary.Read(reader, binary.LittleEndian, &packetData.packetFuel); err != nil {
				return nil, nil, fmt.Errorf("parse fuel packet error: %v", err)
			}
		}

		// Can data packet
		if utils.CheckBit8(packetType, 6) {
			var canDataLength uint8
			if err := binary.Read(reader, binary.LittleEndian, &canDataLength); err != nil {
				return nil, nil, fmt.Errorf("parse can packet error: %v", err)
			}
			if canDataLength > 1 {
				canData := make([]byte, canDataLength)
				reader.Read(canData)
				return nil, nil, fmt.Errorf("can packet is not supported, length: %v; canData: %s", canDataLength, hex.EncodeToString(canData))
			}
		}

		// Odometer packet
		if utils.CheckBit8(packetType, 7) {
			if err := binary.Read(reader, binary.LittleEndian, &packetData.Odometer); err != nil {
				return nil, nil, fmt.Errorf("parse odometer packet error: %v", err)
			}
		}

		packets = append(packets, packetData)
	}

	return packets, nil, nil
}
