package tl

import "fmt"

type ErrRegisteredObjectNotFound struct {
	Crc  uint32
	Data []byte
}

func (e ErrRegisteredObjectNotFound) Error() string {
	return fmt.Sprintf("object with provided crc not registered: 0x%x", e.Crc)
}
