// утилитарные функии, которые не сильно зависят от объявленых структур, но при этом много где используются

package utils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/xelaj/go-dry"
)

const (
	wordLen = 4

	// если длина пакета больше или равн 127 слов, то кодируем 4 байтами, 1 это магическое число, оставшиеся 3 — дилна
	// https://core.telegram.org/mtproto/mtproto-transports#abridged
	magicValueSizeMoreThanSingleByte = 0x7f
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
	return dry.Sha1Byte(key)[12:20]
}

func PacketLengthMTProtoCompatible(data []byte) []byte {
	packetSizeInWords := len(data) / wordLen
	if packetSizeInWords < 127 {
		return []byte{byte(packetSizeInWords)}
	}
	buf := make([]byte, wordLen)
	binary.LittleEndian.PutUint32(buf, uint32(packetSizeInWords))

	buf = append([]byte{magicValueSizeMoreThanSingleByte}, buf[:3]...)
	return buf
}

var (
	ErrPacketSizeIsBigger = errors.New("packet size is more than 127 bytes, require 4 bytes value")
)

// исходя из переданного числа в bytestoGetInfo считает количество СЛОВ и отдает количество БАЙТ которые нужно прочитать
func GetPacketLengthMTProtoCompatible(bytesToGetInfo []byte) (int, error) {
	if len(bytesToGetInfo) != 1 && len(bytesToGetInfo) != 4 {
<<<<<<< HEAD
		panic(fmt.Sprintf("invalid size of bytes. require only 1 or 4, got %d", len(bytesToGetInfo)))
=======
		return 0, fmt.Errorf("invalid size of bytes. require only 1 or 4, got %v", len(bytesToGetInfo))
>>>>>>> 🏇 multiple changes
	}

	if bytesToGetInfo[0] != magicValueSizeMoreThanSingleByte {
		return int(bytesToGetInfo[0]) * wordLen, nil
	}

	if len(bytesToGetInfo) == 1 {
		return 0, ErrPacketSizeIsBigger
	}

	// 3 последующих байта сейчас прочтем, последний для доведения до uint32, то есть в буффере
	// значение будет 0x00ffffff, где f любой байт, который показывает число
	buf := append(bytesToGetInfo, 0x00)

	value := binary.LittleEndian.Uint32(buf)
	return int(value) * wordLen, nil
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
