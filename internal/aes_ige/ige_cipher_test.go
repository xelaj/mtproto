// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package ige

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCipher_isCorrectData(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "good_one",
			data:    Hexed("0000000000000000000000000000000000000000000000000000000000000000"),
			wantErr: assert.NoError,
		},
		{
			name:    "smaller_than_want",
			data:    Hexed("0000"),
			wantErr: assert.Error,
		},
		{
			name:    "not_divisible_by_blocks",
			data:    Hexed("0000000000000000000000000000000000000000000000000000000000"),
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantErr := tt.wantErr
			if wantErr == nil {
				wantErr = assert.NoError
			}
			err := isCorrectData(tt.data)
			wantErr(t, err)
		})
	}
}
