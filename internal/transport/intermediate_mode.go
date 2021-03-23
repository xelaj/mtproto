package transport

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/xelaj/mtproto/internal/encoding/tl"
)

type intermediateMode struct {
	conn io.ReadWriter
}

var transportModeIntermediate = [...]byte{0xee, 0xee, 0xee, 0xee} // meta:immutable

func NewIntermediateMode(conn io.ReadWriter) (Mode, error) {
	if conn == nil {
		return nil, errors.New("conn is nil")
	}

	_, err := conn.Write(transportModeIntermediate[:])
	if err != nil {
		return nil, errors.Wrap(err, "can't setup connection")
	}

	return &intermediateMode{conn: conn}, nil
}

func (m *intermediateMode) WriteMsg(msg []byte) error {
	if err := checkMsgSize(msg); err != nil {
		return err
	}
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(len(msg)))
	if _, err := m.conn.Write(size); err != nil {
		return err
	}
	if _, err := m.conn.Write(msg); err != nil {
		return err
	}

	return nil
}

func (m *intermediateMode) ReadMsg() ([]byte, error) {
	sizeInBytes := make([]byte, tl.WordLen)
	n, err := m.conn.Read(sizeInBytes)
	if err != nil {
		return nil, err
	}
	if n != tl.WordLen {
		return nil, fmt.Errorf("size is not length of int32, expected 4 bytes, got %d", n)
	}

	size := binary.LittleEndian.Uint32(sizeInBytes)
	msg := make([]byte, int(size))
	n, err = m.conn.Read(msg)
	if err != nil {
		return nil, err
	}
	if n != int(size) {
		return nil, fmt.Errorf("expected to read %d bytes, got %d", size, n)
	}

	return msg, nil
}
