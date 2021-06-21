package mode

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xelaj/mtproto/internal/encoding/tl"
)

type abridged struct {
	conn io.ReadWriter
}

var _ Mode = (*abridged)(nil)

var transportModeAbridged = [...]byte{0xef} // meta:immutable

func (*abridged) getModeAnnouncement() []byte {
	return transportModeAbridged[:]
}

const (
	// If the packet length is greater than or equal to 127 words, we encode 4 bytes length, 1 is a magic
	// number, remaining 3 is real length
	// https://core.telegram.org/mtproto/mtproto-transports#abridged
	magicValueSizeMoreThanSingleByte byte = 0x7f
)

func (m *abridged) WriteMsg(msg []byte) error {
	if len(msg)%4 != 0 {
		return ErrNotMultiple{Len: len(msg)}
	}

	var size []byte

	msgLength := len(msg) / tl.WordLen
	if msgLength < int(magicValueSizeMoreThanSingleByte) {
		size = []byte{byte(msgLength)}
	} else {
		b1 := byte(msgLength)
		b2 := byte(msgLength >> 8)
		b3 := byte(msgLength >> 16)

		size = []byte{magicValueSizeMoreThanSingleByte, b1, b2, b3}
	}

	if _, err := m.conn.Write(size); err != nil {
		return err
	}
	if _, err := m.conn.Write(msg); err != nil {
		return err
	}

	return nil
}

func (m *abridged) ReadMsg() ([]byte, error) {
	sizeBuf := make([]byte, 1)
	n, err := m.conn.Read(sizeBuf)
	if err != nil {
		return nil, err
	}
	if n != 1 {
		return nil, fmt.Errorf("need to read at least 1 byte")
	}

	size := 0

	if sizeBuf[0] == magicValueSizeMoreThanSingleByte {
		sizeBuf = make([]byte, 4)
		n, err := m.conn.Read(sizeBuf[:3])
		if err != nil {
			return nil, err
		}
		if n != 3 {
			return nil, fmt.Errorf("need to read 3 bytes, got %v", n)
		}

		size = int(binary.LittleEndian.Uint32(sizeBuf))
	} else {
		size = int(sizeBuf[0])
	}

	size *= tl.WordLen

	msg := make([]byte, size)

	n, err = m.conn.Read(msg)
	if err != nil {
		return nil, err
	}
	if n != int(size) {
		return nil, fmt.Errorf("expected to read %d bytes, got %d", size, n)
	}

	return msg, nil
}
