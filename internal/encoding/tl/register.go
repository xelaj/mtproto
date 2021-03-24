// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

//nolint:gochecknoglobals required global
package tl

import (
	"fmt"
	"reflect"
)

var (
	// used by decoder, guaranteed that types are convertible to tl.Object
	objectByCrc = make(map[uint32]reflect.Type) // this value setting by registerObject(), DO NOT CALL IT BY HANDS
	enumCrcs    = make(map[uint32]null)
)

func registerObject(o Object) {
	if o == nil {
		panic("object is nil")
	}
	objectByCrc[o.CRC()] = reflect.TypeOf(o)
}

func registerEnum(o Object) {
	registerObject(o)
	enumCrcs[o.CRC()] = null{}
}

func RegisterObjects(obs ...Object) {
	for _, o := range obs {
		if val, found := objectByCrc[o.CRC()]; found {
			panic(fmt.Errorf("object with that crc already registered as %v: 0x%08x", val.String(), o.CRC()))
		}

		registerObject(o)
	}
}

func RegisterEnums(enums ...Object) {
	for _, e := range enums {
		if _, found := enumCrcs[e.CRC()]; found {
			panic(fmt.Errorf("enum with that crc already registered"))
		}

		registerEnum(e)
	}
}
