// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

// утилитарные функии, которые не сильно зависят от объявленых структур, но при этом много где используются

package utils

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"
)

// GenerateMessageId отдает по сути unix timestamp но ужасно специфическим образом
// TODO: нахуя нужно битовое и на -4??
func GenerateMessageId() int64 {
	const billion = 1000 * 1000 * 1000
	unixnano := time.Now().UnixNano()
	seconds := unixnano / billion
	nanoseconds := unixnano % billion
	return (seconds << 32) | (nanoseconds & -4)
}

func AuthKeyHash(key []byte) []byte {
	return Sha1Byte(key)[12:20]
}

func GenerateSessionID() int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63() // nolint: gosec потому что начерта?
}

func FullStack() {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			fmt.Fprintln(os.Stderr, string(buf[:n]))
		}
		buf = make([]byte, 2*len(buf))
	}
}

func Sha1Byte(input []byte) []byte {
	r := sha1.Sum(input)
	return r[:]
}
