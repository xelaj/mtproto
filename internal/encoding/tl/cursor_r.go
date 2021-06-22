// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package tl

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"reflect"

	"github.com/pkg/errors"
)

// A Decoder reads and decodes TL values from an input stream.
type Decoder struct {
	buf *bytes.Reader
	err error

	// see Decoder.ExpectTypesInInterface description
	expectedTypes []reflect.Type
}

// NewDecoder returns a new decoder that reads from r.
// Unfortunately, decoder can't work with part of data, so reader must be read all before decoding.
func NewDecoder(r io.Reader) (*Decoder, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "reading data before decoding")
	}

	return &Decoder{buf: bytes.NewReader(data)}, nil
}

// ExpectTypesInInterface defines, how decoder must parse implicit objects.
// how does expectedTypes works:
// So, imagine: you want parse []int32, but also you can get []int64, or SomeCustomType, or even [][]bool.
// How to deal it?
// expectedTypes store your predictions (like "if you got unknown type, parse it as int32, not int64")
// also, if you have predictions deeper than first unknown type, you can say decoder to use predicted vals
//
// So, next time, when you'll have strucre object with interface{} which expect contains []float64 or sort
// of — use this feature via d.ExpectTypesInInterface()
func (d *Decoder) ExpectTypesInInterface(types ...reflect.Type) {
	d.expectedTypes = types
}

func (d *Decoder) read(buf []byte) {
	if d.err != nil {
		return
	}

	n, err := d.buf.Read(buf)
	if err != nil {
		d.unread(n)
		d.err = err
		return
	}

	if n != len(buf) {
		d.unread(n)
		d.err = fmt.Errorf("buffer weren't fully read: want %v bytes, got %v", len(buf), n)
		return
	}
}

func (d *Decoder) unread(count int) {
	for i := 0; i < count; i++ {
		if d.buf.UnreadByte() != nil {
			return
		}
	}
}

func (d *Decoder) PopLong() int64 {
	val := make([]byte, LongLen)
	d.read(val)
	if d.err != nil {
		return 0
	}

	return int64(binary.LittleEndian.Uint64(val))
}

func (d *Decoder) PopDouble() float64 {
	val := make([]byte, DoubleLen)
	d.read(val)
	if d.err != nil {
		return 0
	}

	return math.Float64frombits(binary.LittleEndian.Uint64(val))
}

func (d *Decoder) PopUint() uint32 {
	val := make([]byte, WordLen)
	d.read(val)
	if d.err != nil {
		return 0
	}

	return binary.LittleEndian.Uint32(val)
}

func (d *Decoder) PopRawBytes(size int) []byte {
	val := make([]byte, size)
	d.read(val)
	if d.err != nil {
		return nil
	}

	return val
}

func (d *Decoder) PopBool() bool {
	crc := d.PopUint()
	if d.err != nil {
		return false
	}

	switch crc {
	case CrcTrue:
		return true
	case CrcFalse:
		return false
	default:
		d.err = fmt.Errorf("not a bool value, actually: %#v", crc)
		return false
	}
}

func (d *Decoder) PopNull() {
	crc := d.PopUint()
	if d.err != nil {
		return
	}

	if crc != CrcNull {
		d.err = fmt.Errorf("not a null value, actually: %#v", crc)
		return
	}
}

func (d *Decoder) PopCRC() uint32 {
	return d.PopUint() // я так и не понял, кажется что crc это bigendian, но видимо нет
}

func (d *Decoder) PopInt() int32 {
	return int32(d.PopUint())
}

func (d *Decoder) GetRestOfMessage() ([]byte, error) {
	return ioutil.ReadAll(d.buf)
}

func (d *Decoder) DumpWithoutRead() ([]byte, error) {
	data, err := ioutil.ReadAll(d.buf)
	if err != nil {
		return nil, err
	}

	d.unread(len(data))
	return data, nil
}

func (d *Decoder) PopVector(as reflect.Type) any {
	return d.popVector(as, false)
}

func (d *Decoder) popVector(as reflect.Type, ignoreCRC bool) any {
	if d.err != nil {
		return nil
	}
	if !ignoreCRC {
		crc := d.PopCRC()
		if d.err != nil {
			d.err = errors.Wrap(d.err, "read crc")
			return nil
		}

		if crc != CrcVector {
			d.err = fmt.Errorf("not a vector: 0x%08x, want: 0x%08x", crc, CrcVector)
			return nil
		}
	}

	size := d.PopUint()
	if d.err != nil {
		d.err = errors.Wrap(d.err, "read vector size")
		return nil
	}

	x := reflect.MakeSlice(reflect.SliceOf(as), int(size), int(size))
	for i := 0; i < int(size); i++ {
		var val reflect.Value
		if as.Kind() == reflect.Ptr {
			val = reflect.New(as.Elem())
		} else {
			val = reflect.New(as).Elem()
		}

		d.decodeValue(val)
		if d.err != nil {
			return nil
		}

		x.Index(i).Set(val)
	}

	return x.Interface()
}

func (d *Decoder) PopMessage() []byte {
	val := []byte{0}

	d.read(val)
	if d.err != nil {
		return nil
	}

	// значение первого байта
	firstByte := val[0]
	// сколько РЕАЛЬНЫХ байт занимает сообщение. от 0 до бесконечности
	var realSize int
	// сколько байт занимаем число обозначающее длину массива (1 или 4)
	var lenNumberSize int

	if firstByte != FuckingMagicNumber {
		// это tinyMessage по сути, первый байт является 8битным числом, которое представляет длину сообщения
		realSize = int(firstByte)
		lenNumberSize = 1
	} else {
		// иначе это largeMessage с блядским магитческим числом 0xfe
		val = make([]byte, WordLen-1) // WordLen-1 т.к. 1 байт уже прочитали
		d.read(val)
		if d.err != nil {
			d.err = errors.Wrapf(d.err, "reading last %v bytes of message size", WordLen-1)
			return nil
		}

		val = append(val, 0x0) // добиваем до WordLen

		realSize = int(binary.LittleEndian.Uint32(val))
		lenNumberSize = WordLen
	}

	// этот буффер и будет уже реальным собщением
	buf := make([]byte, realSize)
	d.read(buf)
	if d.err != nil {
		d.err = errors.Wrapf(d.err, "reading message data with len of %v", realSize)
		return nil
	}

	// lenNumberSize это сколько байт мы УЖЕ прочитали. надо узнать, сколько еще байт нам нужно дочитать,
	// что бы их количество делилось на размер слова (4 байта)
	readLen := lenNumberSize + realSize
	if readLen%WordLen != 0 {
		// читаем оставшиеся пустые байты. пустые, потому что длина слова 4 байта, может остаться 1-3 лишних
		// байта
		voidBytes := make([]byte, WordLen-readLen%WordLen)
		d.read(voidBytes)
		if d.err != nil {
			d.err = errors.Wrapf(d.err, "reading %v last void bytes", WordLen-readLen%WordLen)
			return nil
		}

		for _, b := range voidBytes {
			if b != 0 {
				d.err = fmt.Errorf("some of void bytes doesn't equal zero: %#v", voidBytes)
				return nil
			}
		}
	}

	return buf
}
