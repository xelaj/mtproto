package transport

import (
	"io"
)

type ConnConfig interface{}
type Conn io.ReadWriteCloser

type ModeConfig = func(Conn) (Mode, error)
type Mode interface {
	WriteMsg(msg []byte) error // this is not same as the io.Writer
	ReadMsg() ([]byte, error)
}
