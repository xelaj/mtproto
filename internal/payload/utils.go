// Copyright (c) 2020-2024 Xelaj Software
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package payload

import (
	"crypto/sha1" //nolint:gosec // required by standard
)

type Int128 = [16]byte
type Int256 = [32]byte

type Key [256]byte

//nolint:gosec // spec requires to use sha1
func (k *Key) ID() [8]byte { hash := sha1.Sum(k[:]); return [8]byte(hash[12:20]) }
