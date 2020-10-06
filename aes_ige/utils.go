package ige

// побайтовый xor
func xor(dst, src []byte) {
	for i := range dst {
		dst[i] ^= src[i]
	}
}
