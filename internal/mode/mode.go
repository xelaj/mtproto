package mode

import (
	"bytes"
	"context"
	"io"
)

type Variant func(io.ReadWriteCloser, bool) (Mode, error)

// Mode is an interface which handles many ways as the connection sides must determine the size of the
// transmitted messages. Unlike HTTP or UDP connections, raw TCP connections, as well as WebSockets doesn't
// have a standard way to determine the size of the transmitted or received message: their main purpose is
// just to transmit bytes with right order. Mode allows the sides of the connection don't analyze traffic or
// use any end message sequence. In fact, in MTProto world, Mode works like microprotocol, which is packaging
// messages in the container that announces its size in advance
type Mode interface {
	WriteMsg([]byte) error // this is not same as the io.Writer

	// must return io.EOF if connection is open, but catched timeout
	ReadMsg(context.Context) ([]byte, error)
	Close() error
}

// Detect detects mode based on first byte sequence returned from conn
func Detect(conn io.ReadWriteCloser) (Mode, error) {
	if conn == nil {
		return nil, ErrInterfaceIsNil
	}
	b := []byte{0x0}
	_, err := conn.Read(b)
	if err != nil {
		return nil, err
	}

	var detectedMode func(io.ReadWriteCloser, bool) (Mode, error)
	switch b[0] {
	case transportModeAbridged[0]:
		detectedMode = Abridged
	case transportModeIntermediate[0]:
		modeAnnounce := make([]byte, 4)
		copy(modeAnnounce, b)
		_, err = conn.Read(modeAnnounce[1:])
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(modeAnnounce, transportModeIntermediate[:]) {
			return nil, ErrAmbiguousModeAnnounce
		}
		detectedMode = Intermediate
	default:
		return nil, ErrModeNotSupported
	}

	return detectedMode(conn, false)
}

type ReadWriteCloserCtx interface {
	Read(ctx context.Context, p []byte) (n int, err error)
	Write(ctx context.Context, p []byte) (n int, err error)
	Close(ctx context.Context) error
}
