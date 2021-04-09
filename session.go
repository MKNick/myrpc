package myrpc

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

// 完整包校验，解决tcp粘包问题
// 格式len(4 byte)+data(len byte)

const (
	SessionHeadLen = 4
)

type Session struct {
	conn net.Conn
}

func NewSession(c net.Conn) *Session {
	return &Session{conn: c}
}

func (s *Session) write(data []byte) error {
	dataLen := len(data)
	buf := make([]byte, SessionHeadLen+dataLen)
	binary.BigEndian.PutUint32(buf[:SessionHeadLen], uint32(dataLen))
	copy(buf[SessionHeadLen:], data)
	if _, err := s.conn.Write(buf); err != nil {
		return err
	}

	return nil
}

func (s *Session) read() ([]byte, error) {
	header := make([]byte, SessionHeadLen)
	if _, err := io.ReadFull(s.conn, header); err != nil {
		return nil, err
	}

	dataLen := binary.BigEndian.Uint32(header)
	buf := make([]byte, dataLen)
	if _, err := io.ReadFull(s.conn, buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *Session) SendMsg(m RpcMessage) error {
	buf, err := Encode(m)
	if err != nil {
		return errors.New("Encode msg fail, err:" + err.Error())
	}
	err = s.write(buf)
	if err != nil {
		return errors.New("write msg fail, err:" + err.Error())
	}
	return nil
}

func (s *Session) RecvMsg() (RpcMessage, error) {
	buf, err := s.read()
	if err != nil {
		return RpcMessage{}, err
	}

	message, err := Decode(buf)
	if err != nil {
		return RpcMessage{}, err
	}

	return message, nil
}
