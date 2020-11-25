package mtproto

// internal errors for internal purposes

type errorSessionConfigsChanged struct{}

func (*errorSessionConfigsChanged) Error() string {
	return "session configuration was changed, need to repeat request"
}

func (*errorSessionConfigsChanged) CRC() uint32 {
	panic("makes no sense")
}


