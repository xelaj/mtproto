package gen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentificatorNormalization(t *testing.T) {
	type testCase struct {
		in  string
		out string
	}
	tests := []testCase{
		{"user_id", "UserID"},
		{"phone.sendSignalingData", "PhoneSendSignalingData"},
		{"g_a_or_b", "GAOrB"},
		{"embed_url", "EmbedURL"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			assert.Equal(t, tt.out, noramlizeIdentificator(tt.in)) //nolint:golint i see this warn first time
		})
	}
}
