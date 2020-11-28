// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package tl_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xelaj/mtproto/encoding/tl"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name    string
		obj     any
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Rights",
			obj: &Rights{
				DeleteMessages: true,
				BanUsers:       true,
			},
			want:    Hexed("D524B25F18000000"),
			wantErr: assert.NoError,
		},
		{
			name: "AccountInstallThemeParams",
			obj: &AccountInstallThemeParams{
				Dark:   true,
				Format: "abc",
				Theme: &InputThemeObj{
					ID:         123,
					AccessHash: 321,
				},
			},
			want:    Hexed("3737E47A0300000003616263E993563C7B000000000000004101000000000000"),
			wantErr: assert.NoError,
		},
		{
			name: "AccountUnregisterDeviceParams",
			obj: &AccountUnregisterDeviceParams{
				TokenType: 1,
				Token:     "foo",
				OtherUids: []int32{
					1337, 228, 322,
				},
			},
			want:    Hexed("BFC476300100000003666F6F15C4B51C0300000039050000E400000042010000"),
			wantErr: assert.NoError,
		},
		{
			name: "respq",
			obj: &ResPQ{
				Nonce:        &tl.Int128{Int: big.NewInt(123)},
				ServerNonce:  &tl.Int128{Int: big.NewInt(321)},
				Pq:           []byte{1, 2, 3},
				Fingerprints: []int64{322, 1337},
			},
			want: Hexed("632416050000000000000000000000000000007B00000000000000000000" +
				"0000000001410301020315C4B51C0200000042010000000000003905000000000000"),
			wantErr: assert.NoError,
		},
		{
			name: "InitConnectionParams",
			obj: &InvokeWithLayerParams{
				Layer: int32(322),
				Query: &InitConnectionParams{
					APIID:          int32(1337),
					DeviceModel:    "abc",
					SystemVersion:  "def",
					AppVersion:     "123",
					SystemLangCode: "en",
					LangCode:       "en",
					Query:          &SomeNullStruct{},
				},
			},
			want: Hexed("0d0d9bda42010000a95ecdc1000000003905000003616263036465660331" +
				"323302656e000000000002656e006b18f9c4"),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tl.Marshal(tt.obj)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
