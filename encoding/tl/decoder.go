package tl

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/k0kubun/pp"
)

func Decode(data []byte, v interface{}) error {
	if v == nil {
		return fmt.Errorf("can't unmarshal to nil value")
	}

	if err := decodeValue(NewReadCursor(bytes.NewReader(data)), reflect.ValueOf(v)); err != nil {
		pp.Println("failed_decode:", data)
		return fmt.Errorf("decode %T: %w", v, err)
	}

	return nil
}

func decodeObject(cur *ReadCursor, o Object, ignoreCRC bool) error {
	if !ignoreCRC {
		crcCode, err := cur.PopCRC()
		if err != nil {
			return fmt.Errorf("read crc: %w", err)
		}

		if crcCode != o.CRC() {
			return fmt.Errorf("invalid crc code: %#v, want: %#v", crcCode, o.CRC())
		}
	}

	value := reflect.ValueOf(o)
	if value.Kind() != reflect.Ptr {
		panic("not a pointer")
	}

	value = reflect.Indirect(value)
	if value.Kind() != reflect.Struct {
		panic("not receiving on struct: " + value.Type().String() + " -> " + value.Kind().String())
	}

	vtyp := value.Type()
	var optionalBitSet uint32

	popFlag, cached := haveFlagCache[o.CRC()]
	if !cached {
		popFlag = haveFlag(value.Interface())
	}

	if popFlag {
		bitset, err := cur.PopUint()
		if err != nil {
			return fmt.Errorf("read bitset: %w", err)
		}

		optionalBitSet = bitset
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)

		if tag, found := vtyp.Field(i).Tag.Lookup(tagName); found {
			info, err := parseTag(tag)
			if err != nil {
				return fmt.Errorf("parse tag: %w", err)
			}

			if optionalBitSet&(1<<info.index) == 0 {
				continue
			}

			if info.encodedInBitflags {
				field.Set(reflect.ValueOf(true).Convert(field.Type()))
				continue
			}
		}

		if field.Kind() == reflect.Ptr { // && field.IsNil()
			val := reflect.New(field.Type().Elem())
			field.Set(val)
		}

		// fmt.Printf("decoding field '%s'\n", vtyp.Field(i).Name)
		if err := decodeValue(cur, field); err != nil {
			return fmt.Errorf("decode field '%s': %w", vtyp.Field(i).Name, err)
		}
	}

	return nil
}

func decodeValue(cur *ReadCursor, value reflect.Value) error {
	if m, ok := value.Interface().(Unmarshaler); ok {
		return m.UnmarshalTL(cur)
	}

	switch value.Kind() {
	case reflect.Float64:
		val, err := cur.PopDouble()
		if err != nil {
			return err
		}

		value.Set(reflect.ValueOf(val).Convert(value.Type()))
	case reflect.Int64:
		val, err := cur.PopLong()
		if err != nil {
			return err
		}

		value.Set(reflect.ValueOf(val).Convert(value.Type()))
	case reflect.Uint32: // это применимо так же к енумам
		val, err := cur.PopUint()
		if err != nil {
			return err
		}

		value.Set(reflect.ValueOf(val).Convert(value.Type()))
	case reflect.Int32:
		val, err := cur.PopUint()
		if err != nil {
			return err
		}

		value.Set(reflect.ValueOf(int32(val)).Convert(value.Type()))
	case reflect.Bool:
		val, err := cur.PopBool()
		if err != nil {
			return err
		}

		value.Set(reflect.ValueOf(val).Convert(value.Type()))
	case reflect.String:
		msg, err := decodeMessage(cur)
		if err != nil {
			return err
		}

		value.Set(reflect.ValueOf(string(msg)).Convert(value.Type()))
	case reflect.Struct:
		panic("struct does not supported")
	case reflect.Slice:
		if _, ok := value.Interface().([]byte); ok {
			msg, err := decodeMessage(cur)
			if err != nil {
				return err
			}

			value.Set(reflect.ValueOf(msg))
			break
		}

		vec, err := decodeVector(cur, value.Type().Elem())
		if err != nil {
			return err
		}

		value.Set(reflect.ValueOf(vec))
	case reflect.Ptr:
		if o, ok := value.Interface().(Object); ok {
			return decodeObject(cur, o, false)
		}

		return decodeValue(cur, value.Elem())
		panic("неизвестная штука: " + value.Type().String())
	case reflect.Interface:
		obj, err := decodeRegisteredObject(cur)
		if err != nil {
			return fmt.Errorf("decode interface: %w", err)
		}

		value.Set(reflect.ValueOf(obj))
	default:
		panic("неизвестная штука: " + value.Type().String())
	}

	return nil
}

