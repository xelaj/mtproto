// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package tl

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type Encoder struct {
	w io.Writer
	// this error is last unsuccessful write into w. if this err != nil,
	// write() method will not write enay data
	err error
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) write(b []byte) {
	if e.err != nil {
		return
	}

	n, err := e.w.Write(b)
	if err != nil {
		e.err = err
		return
	}

	if n != len(b) {
		e.err = &ErrorPartialWrite{Has: n, Want: len(b)}
	}
}

// CheckErr must call after encoding has been finished. if this func returns not nil value, encoding has
// failed, and you shouldn't use its result
func (e *Encoder) CheckErr() error {
	return e.err
}

// PutBool очень специфичный тип, т.к. есть отдельный конструктор под true и false,
// то можно считать, что это две crc константы
func (e *Encoder) PutBool(v bool) {
	crc := CrcFalse
	if v {
		crc = CrcTrue
	}

	e.PutUint(uint32(crc))
}

func (e *Encoder) PutUint(v uint32) {
	buf := make([]byte, WordLen)
	binary.LittleEndian.PutUint32(buf, v)
	e.write(buf)
}

// PutCRC is an alias for Encoder.PutUint. It uses only for understanding what your code do (like
// self-documented code)
func (e *Encoder) PutCRC(v uint32) {
	e.PutUint(v) // я так и не понял, кажется что crc это bigendian, но видимо нет
}

func (e *Encoder) PutInt(v int32) {
	e.PutUint(uint32(v))
}

func (e *Encoder) PutLong(v int64) {
	buf := make([]byte, LongLen)
	binary.LittleEndian.PutUint64(buf, uint64(v))
	e.write(buf)
}

func (e *Encoder) PutDouble(v float64) {
	buf := make([]byte, DoubleLen)
	binary.LittleEndian.PutUint64(buf, math.Float64bits(v))
	e.write(buf)
}

func (e *Encoder) PutMessage(msg []byte) {
	if len(msg) < FuckingMagicNumber {
		e.putTinyBytes(msg)
	} else {
		e.putLargeBytes(msg)
	}
}

func (e *Encoder) PutString(msg string) {
	e.PutMessage([]byte(msg))
}

func (e *Encoder) putTinyBytes(msg []byte) {
	if len(msg) >= FuckingMagicNumber {
		// it's panicing, cause, you shouldn' call this func by your
		// hands. panic required for internal purposes
		panic("tiny messages supports maximum 253 elements")
	}

	// здесь мы считаем, что длинна итогового сообщения должна делиться на 4
	// (32/8 = 4, 4 байта одно слово)
	// поэтому мы создаем buf с размером, достаточным для пихания
	// массива + 0-3 доп байта что бы итог делился на 4
	realBytesLen := 1 + len(msg) // adding 1, cause we need to store length, realBytesLen doesn't store
	factBytesLen := realBytesLen
	if factBytesLen%WordLen > 0 {
		factBytesLen += WordLen - factBytesLen%WordLen
	}

	buf := make([]byte, factBytesLen)
	buf[0] = byte(len(msg)) // пихаем в первый байт размер сообщения
	copy(buf[1:], msg)

	e.write(buf)
}

func (e *Encoder) putLargeBytes(msg []byte) {
	if len(msg) < FuckingMagicNumber {
		// it's panicing, cause, you shouldn' call this func by your
		// hands. panic required for internal purposes
		panic("can't save binary stream with length less than 253 bytes")
	}

	maxLen := 1 << 24 // 3 байта 24 бита, самый первый это 0xfe оставшиеся 3 как раз длина
	if len(msg) > maxLen {
		e.err = fmt.Errorf("message entity too large: expect less than %v, got %v", maxLen, len(msg))
		return
	}

	realBytesLen := WordLen + len(msg) // первым идет магический байт и 3 байта длины
	factBytesLen := realBytesLen
	if factBytesLen%WordLen > 0 {
		factBytesLen += WordLen - factBytesLen%WordLen
	}

	// FIXME: this thing is uint number too. so, it can decode more simpler
	littleEndianLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(littleEndianLength, uint32(len(msg)))

	buf := make([]byte, factBytesLen)
	buf[0] = byte(FuckingMagicNumber)
	buf[1] = littleEndianLength[0]
	buf[2] = littleEndianLength[1]
	buf[3] = littleEndianLength[2]
	copy(buf[WordLen:], msg)

	e.write(buf)
}

func (e *Encoder) PutRawBytes(b []byte) {
	e.write(b)
}

func (e *Encoder) PutVector(v any) {
	e.encodeVector(sliceToInterfaceSlice(v)...)
}
