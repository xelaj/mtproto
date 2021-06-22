// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package tl

import (
	"math/big"

	"github.com/xelaj/go-dry"
)

// Int128 is alias-like type for fixed size of big int (1024 bit value). It using only for tl objects encoding
// cause native big.Int isn't supported for en(de)coding
type Int128 struct {
	*big.Int
}

// NewInt128 creates int128 with zero value
func NewInt128() *Int128 {
	return &Int128{Int: big.NewInt(0).SetBytes(make([]byte, Int128Len))}
}

// NewInt128 creates int128 with random value
func RandomInt128() *Int128 {
	i := &Int128{Int: big.NewInt(0)}
	i.SetBytes(dry.RandomBytes(Int128Len))
	return i
}

// func reflectIsInt128(v reflect.Value) bool {
// 	_, ok := v.Interface().(*Int128)
// 	return ok
// }

// MarshalTL implements tl marshaler from this package. Just don't use it by your hands, tl.Encoder does all
// what you need
func (i *Int128) MarshalTL(e *Encoder) error {
	e.PutRawBytes(dry.BigIntBytes(i.Int, Int128Len*bitsInByte))
	return nil
}

// UnmarshalTL implements tl unmarshaler from this package. Just don't use it by your hands, tl.Decoder does
// all what you need
func (i *Int128) UnmarshalTL(d *Decoder) error {
	val := d.PopRawBytes(Int128Len)
	if d.err != nil {
		return d.err
	}
	i.Int = big.NewInt(0).SetBytes(val)
	return nil
}

// Int256 is alias-like type for fixed size of big int (2048 bit value). It using only for tl objects encoding
// cause native big.Int isn't supported for en(de)coding
type Int256 struct {
	*big.Int
}

// NewInt256 creates int256 with zero value
func NewInt256() *Int256 {
	return &Int256{Int: big.NewInt(0).SetBytes(make([]byte, Int256Len))}
}

// NewInt256 creates int256 with random value
func RandomInt256() *Int256 {
	i := &Int256{big.NewInt(0)}
	i.SetBytes(dry.RandomBytes(Int256Len))
	return i
}

// func reflectIsInt256(v reflect.Value) bool {
// 	_, ok := v.Interface().(*Int256)
// 	return ok
// }

// MarshalTL implements tl marshaler from this package. Just don't use it by your hands, tl.Encoder does all
// what you need
func (i *Int256) MarshalTL(e *Encoder) error {
	e.PutRawBytes(dry.BigIntBytes(i.Int, Int256Len*bitsInByte))
	return nil
}

// UnmarshalTL implements tl unmarshaler from this package. Just don't use it by your hands, tl.Decoder does
// all what you need
func (i *Int256) UnmarshalTL(d *Decoder) error {
	val := d.PopRawBytes(Int256Len)
	if d.err != nil {
		return d.err
	}
	i.Int = big.NewInt(0).SetBytes(val)
	return nil
}