func DecodeRegistered(data []byte) (Object, error) {
	ob, err := decodeRegisteredObject(
		NewReadCursor(bytes.NewReader(data)),
	)
	if err != nil {
		return nil, fmt.Errorf("decode registered object: %w", err)
	}

	return ob, nil
}

func decodeRegisteredObject(cur *ReadCursor) (Object, error) {
	crc, err := cur.PopCRC()
	if err != nil {
		return nil, fmt.Errorf("read crc: %w", err)
	}

	o, ok := objectByCrc[crc]
	if !ok {
		msg, err := cur.DumpWithoutRead()
		if err != nil {
			return nil, err
		}

		return nil, ErrRegisteredObjectNotFound{
			Crc:  crc,
			Data: msg,
		}
		// return nil, fmt.Errorf("object with crc %#v not found", crc)
	}

	if o == nil {
		panic("nil object")
	}

	if m, ok := o.(Unmarshaler); ok {
		return o, m.UnmarshalTL(cur)
	}

	if _, isEnum := enumCrcs[crc]; !isEnum {
		err := decodeObject(cur, o, true)
		if err != nil {
			return nil, fmt.Errorf("decode registered object %T: %w", o, err)
		}
	}

	return o, nil
}

func decodeMessage(c *ReadCursor) ([]byte, error) {
	var firstByte byte
	val := []byte{0}

	if err := c.read(val); err != nil {
		return nil, err
	}

	firstByte = val[0]

	realSize := 0
	lenNumberSize := 0 // сколько байт занимаем число обозначающее длину массива
	if firstByte != FuckingMagicNumber {
		realSize = int(firstByte) // это tinyMessage по сути, первый байт является 8битным числом, которое представляет длину сообщения
		lenNumberSize = 1
	} else {
		// иначе это largeMessage с блядским магитческим числом 0xfe
		realSizeBuf := make([]byte, WordLen-1) // WordLen-1 т.к. 1 байт уже прочитали
		if err := c.read(realSizeBuf); err != nil {
			return nil, err
		}

		realSizeBuf = append(realSizeBuf, 0x0) // добиваем до WordLen

		realSize = int(binary.LittleEndian.Uint32(realSizeBuf))
		lenNumberSize = WordLen
	}

	buf := make([]byte, realSize)
	if err := c.read(buf); err != nil {
		return nil, err
	}

	readLen := lenNumberSize + realSize // lenNumberSize это сколько байт ушло на описание длины а realsize это сколько мы по факту прочитали
	if readLen%WordLen != 0 {
		voidBytes := make([]byte, 4-readLen%WordLen)
		if err := c.read(voidBytes); err != nil { // читаем оставшиеся пустые байты. пустые, потому что длина слова 4 байта, может остаться 1,2 или 3 лишних байта
			return nil, err
		}

		for _, b := range voidBytes {
			if b != 0 {
				return nil, fmt.Errorf("some of bytes doesn't equal zero: %#v", voidBytes)
			}
		}
	}

	return buf, nil
}

func decodeVector(c *ReadCursor, as reflect.Type) (interface{}, error) {
	crc, err := c.PopCRC()
	if err != nil {
		return nil, fmt.Errorf("read crc: %w", err)
	}

	if crc != CrcVector {
		return nil, fmt.Errorf("not a vector: %#v, want: %#v", crc, CrcVector)
	}

	size, err := c.PopUint()
	if err != nil {
		return nil, fmt.Errorf("read vector size: %w", err)
	}

	x := reflect.MakeSlice(reflect.SliceOf(as), int(size), int(size))
	for i := 0; i < int(size); i++ {
		var val reflect.Value
		if as.Kind() == reflect.Ptr {
			val = reflect.New(as.Elem())
		} else {
			val = reflect.New(as).Elem()
		}

		if err := decodeValue(c, val); err != nil {
			return nil, err
		}

		x.Index(i).Set(val)
	}

	return x.Interface(), nil
}
