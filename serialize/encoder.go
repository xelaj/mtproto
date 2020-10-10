package serialize

import (
	"encoding/binary"
	"math"
	"math/big"
	"reflect"

	"github.com/xelaj/go-dry"
)

type TLEncoder interface {
	Encode() []byte
}

type Encoder struct {
	buf []byte
}

func NewEncoder() *Encoder {
	return &Encoder{make([]byte, 0, 512)} // 512 это капасити, сделано на всякий случай, что бы не тормозить кодировку выделением памяти
}

func (e *Encoder) Result() []byte {
	return e.buf
}

// PutBool очень специфичный тип, т.к. есть отдельный конструктор под true и false,
// то можно считать, что это две crc константы
func (e *Encoder) PutBool(v bool) {
	buf := make([]byte, WordLen)
	crc := crcFalse
	if v {
		crc = crcTrue
	}

	binary.LittleEndian.PutUint32(buf, uint32(crc))

	e.buf = append(e.buf, buf...)
}

func (e *Encoder) PutInt(v int32) {
	e.PutUint(uint32(v))
}

func (e *Encoder) PutUint(v uint32) {
	buf := make([]byte, WordLen)
	binary.LittleEndian.PutUint32(buf, v)
	e.buf = append(e.buf, buf...)
}

func (e *Encoder) PutCRC(v uint32) {
	e.PutUint(v) // я так и не понял, кажется что crc это bigendian, но видимо нет
}

func (e *Encoder) PutLong(v int64) {
	buf := make([]byte, WordLen*2)
	binary.LittleEndian.PutUint64(buf, uint64(v))
	e.buf = append(e.buf, buf...)
}

func (e *Encoder) PutDouble(v float64) {
	buf := make([]byte, WordLen*2)
	binary.LittleEndian.PutUint64(buf, math.Float64bits(v))
	e.buf = append(e.buf, buf...)
}

func (e *Encoder) PutBigInt(s *big.Int) {
	e.PutRawBytes(s.Bytes())
}

func (e *Encoder) PutInt128(s *Int128) {
	// RawBytes потому что... ???
	e.PutRawBytes(dry.BigIntBytes(s.Int, 128))
}

func (e *Encoder) PutInt256(s *Int256) {
	// RawBytes потому что... ???
	e.PutRawBytes(dry.BigIntBytes(s.Int, 256))
}

func (e *Encoder) PutString(msg string) {
	e.PutMessage([]byte(msg))
}

func (e *Encoder) PutMessage(msg []byte) {
	if len(msg) < FuckingMagicNumber {
		e.putTinyBytes(msg)
	} else {
		e.putLargeBytes(msg)
	}
}

func (e *Encoder) putTinyBytes(msg []byte) {
	if len(msg) >= FuckingMagicNumber {
		panic("tiny messages supports maximum 253 elements")
	}

	// здесь мы считаем, что длинна итогового сообщения должна делиться на 4
	// (32/8 = 4, 4 байта одно слово)
	// поэтому мы создаем buf с размером, достаточным для пихания массива + 0-3 доп байта что бы итог делился на 4
	realBytesLen := (1 + len(msg))
	factBytesLen := realBytesLen
	if factBytesLen%WordLen > 0 {
		factBytesLen += WordLen - factBytesLen%WordLen
	}

	buf := make([]byte, factBytesLen)
	buf[0] = byte(len(msg)) // пихаем в первый байт размер сообщения
	copy(buf[1:], msg)

	e.buf = append(e.buf, buf...)
}

func (e *Encoder) putLargeBytes(msg []byte) {
	if len(msg) < FuckingMagicNumber {
		panic("режим работы в маленьких сообщениях не гарантирован")
	}

	maxLen := 1 << 24 // 3 байта 24 бита, самый первый это 0xfe оставшиеся 3 как раз длина
	if len(msg) > maxLen {
		panic("message entity too large")
	}

	realBytesLen := (WordLen + len(msg)) // первым идет магический байт и 3 байта длины
	factBytesLen := realBytesLen
	if factBytesLen%WordLen > 0 {
		factBytesLen += WordLen - factBytesLen%WordLen
	}

	littleEndianLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(littleEndianLength, uint32(len(msg)))

	buf := make([]byte, factBytesLen)
	buf[0] = byte(ByteLenMagicNumber)
	buf[1] = littleEndianLength[0]
	buf[2] = littleEndianLength[1]
	buf[3] = littleEndianLength[2]
	copy(buf[WordLen:], msg)

	e.buf = append(e.buf, buf...)
}

func (e *Encoder) PutRawBytes(s []byte) {
	e.buf = append(e.buf, s...)
}

func (e *Encoder) PutVector(v interface{}) {
	slice := dry.SliceToInterfaceSlice(v)
	buf := NewEncoder()
	if v == nil {
		buf.PutCRC(crcVector)
		buf.PutUint(0)
		return
	}

	buf.PutCRC(crcVector)
	buf.PutUint(uint32(len(slice)))

	for _, item := range slice {
		switch val := item.(type) {
		case int8:
			buf.PutInt(int32(val))
		case int16:
			buf.PutInt(int32(val))
		case int32:
			buf.PutInt(val)
		case int64:
			buf.PutLong(val)
		case uint8:
			buf.PutUint(uint32(val))
		case uint16:
			buf.PutUint(uint32(val))
		case uint32:
			buf.PutUint(val)
		case uint64:
			buf.PutLong(int64(val))
		case bool:
			buf.PutBool(val)
		case string:
			buf.PutString(val)
		case []byte:
			buf.PutMessage(val)
		case TLEncoder:
			buf.PutRawBytes(val.Encode())
		default:
			panic("unserializable type: " + reflect.TypeOf(val).String())
		}

	}

	e.buf = append(e.buf, buf.buf...)
}

func (e *Encoder) GetBuffer() []byte {
	return e.buf
}
