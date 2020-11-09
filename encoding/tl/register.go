package tl

import "fmt"

var (
	// used by decoder
	objectByCrc map[uint32]Object
	crcByObject map[Object]uint32
	enumCrcs    map[uint32]struct{}
)

func init() {
	objectByCrc = make(map[uint32]Object)
	crcByObject = make(map[Object]uint32)
	enumCrcs = make(map[uint32]struct{})
}

func registerObject(o Object) {
	if _, found := objectByCrc[o.CRC()]; found {
		panic(fmt.Errorf("object with that crc already registered: %d", o.CRC()))
	}

	if another, found := crcByObject[o]; found {
		panic(fmt.Errorf("crc already associated with another object: %T", another))
	}

	objectByCrc[o.CRC()] = o
	crcByObject[o] = o.CRC()
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
