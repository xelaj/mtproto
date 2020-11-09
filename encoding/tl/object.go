package tl

type Object interface {
	CRC() uint32
}

type Marshaler interface {
	MarshalTL(*WriteCursor) error
}

type Unmarshaler interface {
	UnmarshalTL(*ReadCursor) error
}
