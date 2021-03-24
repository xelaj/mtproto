// Copyright£ (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package tl

import "fmt"

type ErrRegisteredObjectNotFound struct {
	Crc  uint32
	Data []byte
}

func (e *ErrRegisteredObjectNotFound) Error() string {
	return fmt.Sprintf("object with provided crc not registered: 0x%08x", e.Crc)
}

type ErrMustParseSlicesExplicitly null

func (e *ErrMustParseSlicesExplicitly) Error() string {
	return "got vector CRC code when parsing unknown object: vectors can't be parsed as predicted objects"
}

type ErrorPartialWrite struct {
	Has  int
	Want int
}

func (e *ErrorPartialWrite) Error() string {
	return fmt.Sprintf("write failed: writed only %v bytes, expected %v", e.Has, e.Want)
}
