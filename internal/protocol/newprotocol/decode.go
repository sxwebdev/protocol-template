package newprotocol

import (
	"github.com/sxwebdev/protocol-template/internal/model"
	"github.com/sxwebdev/protocol-template/internal/protocol/base"
)

func (s *NewProtocol) Decode(conn *base.Conn, bs []byte) (model.Locations, error) {
	locations := model.NewLocations()

	return locations, nil
}
