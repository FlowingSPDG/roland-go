package v60hd

import (
	"fmt"
	"net"
)

type V60HD struct {
	conn net.Conn
}

const (
	stx = "\x02"
)

func NewV60HD(ipAddress string, port string) (*V60HD, error) {
	address := net.JoinHostPort(ipAddress, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &V60HD{conn: conn}, nil
}

func (v *V60HD) Close() error {
	return v.conn.Close()
}

func (v *V60HD) PGM(channel int) error {
	command := fmt.Sprintf("%sPGM:%d;", stx, channel)
	_, err := v.conn.Write([]byte(command))
	return err
}

// TODO: 受信処理

// TODO: 他コマンドの対応
