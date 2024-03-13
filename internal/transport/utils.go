package transport

import (
	"fmt"

	"github.com/xelaj/tl"

	"github.com/xelaj/mtproto/v2/internal/payload"
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

type ackMsg struct {
	MsgIDs []payload.MsgID
}

func (*ackMsg) CRC() uint32 { return 0x62d6b459 }

// msgs_ack#62d6b459 msg_ids:Vector<long> = MsgsAck;
func encodeAck(ids []payload.MsgID) []byte {
	b, err := tl.Marshal(&ackMsg{MsgIDs: ids})
	if err != nil {
		panic(err)
	}

	return b
}

type badServerSalt struct {
	BadMsgID    payload.MsgID
	BadMsgSeqNo uint32
	ErrorCode   int32
	NewSalt     uint64
}

const crcBadServerSalt = 0xedab447b

func (*badServerSalt) CRC() uint32 { return crcBadServerSalt }
