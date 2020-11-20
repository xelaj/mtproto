// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package tl_test

import (
	"testing"

	"github.com/k0kubun/pp"
	"github.com/xelaj/go-dry"

	"github.com/stretchr/testify/assert"
	"github.com/xelaj/mtproto/encoding/tl"
)

type any = interface{}

// null struct{}

func TestDecode(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		v        any
		expected any
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "authSentCode",
			//           |  CRC || Flag || CRC  |
			data: Hexed("0225005E020000008659BB3D0500000012316637366461306431353531313539363336008C15A372"),
			v:    &AuthSentCode{},
			expected: &AuthSentCode{
				Type: &AuthSentCodeTypeApp{
					Length: 5,
				},
				PhoneCodeHash: "1f76da0d1551159636",
				NextType:      0x72a3158c,
				Timeout:       0,
			},
		},
		{
			name: "poll-results",
			data: Hexed("a3c1dcba1e00000015c4b51c02000000d2da6d3b010000000301020302000000d2da6d3b" +
				"0000000003040506060000000c00000015c4b51c02000000050000000600000005616c616c610000" +
				"15c4b51c00000000"),
			v: &PollResults{},
			expected: &PollResults{
				Min: false,
				Results: []*PollAnswerVoters{
					{
						Chosen:  true,
						Correct: false,
						Option: []byte{
							0x01, 0x02, 0x03,
						},
						Voters: 2,
					},
					{
						Chosen:  false,
						Correct: false,
						Option: []byte{
							0x04, 0x05, 0x06,
						},
						Voters: 6,
					},
				},
				TotalVoters: 12,
				RecentVoters: []int32{
					5,
					6,
				},
				Solution:         "alala",
				SolutionEntities: []MessageEntity{},
			},
		},
		// TODO: отработать возможные ошибки
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				data     = tt.data
				v        = tt.v
				expected = tt.expected
				wantErr  = noErrAsDefault(tt.wantErr)
			)
			err := tl.Decode(data, v)
			if !wantErr(t, err) {
				pp.Println(dry.BytesEncodeHex(string(data)))
				return
			}

			assert.Equal(t, expected, v)
		})
	}
}

func noErrAsDefault(e assert.ErrorAssertionFunc) assert.ErrorAssertionFunc {
	if e == nil {
		return assert.NoError
	}

	return e
}
