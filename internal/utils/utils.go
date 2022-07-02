// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

// утилитарные функии, которые не сильно зависят от объявленых структур, но при этом много где используются

package utils

import (
	"crypto/sha1"
	"math/rand"
	"time"
)

// GenerateMessageId возвращает unix timestamp
func GenerateMessageId() int64 {
	return time.Now().UnixNano()
}

func AuthKeyHash(key []byte) []byte {
	return Sha1Byte(key)[12:20]
}

func GenerateSessionID() int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63() // nolint: gosec потому что начерта?
}

func Sha1Byte(input []byte) []byte {
	r := sha1.Sum(input)
	return r[:]
}
