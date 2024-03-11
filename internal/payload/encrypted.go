// Copyright (c) 2020-2024 Xelaj Software
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package payload

import (
	"crypto/aes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type Encrypted struct {
	Salt      uint64
	SessionID uint64
	ID        MsgID
	SeqNo     uint32
	Msg       []byte
}

func (msg *Encrypted) NeedToAck() bool { return (msg.SeqNo & 1) != 0 }

// Envelope is a data inside encrypted part of encrypted message.
//
// compatible with MTProto 1.0 and 2.0
//
// https://core.telegram.org/mtproto/description#encrypted-message-encrypted-data
type Envelope struct {
	Salt      uint64
	SessionID uint64
	MsgID     MsgID
	SeqNo     uint32
	Msg       []byte
}

func BuildEnvelope(salt, sessionID uint64, seqNo uint32, msgID MsgID, msg []byte, padder io.Reader) (res []byte) {
	if padder != nil {
		padLen := envelopePad(eBodyOffset+len(msg), aes.BlockSize, minPadding)
		if padLen < minPadding || padLen > maxPadding {
			panic("invalid pad size")
		}
		res = make([]byte, eBodyOffset+len(msg)+padLen)
		if _, err := padder.Read(res[eBodyOffset+len(msg):]); err != nil {
			panic(err)
		} else if len(res)%aes.BlockSize != 0 {
			panic("not dividing")
		}
	} else {
		res = make([]byte, eBodyOffset+len(msg))
	}

	var le = binary.LittleEndian
	le.PutUint64(res[eSaltOffset:eSaltOffset+eSaltLen], salt)
	le.PutUint64(res[eSessionIDOffset:eSessionIDOffset+eSessionLen], sessionID)
	le.PutUint64(res[eMsgIDOffset:eMsgIDOffset+eMsgIDLen], uint64(msgID))
	le.PutUint32(res[eSeqNoOffset:eSeqNoOffset+eSeqNoLen], seqNo)
	le.PutUint32(res[eMsgLenOffset:eMsgLenOffset+eMsgLenLen], uint32(len(msg)))
	copy(res[eMsgLenOffset+eMsgLenLen:], msg)

	return res
}

const (
	eSaltOffset      = 0
	eSaltLen         = 8
	eSessionIDOffset = eSaltOffset + eSaltLen
	eSessionLen      = 8
	eMsgIDOffset     = eSessionIDOffset + eSessionLen
	eMsgIDLen        = 8
	eSeqNoOffset     = eMsgIDOffset + eMsgIDLen
	eSeqNoLen        = 4
	eMsgLenOffset    = eSeqNoOffset + eSeqNoLen
	eMsgLenLen       = 4
	eBodyOffset      = eMsgLenOffset + eMsgLenLen
	_                = 0 // msgLen

	minPadding = 12
	maxPadding = 1024
)

func pad(l, n int) int { return (n - l%n) % n }

func envelopePad(l, n, min int) int {
	p := pad(l, n)
	if p >= min {
		return p
	}

	return p + n
}

// if padder set to nil, padding won't be created
func (e *Envelope) Serialize(padder io.Reader) (res []byte) {
	if padder != nil {
		padLen := envelopePad(eBodyOffset+len(e.Msg), aes.BlockSize, minPadding)
		if padLen < minPadding || padLen > maxPadding {
			panic("invalid pad size")
		}
		res = make([]byte, eBodyOffset+len(e.Msg)+padLen)
		if _, err := padder.Read(res[eBodyOffset+len(e.Msg):]); err != nil {
			panic(err)
		} else if len(res)%aes.BlockSize != 0 {
			panic("not dividing")
		}
	} else {
		res = make([]byte, eBodyOffset+len(e.Msg))
	}

	var le = binary.LittleEndian
	le.PutUint64(res[eSaltOffset:eSaltOffset+eSaltLen], e.Salt)
	le.PutUint64(res[eSessionIDOffset:eSessionIDOffset+eSessionLen], e.SessionID)
	le.PutUint64(res[eMsgIDOffset:eMsgIDOffset+eMsgIDLen], uint64(e.MsgID))
	le.PutUint32(res[eSeqNoOffset:eSeqNoOffset+eSeqNoLen], e.SeqNo)
	le.PutUint32(res[eMsgLenOffset:eMsgLenOffset+eMsgLenLen], uint32(len(e.Msg)))
	copy(res[eMsgLenOffset+eMsgLenLen:], e.Msg)

	return res
}

func DeserializeEnvelope(b []byte) (*Envelope, error) {
	if len(b) < eBodyOffset {
		return nil, errors.New("message is too short")
	} else if len(b)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("invalid data length: have %v, want %v", len(b), pad(eBodyOffset+len(b), aes.BlockSize))
	}

	var (
		le        = binary.LittleEndian
		salt      = le.Uint64(b[eSaltOffset : eSaltOffset+eSaltLen])
		sessionID = le.Uint64(b[eSessionIDOffset : eSessionIDOffset+eSessionLen])
		msgID     = le.Uint64(b[eMsgIDOffset : eMsgIDOffset+eMsgIDLen])
		seqNo     = le.Uint32(b[eSeqNoOffset : eSeqNoOffset+eSeqNoLen])
		msgLen    = le.Uint32(b[eMsgLenOffset : eMsgLenOffset+eMsgLenLen])
		body      = b[eBodyOffset : eBodyOffset+msgLen]
	)

	// can't check that body is larger than msgLen, because it has padding
	if len(body) < int(msgLen) {
		return nil, fmt.Errorf("message not equal defined size: have %v, want %v", len(body), int(msgLen))
	}

	// Checking that padding of decrypted message is not too small or too big.
	if paddingLen := len(b) - (eBodyOffset + int(msgLen)); paddingLen < minPadding {
		return nil, fmt.Errorf("padding %d of message is too small", paddingLen)
	} else if paddingLen > maxPadding {
		return nil, fmt.Errorf("padding %d of message is too big", paddingLen)
	}

	return &Envelope{
		Salt:      salt,
		SessionID: sessionID,
		MsgID:     MsgID(msgID),
		SeqNo:     seqNo,
		Msg:       body[:msgLen],
	}, nil
}
