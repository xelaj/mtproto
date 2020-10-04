package serialize

import (
	"math/big"
	"reflect"

	"github.com/xelaj/go-dry"
)

const (
	WordLen   = 4           // размер слова в TL (32 бита)
	LongLen   = WordLen * 2 // int64 8 байт занимает
	DoubleLen = WordLen * 2 // float64 8 байт занимает
	Int128Len = WordLen * 4 // int128 16 байт
	Int256Len = WordLen * 8 // int256 32 байт

	FuckingMagicNumber = 254  // 253 элемента максимум можно закодировать в массиве элементов
	ByteLenMagicNumber = 0xfe // ???

	// https://core.telegram.org/schema/mtproto
	crc_vector    = 0x1cb5c415
	crc_boolFalse = 0xbc799737
	crc_boolTrue  = 0x997275b5
	crc_null      = 0x56730bcc
)

var (
	int64Type = reflect.TypeOf(int64(0))
	//int128Type = reflect.TypeOf(&Int128{})
	//int256Type = reflect.TypeOf(&Int256{})
)

type TL interface {
	CRC() uint32
	TLEncoder
}

type TLDecoder interface {
	TL
	DecodeFrom(d *Decoder)
}

type Int128 struct {
	*big.Int
}

func RandomInt128() *Int128 {
	i := &Int128{big.NewInt(0)}
	i.SetBytes(dry.RandomBytes(Int128Len))
	return i
}

type Int256 struct {
	*big.Int
}

func RandomInt256() *Int256 {
	i := &Int256{big.NewInt(0)}
	i.SetBytes(dry.RandomBytes(Int256Len))
	return i
}

// TL_Null это пустой объект, который нужен для передачи в каналы TL с информацией, что ответа можно не ждать
type Null struct {
}

func (_ *Null) CRC() uint32 {
	return 0x69696969
}

func (t *Null) Encode() []byte {
	return nil
}

func (t *Null) DecodeFrom(d *Decoder) {
}

// ErrorSessionConfigsChanged это пустой объект, который показывает, что конфигурация сессии изменилась, и нужно создавать новую
type ErrorSessionConfigsChanged struct {
}

func (_ *ErrorSessionConfigsChanged) CRC() uint32 {
	panic("not acceptable")
}

func (_ *ErrorSessionConfigsChanged) Encode() []byte {
	panic("not acceptable")
}

func (_ *ErrorSessionConfigsChanged) DecodeFrom(d *Decoder) {
	panic("not acceptable")
}

func (_ *ErrorSessionConfigsChanged) Error() string {
	return "session configuration was changed"
}

type CustomObjectConstructor func(constructorID uint32) (obj TL, isEnum bool, err error)

// список функций-фабрик, которые будет использовать PopObj.
var customDecoders []CustomObjectConstructor

func AddObjectConstructor(c ...CustomObjectConstructor) {
	if customDecoders == nil {
		customDecoders = c
		return
	}

	customDecoders = append(customDecoders, c...)
}

func init() {
	AddObjectConstructor(GenerateCommonObject)
}
