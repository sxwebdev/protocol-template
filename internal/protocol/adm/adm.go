package adm

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/paulmach/orb"
	"github.com/sxwebdev/protocol-template/internal/model"
	"github.com/sxwebdev/protocol-template/internal/protocol/base"
)

// New ...
func New(b *base.Base) base.IBase {
	return &ADM{Base: b}
}

func (s *ADM) ParseData(conn *base.Conn) error {

	// Максимальный размер пакета 1024
	buf := make([]byte, 1024)
	read_bytes, err := bufio.NewReader(conn.Conn).Read(buf)
	if err != nil {
		return err
	}

	idata, err := s.Decode(conn, buf[:read_bytes])
	if err != nil {
		return err
	}

	data, _ := idata.([]ADMPacketData)

	if len(data) == 0 {
		return nil
	}

	if conn.IMEI == "" {
		if len(data[0].GetIMEI()) == 0 {
			return errors.New("received empty imei")
		}
		conn.IMEI = data[0].GetIMEI()
	}

	if conn.DeviceID == "" && data[0].deviceID != 0 {
		conn.DeviceID = strconv.Itoa(int(data[0].deviceID))
	}

	if conn.Firmware == "" && data[0].FW != 0 {
		conn.Firmware = strconv.Itoa(int(data[0].FW))
	}

	if conn.Hardware == "" && data[0].HW != 0 {
		conn.Hardware = strconv.Itoa(int(data[0].HW))
	}

	if _, ok := conn.Params["reply_enabled"]; !ok && data[0].packetFirst.GetIMEI() != "" {
		conn.Params["reply_enabled"] = data[0].packetFirst.ReplyEnabled
	}

	if s.Config.ENV == "dev" {
		lastItem := data[len(data)-1]
		s.Logger.Debugf("last time: %s", lastItem.GetTime().String())
	}

	if conn.IMEI != "" {
		convertedData, err := s.ConvertData(data)
		if err != nil {
			return err
		}
		if err := s.AddData(conn, convertedData); err != nil {
			return err
		}
	}

	if reply_enabled, ok := conn.Params["reply_enabled"]; ok {
		reply_enabled := reply_enabled.(uint8)
		if reply_enabled == 0x02 {
			response := fmt.Sprintf("***%d*", len(data))
			if err := s.Reply(conn, []byte(response)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *ADM) ConvertData(data interface{}) ([]*model.Data, error) {
	admData, ok := data.([]ADMPacketData)
	if !ok {
		return nil, errors.New("converted data interface error")
	}
	result := make([]*model.Data, 0, len(admData))

	for _, item := range admData {
		timestamp := time.Time{}
		if item.Timestamp != 0 {
			timestamp = item.GetTime()
		}

		params := make(model.DataParams)
		params[model.DPKey_Acceleration] = item.Acceleration / 10
		params[model.DPKey_SatteliteCount] = item.GetSatteliteCount()
		params[model.DPKey_Direction] = item.Course / 10
		params[model.DPKey_HDOP] = float64(item.HDOP) / 10
		params[model.DPKey_Altitude] = item.Height

		for status_name, value := range item.GetStatus() {
			params[status_name] = value
		}

		newDataItem := &model.Data{
			LatLng:     orb.Point{float64(item.Longitude), float64(item.Latitude)},
			Speed:      float64(item.Speed) / 10,
			Timestamp:  timestamp,
			DataParams: params,
		}
		result = append(result, newDataItem)
	}

	return result, nil
}

func (s *ADM) SendCommands(conn *base.Conn, commands []interface{}) error {

	return nil
}
