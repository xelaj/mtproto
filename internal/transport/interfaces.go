package transport

import (
	"io"
)

type Conn io.ReadWriteCloser

type Mode interface {
	WriteMsg(msg []byte) error // this is not same as the io.Writer
	ReadMsg() ([]byte, error)
}
