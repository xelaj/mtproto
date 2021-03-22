package transport

import (
	"fmt"
	"github.com/xelaj/mtproto/internal/encoding/tl"
)

type ErrNotMultiple struct {
	Len int
}

func (e *ErrNotMultiple) Error() string {
	msg := "size of message not multiple of 4"
	if e.Len != 0 {
		return fmt.Sprintf(msg+" (got %v)", e.Len)
	}
	return msg
}

func checkMsgSize(msg []byte) error {
	if len(msg)%tl.WordLen != 0 {
		return &ErrNotMultiple{Len: len(msg)}
	}
	return nil
}
