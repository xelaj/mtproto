// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package srp

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xelaj/go-dry"
)

func Test2FA(t *testing.T) {
	type args struct {
		password string
		srpB     []byte
		mp       *ModPow
		// bytes are always 1, telegram doesn't care about it
	}
	tests := []struct {
		args        args
		want        *SrpAnswer
		expectError assert.ErrorAssertionFunc
	}{
		{
			args: args{
				password: "123123",
				srpB: Hexed("9C52401A6A8084EC82F01C3725D3FB448BD2F0C909F9D97726EAC4B7A74172D9" +
					"52F02466BE6734FA274D2B7429E27397F10372D66B400B80A5C5AE3F28B17BF3" +
					"105D7A2D2A885998CDC2DEFC208AEC217AB58859A9ABC2374AD93DC285F4B3FB" +
					"CAFF4143D7888F2425BD2FB711B25609CEB21757D935B1EF2F042173AD0CE2FE" +
					"0E474DAC53914BD25A8A9AED4AEA8953D55CB88621DB37B871EA0D04393AC098" +
					"7F68094CCC9DE8239251375D8FFFD263316CD528C097B7BC9FB919FBEDB76C52" +
					"5DF3413C374EE076D97A1E6D352BB7CC80FD13651B04B32E2E48C5268150842C" +
					"FD07CF855958B1B5EA9C36FDAD697FE3AEC8DCC6B1EFEC36874AF226204676CF"),
				mp: &ModPow{
					Salt1: Hexed("4D11FB6BEC38F9D2546BB0F61E4F1C99A1BC0DB8F0D5F35B1291B37B213123D7ED48F3C6794D495B"),
					Salt2: Hexed("A1B181AAFE88188680AE32860D60BB01"),
					G:     3,
					P: Hexed("C71CAEB9C6B1C9048E6C522F70F13F73980D40238E3E21C14934D037563D930F" +
						"48198A0AA7C14058229493D22530F4DBFA336F6E0AC925139543AED44CCE7C37" +
						"20FD51F69458705AC68CD4FE6B6B13ABDC9746512969328454F18FAF8C595F64" +
						"2477FE96BB2A941D5BCD1D4AC8CC49880708FA9B378E3C4F3A9060BEE67CF9A4" +
						"A4A695811051907E162753B56B0F6B410DBA74D8A84B2A14B3144E0EF1284754" +
						"FD17ED950D5965B4B9DD46582DB1178D169C6BC465B0D6FF9CA3928FEF5B9AE4" +
						"E418FC15E83EBEA0F87FA9FF5EED70050DED2849F47BF959D956850CE929851F" +
						"0D8115F635B105EE2E4E15D04B2454BF6F4FADF034B10403119CD8E3B92FCC5B"),
				},
			},
			want: &SrpAnswer{
				GA: setByte(256, 3),
				M1: Hexed("999DF906BDA2C6CBB52F503406EBA2D0D0503ACE0CC302C38F13EE5010AD4051"),
			},
			expectError: assert.NoError,
		},
	}
	for i := range tests {
		tcase := tests[i]
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			// random это байтовое представление big.NewInt(1), telegram doesn't care about real random
			random := setByte(randombyteLen, 1)
			got, err := getInputCheckPassword(tcase.args.password, tcase.args.srpB, tcase.args.mp, random)
			if !tcase.expectError(t, err) {
				return
			}

			if !assert.Equal(t, tcase.want, got) {
				return
			}
		})
	}
}

func Hexed(in string) []byte {
	res, err := hex.DecodeString(in)
	dry.PanicIfErr(err)
	return res
}

func setByte(size, value int) []byte {
	res := make([]byte, size)
	binary.BigEndian.PutUint32(res[size-4:], uint32(value))
	return res
}
