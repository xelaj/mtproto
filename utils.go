package mtproto

import (
	"github.com/xelaj/mtproto/serialize"
)

func MessageRequireToAck(msg serialize.TL) bool {
	switch msg.(type) {
	case /**serialize.Ping,*/ *serialize.MsgsAck:
		return false
	default:
		return true
	}
}
