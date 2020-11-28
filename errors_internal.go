// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package mtproto

// internal errors for internal purposes

type errorSessionConfigsChanged struct{}

func (*errorSessionConfigsChanged) Error() string {
	return "session configuration was changed, need to repeat request"
}

func (*errorSessionConfigsChanged) CRC() uint32 {
	panic("makes no sense")
}
