package deeplinks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xelaj/mtproto/telegram/deeplinks"
)

func TestResolveLink(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		want    deeplinks.Deeplink
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "some bot link",
			link: "t.me/BotFather",
			want: &deeplinks.ResolveParameters{
				Domain: "botfather",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr = noErrAsDefault(tt.wantErr)

			res, err := deeplinks.Resolve(tt.link)
			if !tt.wantErr(t, err) {
				return
			}

			if err == nil {
				assert.Equal(t, tt.want, res)
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
