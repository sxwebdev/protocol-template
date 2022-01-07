package newprotocol

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/sxwebdev/protocol-template/internal/model"
	"github.com/sxwebdev/protocol-template/internal/protocol/base"
)

// New ...
func New(b *base.Base) base.IBase {
	return &NewProtocol{Base: b}
}

// StartServer ...
func (s *NewProtocol) StartServer() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Config.Port))
	if err != nil {
		return err
	}

	s.Logger.Infof("Server %s started", s.Config.ProtocolType)

	for {
		connection, err := listener.Accept()
		if err != nil {
			s.Logger.Errorf("Error accepting connection %v", err)
			continue
		}

		s.Logger.Debugf("Accepted connection from %v", connection.RemoteAddr())

		go func(connection net.Conn) {
			if err := s.handleConnection(s.NewConn(connection)); err != nil {
				s.Logger.Errorf("Parse packet error: %v", err)
			}
		}(connection)
	}
}

// StopServer ...
func (s *NewProtocol) StopServer() error {
	for c := range s.Conns {
		s.DeleteConn(c)
	}

	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}

func (s *NewProtocol) handleConnection(conn *base.Conn) error {

	defer s.DeleteConn(conn)
	s.AddConn(conn)

	done := make(chan error, 1)
	timer := time.NewTimer(conn.IdleTimeout)
	defer timer.Stop()

	r := bufio.NewReader(conn.Conn)

	for {

		go func(r *bufio.Reader, done chan error) {

			err := func() error {
				// Edit it
				buf := make([]byte, 2)

				return s.parsePacket(conn, buf)
			}()

			done <- err

		}(r, done)

		select {
		case <-timer.C:
			return errors.New("stop by timeout")
		case err := <-done:
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(conn.IdleTimeout)
			if err != nil {
				return err
			}
		}
	}
}

// edit it
func (s *NewProtocol) parsePacket(conn *base.Conn, packet []byte) error {

	return nil
}

// edit it
func (s *NewProtocol) ConvertData(data interface{}) ([]*model.Data, error) {
	result := make([]*model.Data, 0)

	return result, nil
}

// Decode - парсит пакет и возвращает массив с данными
//
// Можно использовать для тестирования входящих пакетов не только с устройства,
// но и отдавать разработчикам результат декодирования на диагностику
//  newprotocol.Decode()
// edit it
func Decode(pkg []byte, isFullPacket bool) ([]ProtocolData, error) {
	result := make([]ProtocolData, 0)

	return result, nil
}

// Reply ...
// edit it
func (s *NewProtocol) Reply(conn *base.Conn, data []byte) error {
	_, err := conn.Conn.Write(data)
	return err
}

// SendCommand ...
// edit it
func (s *NewProtocol) SendCommands(conn *base.Conn, commands []interface{}) error {

	return nil
}
