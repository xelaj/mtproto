package tl

import (
	"fmt"
	"reflect"
)

var (
	// used by decoder
	objectByCrc   map[uint32]Object
	enumCrcs      map[uint32]struct{}
	haveFlagCache map[uint32]bool
)

func init() {
	objectByCrc = make(map[uint32]Object)
	enumCrcs = make(map[uint32]struct{})
	haveFlagCache = map[uint32]bool{}
}

func registerObject(o Object) {
	if _, found := objectByCrc[o.CRC()]; found {
		panic(fmt.Errorf("object with that crc already registered: %d", o.CRC()))
	}

	if elem := reflect.ValueOf(o); elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
		if elem.Kind() == reflect.Struct {
			haveFlagCache[o.CRC()] = haveFlag(elem.Interface())
		}
	}

	objectByCrc[o.CRC()] = o
}

func registerEnum(o Object) {
	registerObject(o)
	if _, found := enumCrcs[o.CRC()]; found {
		panic(fmt.Errorf("enum with that crc already registered"))
	}

	enumCrcs[o.CRC()] = struct{}{}
}

func RegisterObjects(obs ...Object) {
	for _, o := range obs {
		registerObject(o)
	}
}

func RegisterEnums(enums ...Object) {
	for _, e := range enums {
		registerEnum(e)
	}
}
