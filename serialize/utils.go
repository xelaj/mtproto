package serialize

import (
	"fmt"
	"math/big"
)

var bitlen = []int{
	1 << 3,  // 8
	1 << 4,  // 16
	1 << 5,  // 32
	1 << 6,  // 64
	1 << 7,  // 128
	1 << 8,  // 256
	1 << 9,  // 512
	1 << 10, // 1024
	1 << 11, // 2048
}

func bigIntToBytes(i *big.Int, bitSize int) ([]byte, error) {
	vbytes := i.Bytes()
	for i, b := range bitlen {
		if b == bitSize {
			break
		}

		if i == len(bitlen)-1 {
			return nil, fmt.Errorf("bitsize not squaring by 2: bitsize %v", bitSize)
		}
	}

	offset := bitSize/8 - len(vbytes)
	if offset < 0 {
		return nil, fmt.Errorf("bitsize too small: have %v, want at least %v", bitSize, vbytes)
	}

	return append(make([]byte, offset), vbytes...), nil
}
