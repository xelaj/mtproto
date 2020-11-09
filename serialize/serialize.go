package serialize

import (
	"math/big"
	"reflect"

	"github.com/xelaj/go-dry"
	"github.com/xelaj/mtproto/encoding/tl"
)

const (
	int128Len = 4 * 4 // int128 16 байт
	int256Len = 4 * 8 // int256 32 байт
)

type Int128 struct {
	*big.Int
}

func (i *Int128) MarshalTL(w *tl.WriteCursor) error {
	b, err := bigIntToBytes(i.Int, 128)
	if err != nil {
		return err
	}

	return w.PutRawBytes(b)
}

func (i *Int128) UnmarshalTL(r *tl.ReadCursor) error {
	buf, err := r.PopRawBytes(int128Len)
	if err != nil {
		return err
	}

	i.Int = big.NewInt(0).SetBytes(buf)
	return nil
}

type Int256 struct {
	*big.Int
}

func (i *Int256) MarshalTL(w *tl.WriteCursor) error {
	b, err := bigIntToBytes(i.Int, 256)
	if err != nil {
		return err
	}

	return w.PutRawBytes(b)
}

func (i *Int256) UnmarshalTL(r *tl.ReadCursor) error {
	buf, err := r.PopRawBytes(int256Len)
	if err != nil {
		return err
	}

	i.Int = big.NewInt(0).SetBytes(buf)
	return nil
}

func RandomInt128() *Int128 {
	i := &Int128{big.NewInt(0)}
	i.SetBytes(dry.RandomBytes(int128Len))
	return i
}

func reflectIsInt128(v reflect.Value) bool {
	_, ok := v.Interface().(*Int128)
	return ok
}

func RandomInt256() *Int256 {
	i := &Int256{big.NewInt(0)}
	i.SetBytes(dry.RandomBytes(int256Len))
	return i
}

func reflectIsInt256(v reflect.Value) bool {
	_, ok := v.Interface().(*Int256)
	return ok
}

// ErrorSessionConfigsChanged это пустой объект, который показывает, что конфигурация сессии изменилась, и нужно создавать новую
type ErrorSessionConfigsChanged struct {
}

func (*ErrorSessionConfigsChanged) CRC() uint32 {
	panic("don't use me")
}

func (ErrorSessionConfigsChanged) Error() string {
	return "session configuration was changed"
}