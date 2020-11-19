package tl

import "fmt"

var (
	// used by decoder
	objectByCrc = make(map[uint32]Object) // this value setting by registerObject(), DO NOT CALL IT BY HANDS
	enumCrcs    = make(map[uint32]null)
)

func registerObject(o Object) {
	objectByCrc[o.CRC()] = o
}

func registerEnum(o Object) {
	registerObject(o)
	enumCrcs[o.CRC()] = null{}
}

func RegisterObjects(obs ...Object) {
	for _, o := range obs {
		if _, found := objectByCrc[o.CRC()]; found {
			panic(fmt.Errorf("object with that crc already registered: %d", o.CRC()))
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
