package mtproto

import (
	"github.com/xelaj/mtproto/encoding/tl"
	"github.com/xelaj/mtproto/serialize"
)

func MessageRequireToAck(msg tl.Object) bool {
	switch msg.(type) {
	case /**serialize.Ping,*/ *serialize.MsgsAck:
		return false
	default:
		return true
	}
}
