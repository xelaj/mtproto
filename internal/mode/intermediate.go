package mode

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func Intermediate(conn io.ReadWriteCloser, announce bool) (Mode, error) {
	if conn == nil {
		return nil, ErrInterfaceIsNil
	}

	mode := &intermediate{conn: conn}
	if announce {
		if err := mode.announce(); err != nil {
			return nil, err
		}
	}

	return mode, nil
}

type intermediate struct {
	conn io.ReadWriteCloser
}

var _ Mode = (*intermediate)(nil)

var transportModeIntermediate = [...]byte{0xee, 0xee, 0xee, 0xee} // meta:immutable

func (i *intermediate) announce() error {
	_, err := i.conn.Write(transportModeIntermediate[:])
	return err
}

func (m *intermediate) WriteMsg(msg []byte) error {
	_, err := m.conn.Write(append(
		binary.LittleEndian.AppendUint32(nil, uint32(len(msg))),
		msg...,
	))

	return err
}

func (m *intermediate) ReadMsg(ctx context.Context) (msg []byte, err error) {
	for msg == nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		sizeBuf := make([]byte, 4)
		n, err := m.conn.Read(sizeBuf)
		if errors.Is(err, context.DeadlineExceeded) {
			continue
		} else if err != nil {
			return nil, err
		}
		if n != 4 {
			return nil, fmt.Errorf("size is not length of int32, expected 4 bytes, got %d", n)
		}

		size := binary.LittleEndian.Uint32(sizeBuf)
		msg = make([]byte, int(size))

		if n, err = io.ReadFull(m.conn, msg); err != nil {
			return nil, err
		}
		if n != int(size) {
			return nil, fmt.Errorf("expected to read %d bytes, got %d", size, n)
		}
	}

	return msg, nil
}

func (m *intermediate) Close() error { return m.conn.Close() }
