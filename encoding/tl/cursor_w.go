package tl

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type WriteCursor struct {
	w io.Writer
}

func NewWriteCursor(w io.Writer) *WriteCursor {
	return &WriteCursor{w: w}
}

func (c *WriteCursor) write(b []byte) error {
	n, err := c.w.Write(b)
	if err != nil {
		return err
	}

	if n != len(b) {
		return fmt.Errorf("can't write")
	}

	return nil
}

// PutBool очень специфичный тип, т.к. есть отдельный конструктор под true и false,
// то можно считать, что это две crc константы
func (c *WriteCursor) PutBool(v bool) error {
	buf := make([]byte, WordLen)
	crc := CrcFalse
	if v {
		crc = CrcTrue
	}

	binary.LittleEndian.PutUint32(buf, uint32(crc))
	return c.write(buf)
}

func (c *WriteCursor) PutUint(v uint32) error {
	buf := make([]byte, WordLen)
	binary.LittleEndian.PutUint32(buf, v)
	return c.write(buf)
}

func (c *WriteCursor) PutCRC(v uint32) error {
	return c.PutUint(v) // я так и не понял, кажется что crc это bigendian, но видимо нет
}

func (c *WriteCursor) PutLong(v int64) error {
	buf := make([]byte, WordLen*2)
	binary.LittleEndian.PutUint64(buf, uint64(v))
	return c.write(buf)
}

func (c *WriteCursor) PutDouble(v float64) error {
	buf := make([]byte, WordLen*2)
	binary.LittleEndian.PutUint64(buf, math.Float64bits(v))
	return c.write(buf)
}

func (c *WriteCursor) PutString(msg string) error {
	return c.PutMessage([]byte(msg))
}

func (c *WriteCursor) PutMessage(msg []byte) error {
	if len(msg) < FuckingMagicNumber {
		return c.putTinyBytes(msg)
	}

	return c.putLargeBytes(msg)
}

func (c *WriteCursor) putTinyBytes(msg []byte) error {
	if len(msg) >= FuckingMagicNumber {
		return fmt.Errorf("tiny messages supports maximum 253 elements")
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

	return c.write(buf)
}

func (c *WriteCursor) putLargeBytes(msg []byte) error {
	if len(msg) < FuckingMagicNumber {
		return fmt.Errorf("режим работы в маленьких сообщениях не гарантирован")
	}

	maxLen := 1 << 24 // 3 байта 24 бита, самый первый это 0xfe оставшиеся 3 как раз длина
	if len(msg) > maxLen {
		return fmt.Errorf("message entity too large")
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

	return c.write(buf)
}

func (c *WriteCursor) PutRawBytes(b []byte) error {
	return c.write(b)
}

func (c *WriteCursor) PutVector(v interface{}) error {
	return encodeVector(c, sliceToInterfaceSlice(v))
}
