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

	// Блядские магические числа
	FuckingMagicNumber = 254  // 253 элемента максимум можно закодировать в массиве элементов
	ByteLenMagicNumber = 0xfe // ???

	// https://core.telegram.org/schema/mtproto
	crcVector = 0x1cb5c415
	crcFalse  = 0xbc799737
	crcTrue   = 0x997275b5
	crcNull   = 0x56730bcc

	// CrcRpcResult публичная переменная, т.к. это специфический конструктор
	CrcRpcResult  = 0xf35c6d01 // nolint
	CrcGzipPacked = 0x3072cfa1
)

var (
	int64Type = reflect.TypeOf(int64(0))
	// int128Type = reflect.TypeOf(&Int128{})
	// int256Type = reflect.TypeOf(&Int256{})
)

type TL interface {
	CRC() uint32
	TLEncoder
}

func reflectIsTL(v reflect.Value) bool {
	_, ok := v.Interface().(TL)
	return ok
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

func reflectIsInt128(v reflect.Value) bool {
	_, ok := v.Interface().(*Int128)
	return ok
}

type Int256 struct {
	*big.Int
}

func RandomInt256() *Int256 {
	i := &Int256{big.NewInt(0)}
	i.SetBytes(dry.RandomBytes(Int256Len))
	return i
}

func reflectIsInt256(v reflect.Value) bool {
	_, ok := v.Interface().(*Int256)
	return ok
}

// TL_Null это пустой объект, который нужен для передачи в каналы TL с информацией, что ответа можно не ждать
type Null struct {
}

func (*Null) CRC() uint32 {
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

func (*ErrorSessionConfigsChanged) CRC() uint32 {
	panic("not acceptable")
}

func (*ErrorSessionConfigsChanged) Encode() []byte {
	panic("not acceptable")
}

func (*ErrorSessionConfigsChanged) DecodeFrom(d *Decoder) {
	panic("not acceptable")
}

func (*ErrorSessionConfigsChanged) Error() string {
	return "session configuration was changed"
}

// --------------------------------------------------------------------------------------

// dummy bool struct for methods generation
type Bool struct{}

func (*Bool) CRC() uint32 {
	panic("it's a dummy constructor!")
}

func (t *Bool) Encode() []byte {
	panic("it's a dummy constructor!")
}

func (t *Bool) DecodeFrom(d *Decoder) {
	panic("it's a dummy constructor!")
}

// dummy bool struct for methods generation
type Long struct{}

func (*Long) CRC() uint32 {
	panic("it's a dummy constructor!")
}

func (*Long) Encode() []byte {
	panic("it's a dummy constructor!")
}

func (*Long) DecodeFrom(d *Decoder) {
	panic("it's a dummy constructor!")
}

// dummy bool struct for methods generation
type Int struct{}

func (*Int) CRC() uint32 {
	panic("it's a dummy constructor!")
}

func (*Int) Encode() []byte {
	panic("it's a dummy constructor!")
}

func (*Int) DecodeFrom(d *Decoder) {
	panic("it's a dummy constructor!")
}

// блять! вектор ведь это тоже структура! короче вот эта структура просто в себе хранит
// слайс либо стандартных типов ([]int32, []float64, []bool и прочее), либо тл объекта
// ([]TL). алгоритм который использует эту структуру должен гарантировать, что параметр
// является слайсом, а элементы слайса являются либо стандартные типы, либо TL.
type InnerVectorObject struct {
	I interface{}
}

func (*InnerVectorObject) CRC() uint32 {
	panic("it's a dummy constructor!")
}

func (*InnerVectorObject) Encode() []byte {
	panic("it's a dummy constructor!")
}

// --------------------------------------------------------------------------------------

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
