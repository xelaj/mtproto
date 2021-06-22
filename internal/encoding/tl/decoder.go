// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package tl

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

func Decode(data []byte, res any) error {
	if res == nil {
		return errors.New("can't unmarshal to nil value")
	}
	if reflect.TypeOf(res).Kind() != reflect.Ptr {
		return fmt.Errorf("res value is not pointer as expected. got %v", reflect.TypeOf(res))
	}

	d, err := NewDecoder(bytes.NewReader(data))
	if err != nil {
		return err
	}

	d.decodeValue(reflect.ValueOf(res))
	if d.err != nil {
		return errors.Wrapf(d.err, "decode %T", res)
	}

	return nil
}

// DecodeUnknownObject decodes object from message, when you don't actually know, what message contains.
// due to TL doesn't provide mechanism for understanding is message a int or string, you MUST guarantee, that
// input stream DOES NOT contains any type WITHOUT its CRC code. So, strings, ints, floats, etc. CAN'T BE
// automatically parsed.
//
// expectNextTypes is your predictions how decoder must parse objects hidden under interfaces.
// See Decoder.ExpectTypesInInterface description
func DecodeUnknownObject(data []byte, expectNextTypes ...reflect.Type) (Object, error) {
	d, err := NewDecoder(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if len(expectNextTypes) > 0 {
		d.ExpectTypesInInterface(expectNextTypes...)
	}

	obj := d.decodeRegisteredObject()
	if d.err != nil {
		return nil, errors.Wrap(d.err, "decoding predicted object")
	}
	return obj, nil
}

func (d *Decoder) decodeObject(o Object, ignoreCRC bool) {
	if d.err != nil {
		return
	}

	if !ignoreCRC {
		crcCode := d.PopCRC()
		if d.err != nil {
			d.err = errors.Wrap(d.err, "read crc")
			return
		}

		if crcCode != o.CRC() {
			d.err = fmt.Errorf("invalid crc code: %#v, want: %#v", crcCode, o.CRC())
			return
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
	var flagsetIndex = -1
	if haveFlag(value.Interface()) {
		// getting new cause we need idempotent response
		indexGetter, ok := reflect.New(vtyp).Interface().(FlagIndexGetter)
		if !ok {
			panic("type " + value.Type().String() + " has type bit flag tags, but doesn't inplement tl.FlagIndexGetter")
		}
		flagsetIndex = indexGetter.FlagIndex()
		if flagsetIndex < 0 {
			panic("flag index is below zero, must be index of parameters")
		}

	}

	var bitsetParsed bool
	loopCycles := value.NumField()
	if flagsetIndex >= 0 {
		loopCycles++
	}
	for i := 0; i < loopCycles; i++ {
		// parsing flag is necessary
		if flagsetIndex == i {
			optionalBitSet = d.PopUint()
			if d.err != nil {
				d.err = errors.Wrap(d.err, "read bitset")
				return
			}
			bitsetParsed = true
			continue
		}

		fieldIndex := i
		if bitsetParsed {
			fieldIndex--
		}
		field := value.Field(fieldIndex)

		if _, found := vtyp.Field(fieldIndex).Tag.Lookup(tagName); found {
			info, err := parseTag(vtyp.Field(fieldIndex).Tag)
			if err != nil {
				d.err = errors.Wrap(err, "parse tag")
				return
			}

			if optionalBitSet&(1<<info.index) == 0 {
				continue
			}

			if info.encodedInBitflag {
				field.Set(reflect.ValueOf(true).Convert(field.Type()))
				continue
			}
		}

		if field.Kind() == reflect.Ptr { // && field.IsNil()
			val := reflect.New(field.Type().Elem())
			field.Set(val)
		}

		d.decodeValue(field)
		if d.err != nil {
			d.err = errors.Wrapf(d.err, "decode field '%s'", vtyp.Field(fieldIndex).Name)
			break
		}
	}
}

func (d *Decoder) decodeValue(value reflect.Value) {
	if d.err != nil {
		return
	}
	if m, ok := value.Interface().(Unmarshaler); ok {
		err := m.UnmarshalTL(d)
		if err != nil {
			d.err = err
		}
		return
	}

	val := d.decodeValueGeneral(value)
	if val != nil {
		value.Set(reflect.ValueOf(val).Convert(value.Type()))
		return
	}

	switch value.Kind() { //nolint:exhaustive has default case + more types checked
	// Float64,Int64,Uint32,Int32,Bool,String,Chan, Func, Uintptr, UnsafePointer,Struct,Map,Array,Int,
	// Int8,Int16,Uint,Uint8,Uint16,Uint64,Float32,Complex64,Complex128
	// these values are checked already

	case reflect.Slice:
		if _, ok := value.Interface().([]byte); ok {
			val = d.PopMessage()
		} else {
			val = d.PopVector(value.Type().Elem())
		}

	case reflect.Ptr:
		if o, ok := value.Interface().(Object); ok {
			d.decodeObject(o, false)
		} else {
			d.decodeValue(value.Elem())
		}

		return

	case reflect.Interface:
		val = d.decodeRegisteredObject()

		// if we got slice, we must unwrap it, cause WrappedSlice allowed to exist ONLY in root of returned object
		if v, ok := val.(*WrappedSlice); ok {
			if reflect.TypeOf(v.data).ConvertibleTo(value.Type()) {
				val = v.data
			}
		}

		if d.err != nil {
			d.err = errors.Wrap(d.err, "decode interface")
			return
		}
	default:
		panic("неизвестная штука: " + value.Type().String())
	}

	if d.err != nil {
		return
	}

	value.Set(reflect.ValueOf(val).Convert(value.Type()))
}

// декодирует базовые типы, строчки числа, вот это. если тип не найден возвращает nil
func (d *Decoder) decodeValueGeneral(value reflect.Value) any {
	// value, which is setting into value arg
	var val any

	switch value.Kind() { //nolint:exhaustive has default case
	case reflect.Float64:
		val = d.PopDouble()

	case reflect.Int64:
		val = d.PopLong()

	case reflect.Uint32: // это применимо так же к енумам
		val = d.PopUint()

	case reflect.Int32:
		val = int32(d.PopUint())

	case reflect.Bool:
		val = d.PopBool()

	case reflect.String:
		val = string(d.PopMessage())

	case reflect.Chan, reflect.Func, reflect.Uintptr, reflect.UnsafePointer:
		panic(value.Kind().String() + " does not supported")

	case reflect.Struct:
		d.err = fmt.Errorf("%v must implement tl.Object for decoding (also it must be pointer)", value.Type())

	case reflect.Map:
		d.err = errors.New("map is not ordered object (must order like structs)")

	case reflect.Array:
		d.err = errors.New("array must be slice")

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint64:
		d.err = fmt.Errorf("int kind: %v (must converted to int32, int64 or uint32 explicitly)", value.Kind())
		return nil

	case reflect.Float32, reflect.Complex64, reflect.Complex128:
		d.err = fmt.Errorf("float kind: %s (must be converted to float64 explicitly)", value.Kind())
		return nil

	default:
		// not basic type
		return nil
	}

	return val
}

func (d *Decoder) decodeRegisteredObject() Object {
	crc := d.PopCRC()
	if d.err != nil {
		d.err = errors.Wrap(d.err, "read crc")
	}

	var _typ reflect.Type

	// firstly, we are checking specific crc situations.
	// See https://github.com/xelaj/mtproto/issues/51
	switch crc {
	case CrcVector:
		if len(d.expectedTypes) == 0 {
			d.err = &ErrMustParseSlicesExplicitly{}
			return nil
		}
		_typ = d.expectedTypes[0]
		d.expectedTypes = d.expectedTypes[1:]

		res := d.popVector(_typ.Elem(), true)
		if d.err != nil {
			return nil
		}

		return &WrappedSlice{res}

	case CrcFalse:
		return &PseudoFalse{}

	case CrcTrue:
		return &PseudoTrue{}

	case CrcNull:
		return &PseudoNil{}
	}

	// in other ways we're trying to get object from registred crcs
	var ok bool
	_typ, ok = objectByCrc[crc]
	if !ok {
		msg, err := d.DumpWithoutRead()
		if err != nil {
			return nil
		}

		d.err = &ErrRegisteredObjectNotFound{
			Crc:  crc,
			Data: msg,
		}

		return nil
	}

	o := reflect.New(_typ.Elem()).Interface().(Object)

	if m, ok := o.(Unmarshaler); ok {
		err := m.UnmarshalTL(d)
		if err != nil {
			d.err = err
			return nil
		}
		return o
	}

	if _, isEnum := enumCrcs[crc]; !isEnum {
		d.decodeObject(o, true)
		if d.err != nil {
			d.err = errors.Wrapf(d.err, "decode registered object %T", o)
			return nil
		}
	}

	return o
}
