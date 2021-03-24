// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package ige

// побайтовый xor
func xor(dst, src []byte) {
	for i := range dst {
		dst[i] ^= src[i]
	}
}
