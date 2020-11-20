// Copyright (c) 2020 KHS Films
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

func NewInt128() *Int128 {
	return &Int128{Int: big.NewInt(0).SetBytes(make([]byte, Int128Len))}
}

func RandomInt128() *Int128 {
	i := &Int128{Int: big.NewInt(0)}
	i.SetBytes(dry.RandomBytes(Int128Len))
	return i
}

// func reflectIsInt128(v reflect.Value) bool {
// 	_, ok := v.Interface().(*Int128)
// 	return ok
// }

func (i *Int128) MarshalTL(e *Encoder) error {
	e.PutRawBytes(dry.BigIntBytes(i.Int, Int128Len*bitsInByte))
	return nil
}

func (i *Int128) UnmarshalTL(d *Decoder) error {
	val := d.PopRawBytes(Int128Len)
	if d.err != nil {
		return d.err
	}
	i.Int = big.NewInt(0).SetBytes(val)
	return nil
}

type Int256 struct {
	*big.Int
}

func NewInt256() *Int256 {
	return &Int256{Int: big.NewInt(0).SetBytes(make([]byte, Int256Len))}
}

func RandomInt256() *Int256 {
	i := &Int256{big.NewInt(0)}
	i.SetBytes(dry.RandomBytes(Int256Len))
	return i
}

// func reflectIsInt256(v reflect.Value) bool {
// 	_, ok := v.Interface().(*Int256)
// 	return ok
// }

func (i *Int256) MarshalTL(e *Encoder) error {
	e.PutRawBytes(dry.BigIntBytes(i.Int, Int256Len*bitsInByte))
	return nil
}

func (i *Int256) UnmarshalTL(d *Decoder) error {
	val := d.PopRawBytes(Int256Len)
	if d.err != nil {
		return d.err
	}
	i.Int = big.NewInt(0).SetBytes(val)
	return nil
}
