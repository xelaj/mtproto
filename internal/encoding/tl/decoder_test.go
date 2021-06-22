// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package tl_test

import (
	"reflect"
	"testing"

	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
	"github.com/xelaj/go-dry"

	"github.com/xelaj/mtproto/internal/encoding/tl"
)

var (
	True = true // for pointer
)

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
		{
			name:     "nakedBooleanObject#00",
			data:     Hexed("b5757299"),
			v:        &True,
			expected: &tl.PseudoTrue{},
		},
		{
			name:     "nakedBooleanObject#01",
			data:     Hexed("379779bc"),
			v:        &True,
			expected: &tl.PseudoFalse{},
		},
		// TODO: отработать возможные ошибки
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr = noErrAsDefault(tt.wantErr)

			err := tl.Decode(tt.data, tt.v)
			if !tt.wantErr(t, err) {
				pp.Println(dry.BytesEncodeHex(string(tt.data)))
				return
			}
			if err != nil {
				assert.Equal(t, tt.expected, tt.v)
			}
		})
	}
}

func TestDecodeUnknown(t *testing.T) {
	tests := []struct {
		name            string
		data            []byte
		hintsForDecoder []reflect.Type
		expected        any
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "authSentCode",
			//           |  CRC || Flag || CRC  |
			data: Hexed("0225005E020000008659BB3D0500000012316637366461306431353531313539363336008C15A372"),
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
		{
			name:            "predicting-[]int64",
			data:            Hexed("15c4b51c00000000"),
			expected:        tl.ExampleWrappedInt64Slice,
			hintsForDecoder: []reflect.Type{reflect.TypeOf([]int64{})},
		},
		{
			name: "predicting-[]int64ButDataNotWrapped",
			data: Hexed("4646fdfdff00000015c4b51c00000000"),
			expected: &AnyStructWithAnyType{
				SomeInt: 255,
				Data:    []int64{},
			},
			hintsForDecoder: []reflect.Type{reflect.TypeOf([]int64{})},
		},
		{
			name: "predicting-[]int64ButDataNotWrappedAndItsObject",
			data: Hexed("46fd46fdff00000015c4b51c00000000"),
			expected: &AnyStructWithAnyObject{
				SomeInt: 255,
				Data:    tl.ExampleWrappedInt64Slice,
			},
			hintsForDecoder: []reflect.Type{reflect.TypeOf([]int64{})},
		},
		{
			name:     "nakedBooleanObject#00",
			data:     Hexed("b5757299"),
			expected: &tl.PseudoTrue{},
		},
		{
			name:     "nakedBooleanObject#01",
			data:     Hexed("379779bc"),
			expected: &tl.PseudoFalse{},
		},
		{
			name: "issue_59", // https://github.com/xelaj/mtproto/issues/59
			//           crc     id              flag    question string
			data: Hexed("6181e186100000006115f84a0000000015d094d0bed181d182d0b0d182d0bed1" +
				//               slice   len3    crc     long message
				"87d0bdd0be3f000015c4b51c03000000e9c2a96c32d0b4d0bed181d182d0b0d1" +
				// still message
				"82d0bed187d0bdd0be20d182d0bed0bbd18cd0bad0be20d180d0b0d181d0bfd0" +
				// wow, it ends< 1 byte  crc     long message
				"b8d181d0bad0b80001300000e9c2a96c56d0bfd0bed0bcd0b8d0bcd0be20d180" +
				// still message
				"d0b0d181d0bfd0b8d181d0bad0b820d0bdd183d0b6d0bdd18b20d181d0b2d0b8" +
				// still message
				"d0b4d0b5d182d0b5d0bbd18cd181d0bad0b8d0b520d0bfd0bed0bad0b0d0b7d0" +
				// done! wow     1 byte  crc
				"b0d0bdd0b8d18f0001310000e9c2a96ca7d0bad180d0bed0bcd0b520d180d0b0" +
				// still message
				"d181d0bfd0b8d181d0bad0b820d0bad180d0b5d0b4d0b8d182d0bed180d18320" +
				// why it so long??
				"d0bdd183d0b6d0bdd0be20d0b4d0bed0bad0b0d0b7d0b0d182d18c20d0bdd0b0" +
				// omg
				"d0bbd0b8d187d0b8d0b520d182d0b0d0bad0bed0b920d181d183d0bcd0bcd18b" +
				// please stop
				"20d0bdd0b020d0bcd0bed0bcd0b5d0bdd18220d0b7d0b0d0bad0bbd18ed187d0" +
				// HOORAY MESSAGE ENDS!!!                        1 byte
				"b5d0bdd0b8d18f20d0b4d0bed0b3d0bed0b2d0bed180d0b001320000"),
			expected: &Poll{
				ID:             5402091259386920976,
				Closed:         false,
				PublicVoters:   false,
				MultipleChoice: false,
				Quiz:           false,
				Question:       "Достаточно?",
				Answers: []*PollAnswer{
					&PollAnswer{ // don't mind on these texts, i'm too lazy to edit them
						Text: "достаточно только расписки",
						Option: []uint8{
							0x30,
						},
					},
					&PollAnswer{ // don't mind on these texts, i'm too lazy to edit them
						Text: "помимо расписки нужны свидетельские показания",
						Option: []uint8{
							0x31,
						},
					},
					&PollAnswer{ // don't mind on these texts, i'm too lazy to edit them
						Text: "кроме расписки кредитору нужно доказать наличие такой суммы на момент заключения договора",
						Option: []uint8{
							0x32,
						},
					},
				},
				ClosePeriod: 0,
				CloseDate:   0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr = noErrAsDefault(tt.wantErr)

			res, err := tl.DecodeUnknownObject(tt.data, tt.hintsForDecoder...)
			if !tt.wantErr(t, err) {
				pp.Println(dry.BytesEncodeHex(string(tt.data)))
				return
			}

			if err == nil {
				assert.Equal(t, tt.expected, res)
			}
		})
	}
}

func noErrAsDefault(e assert.ErrorAssertionFunc) assert.ErrorAssertionFunc {
	if e == nil {
		return assert.NoError
	}

	return e
}
