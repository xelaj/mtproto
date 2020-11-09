package serialize

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xelaj/mtproto/encoding/tl"
)

type Vector []tl.Object

const crcVector = 0x1cb5c415

func (*Vector) CRC() uint32 { return crcVector }

type Bool bool

const (
	crcTrue  = 0x997275b5
	crcFalse = 0xbc799737
)

func (b *Bool) MarshalTL(w io.Writer) error {
	buf := make([]byte, 4)
	crc := crcFalse
	if *b {
		crc = crcTrue
	}

	binary.LittleEndian.PutUint32(buf, uint32(crc))
	_, err := w.Write(buf)
	return err
}

func (b *Bool) UnmarshalTL(r io.Reader) error {
	buf := make([]byte, 4)
	_, err := r.Read(buf)
	if err != nil {
		return err
	}

	crc := binary.LittleEndian.Uint32(buf)

	switch crc {
	case crcTrue:
		*b = true
		return nil
	case crcFalse:
		*b = false
		return nil
	default:
		return fmt.Errorf("not a bool value, actually: %#v", crc)
	}
}

type Long uint64

const longLen = 8 // int64 8 байт

func (l *Long) MarshalTL(w io.Writer) error {
	buf := make([]byte, longLen)
	binary.LittleEndian.PutUint64(buf, uint64(*l))
	_, err := w.Write(buf)
	return err
}

func (l *Long) UnmarshalTL(r io.Reader) error {
	buf := make([]byte, longLen)
	_, err := r.Read(buf)
	if err != nil {
		return err
	}

	*l = Long(binary.LittleEndian.Uint64(buf))
	return nil
}
