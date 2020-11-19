package tl

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"reflect"
)

type ReadCursor struct {
	r io.Reader
}

func NewReadCursor(r io.Reader) *ReadCursor {
	return &ReadCursor{r: r}
}

func (c *ReadCursor) read(buf []byte) error {
	n, err := c.r.Read(buf)
	if err != nil {
		return err
	}

	if n != len(buf) {
		return fmt.Errorf("can't read full buffer")
	}

	return nil
}

func (c *ReadCursor) PopLong() (int64, error) {
	val := make([]byte, LongLen)
	if err := c.read(val); err != nil {
		return 0, err
	}

	return int64(binary.LittleEndian.Uint64(val)), nil
}

func (c *ReadCursor) PopDouble() (float64, error) {
	val := make([]byte, DoubleLen)
	if err := c.read(val); err != nil {
		return 0, err
	}

	return math.Float64frombits(binary.LittleEndian.Uint64(val)), nil
}

func (c *ReadCursor) PopUint() (uint32, error) {
	val := make([]byte, WordLen)
	if err := c.read(val); err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(val), nil
}

func (c *ReadCursor) PopRawBytes(size int) ([]byte, error) {
	val := make([]byte, size)
	if err := c.read(val); err != nil {
		return nil, err
	}

	return val, nil
}

func (c *ReadCursor) PopBool() (bool, error) {
	crc, err := c.PopUint()
	if err != nil {
		return false, err
	}

	switch crc {
	case CrcTrue:
		return true, nil
	case CrcFalse:
		return false, nil
	default:
		return false, fmt.Errorf("not a bool value, actually: %#v", crc)
	}
}

// func (c *ReadCursor) PopNull() (any, error) {
// 	crc, err := c.PopUint()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if crc != CrcNull {
// 		return nil, fmt.Errorf("not a null value, actually: %#v", crc)
// 	}

// 	return nil, nil
// }

func (c *ReadCursor) PopCRC() (uint32, error) {
	return c.PopUint() // я так и не понял, кажется что crc это bigendian, но видимо нет
}

func (c *ReadCursor) GetRestOfMessage() ([]byte, error) {
	return ioutil.ReadAll(c.r)
}

func (c *ReadCursor) DumpWithoutRead() ([]byte, error) {
	data, err := ioutil.ReadAll(c.r)
	if err != nil {
		return nil, err
	}

	cp := make([]byte, len(data))
	copy(cp, data)
	c.r = ioutil.NopCloser(bytes.NewReader(cp))

	return data, nil
}

func (c *ReadCursor) PopVector(as reflect.Type) (any, error) {
	return decodeVector(c, as)
}

func (c *ReadCursor) PopMessage() ([]byte, error) {
	return decodeMessage(c)
}
