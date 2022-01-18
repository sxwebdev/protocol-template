package adm

import (
	"fmt"

	"github.com/sxwebdev/protocol-template/internal/protocol/base"
)

// New ...
func New(b *base.Base) base.IBase {
	return &ADM{Base: b}
}

func (s *ADM) ParseData(conn *base.Conn) error {

	locations, err := s.Decode(conn, nil)
	if err != nil {
		return err
	}

	s.PrintLastLocTime(locations)

	if conn.IMEI != "" && len(locations) > 0 {
		if err := s.AddData(conn, locations); err != nil {
			return err
		}
	}

	if reply_enabled, ok := conn.Params["reply_enabled"]; ok {
		reply_enabled := reply_enabled.(uint8)
		if reply_enabled == 0x02 {
			response := fmt.Sprintf("***%d*", len(locations))
			if err := s.Reply(conn, []byte(response)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *ADM) SendCommands(conn *base.Conn, commands []interface{}) error {

	return nil
}
