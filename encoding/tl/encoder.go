package tl

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

func Marshal(v any) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	encoder := NewEncoder(buf)
	encoder.encodeValue(reflect.ValueOf(v))
	if err := encoder.CheckErr(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// FIXME: high cyclomatic complexity
func (c *Encoder) encodeValue(value reflect.Value) { //nolint:gocyclo working on it
	if m, ok := value.Interface().(Marshaler); ok {
		if c.err != nil {
			return
		}
		c.err = m.MarshalTL(c)
		return
	}

	switch value.Kind() {
	case reflect.Uint32:
		c.PutUint(uint32(value.Uint()))

	case reflect.Int32:
		c.PutUint(uint32(value.Int()))

	case reflect.Int64:
		c.PutLong(value.Int())

	case reflect.Float64:
		c.PutDouble(value.Float())

	case reflect.Bool:
		c.PutBool(value.Bool())

	case reflect.String:
		c.PutString(value.String())

	case reflect.Ptr, reflect.Interface:
		if value.IsNil() {
			c.err = errors.New("value can't be nil")
			return
		}
		if reflect.Indirect(value).Kind() == reflect.Struct {
			c.encodeStruct(value)
		} else {
			c.err = fmt.Errorf("unsupported type: %s", value.Type().String())
		}

	case reflect.Slice:
		if b, ok := value.Interface().([]byte); ok {
			c.PutMessage(b)
			break
		}

		c.encodeVector(sliceToInterfaceSlice(value.Interface())...)

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint64:
		c.err = fmt.Errorf("int kind: %s (must be converted to int32, int64 or uint32 explicitly)", value.Kind())

	case reflect.Float32, reflect.Complex64, reflect.Complex128:
		c.err = fmt.Errorf("float kind: %s (must be converted to float64 explicitly)", value.Kind())

	default:
		c.err = fmt.Errorf("unsupported type: %s", value.Type().String())
	}
}

// FIXME: high cyclomatic complexity
// v must be pointer to struct
func (c *Encoder) encodeStruct(v reflect.Value) { //nolint:gocyclo WIP
	if c.err != nil {
		return
	}

	o, ok := v.Interface().(Object)
	if !ok {
		c.err = errors.New(v.Type().String() + " doesn't implement tl.Object interface")
		return
	}

	var hasFlagsField bool
	var flag uint32
	var flagIndex int
	g, ok := v.Interface().(FlagIndexGetter)
	if ok {
		hasFlagsField = true
		flagIndex = g.FlagIndex()
	}

	v = reflect.Indirect(v)

	// what we checked and what we know about value:
	// 1) it's not Marshaler (marshaler object already parsing in c.encodeValue())
	// 2) implements tl.Object
	// 3) definitely struct (we don't call encodeStruct(), only in c.encodeValue())
	// 4) not nil (structs can't be nil, only pointers and interfaces)
	c.PutCRC(o.CRC())
	var tmpObjects = make([]reflect.Value, 0)

	vtyp := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if hasFlagsField && flagIndex == i {
			tmpObjects = append(tmpObjects, reflect.ValueOf(0))
		}

		info, err := parseTag(string(vtyp.Field(i).Tag))
		if err != nil {
			c.err = errors.Wrapf(err, "parsing tag of field %v", vtyp.Field(i).Name)
			return
		}

		if info == nil {
			// если тега нет, то это обязательное поле, значит 100% записываем
			tmpObjects = append(tmpObjects, v.Field(i))
			continue
		}

		if info.ignore {
			continue
		}

		if info.encodedInBitflag && vtyp.Field(i).Type.Kind() != reflect.Bool {
			c.err = fmt.Errorf("field '%s': only bool values can be encoded in bitflag", vtyp.Field(i).Name)
			return
		}

		fieldVal := v.Field(i)
		if !fieldVal.IsZero() {
			// тег есть, это 100% опциональное поле
			flag |= 1 << info.index
			if info.encodedInBitflag {
				continue
			}

			tmpObjects = append(tmpObjects, v.Field(i))
			continue
		}
	}

	for i, elem := range tmpObjects {
		if i == flagIndex {
			c.PutUint(flag)
			continue
		}

		c.encodeValue(elem)
		if c.err != nil {
			return
		}
	}
}

func (c *WriteCursor) encodeVector(slice ...any) {
	c.PutCRC(CrcVector)
	c.PutUint(uint32(len(slice)))

	for i, item := range slice {
		c.encodeValue(reflect.ValueOf(item))
		if c.CheckErr() != nil {
			c.err = errors.Wrapf(c.err, "[%v]", i)
		}
	}
}
