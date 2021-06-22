// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package tl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xelaj/mtproto/internal/encoding/tl"
)

func TestUnwrapNativeTypes(t *testing.T) {
	tests := []struct {
		name     string
		v        tl.Object
		expected any
	}{
		{
			name:     "simple_bool",
			v:        &tl.PseudoTrue{},
			expected: true,
		},
		{
			name:     "simple_nil",
			v:        &tl.PseudoNil{},
			expected: nil,
		},
		{
			name:     "simple_slice",
			v:        tl.ExampleWrappedInt64Slice,
			expected: []int64{},
		},
		// TODO: more?
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tl.UnwrapNativeTypes(tt.v))
		})
	}
}
