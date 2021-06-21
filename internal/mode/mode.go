package mode

import (
	"bytes"
	"io"

	"github.com/pkg/errors"
)

// Mode is an interface which handles many ways as the connection sides must determine the size of the
// transmitted messages. Unlike HTTP or UDP connections, raw TCP connections, as well as WebSockets doesn't
// have a standard way to determine the size of the transmitted or received message: their main purpose is
// just to transmit bytes with right order. Mode allows the sides of the connection don't analyze traffic or
// use any end message sequence. In fact, in MTProto world, Mode works like microprotocol, which is packaging
// messages in the container that announces its size in advance
type Mode interface {
	WriteMsg([]byte) error // this is not same as the io.Writer
	ReadMsg() ([]byte, error)

	// getModeAnnouncement returns announce byte sequence to other side
	getModeAnnouncement() []byte
}

type Variant uint8

const (
	Abridged Variant = iota
	Intermediate
	PaddedIntermediate
	Full
)

func New(v Variant, conn io.ReadWriter) (Mode, error) {
	if conn == nil {
		return nil, ErrInterfaceIsNil
	}

	m, err := initMode(v, conn)
	if err != nil {
		return nil, err
	}
	announcement := m.getModeAnnouncement()
	_, err = conn.Write(announcement)
	if err != nil {
		return nil, errors.Wrap(err, "can't setup connection")
	}

	return m, nil
}

func initMode(v Variant, conn io.ReadWriter) (Mode, error) {
	switch v {
	case PaddedIntermediate, Full:
		panic("not supported yet")
	case Abridged:
		return &abridged{conn: conn}, nil
	case Intermediate:
		return &intermediate{conn: conn}, nil
	default:
		return nil, ErrModeNotSupported
	}
}

// Detect detects mode based on first byte sequence returned from conn
func Detect(conn io.ReadWriter) (Mode, error) {
	if conn == nil {
		return nil, ErrInterfaceIsNil
	}
	b := []byte{0x0}
	_, err := conn.Read(b)
	if err != nil {
		return nil, err
	}

	var detectedMode Variant
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

	return initMode(detectedMode, conn)
}

func GetVariant(m Mode) (Variant, error) {
	switch m.(type) {
	case *abridged:
		return Abridged, nil
	case *intermediate:
		return Intermediate, nil
	default:
		return Variant(0xff), errors.New("using custom mode, cant't detect")
	}
}
