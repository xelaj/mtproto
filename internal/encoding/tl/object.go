// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package tl

type Object interface {
	CRC() uint32
}

type FlagIndexGetter interface {
	FlagIndex() int
}

type Marshaler interface {
	MarshalTL(*Encoder) error
}

type Unmarshaler interface {
	UnmarshalTL(*Decoder) error
}

// InterfacedObject is specific struct for handling bool types, slice and null as object.
// See https://github.com/xelaj/mtproto/issues/51
type InterfacedObject struct {
	value interface{}
}

func (*InterfacedObject) CRC() uint32 {
	panic("makes no sense")
}

func (*InterfacedObject) UnmarshalTL(*Decoder) error {
	panic("impossible to (un)marshal hidden object. Use explicit methods")
}

func (*InterfacedObject) MarshalTL(*Encoder) error {
	panic("impossible to (un)marshal hidden object. Use explicit methods")
}

func (i *InterfacedObject) Unwrap() interface{} {
	return i.value
}
