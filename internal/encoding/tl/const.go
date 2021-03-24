// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package tl

const (
	WordLen   = 4           // размер слова в TL (32 бита)
	LongLen   = WordLen * 2 // int64 8 байт занимает
	DoubleLen = WordLen * 2 // float64 8 байт занимает
	Int128Len = WordLen * 4 // int128 16 байт
	Int256Len = WordLen * 8 // int256 32 байт

	// Блядские магические числа
	FuckingMagicNumber = 0xfe // 253 элемента максимум можно закодировать в массиве элементов

	// https://core.telegram.org/schema/mtproto
	CrcVector = 0x1cb5c415
	CrcFalse  = 0xbc799737
	CrcTrue   = 0x997275b5
	CrcNull   = 0x56730bcc

	bitsInByte = 8 // cause we don't want store magic numbers
)
