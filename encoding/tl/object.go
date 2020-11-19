package tl

type Object interface {
	CRC() uint32
}

type FlagIndexGetter interface {
	FlagIndex() int
}

type Marshaler interface {
	MarshalTL(*WriteCursor) error
}

type Unmarshaler interface {
	UnmarshalTL(*ReadCursor) error
}
