package mtproto

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/xelaj/go-dry"

	"github.com/tj/assert"
)

func TestMTProto_SaveSession(t *testing.T) {
	storePath := filepath.Join(os.TempDir(), "session.json")
	defer os.Remove(storePath)

	m := &MTProto{
		authKey:       []byte("some auth key"),
		authKeyHash:   []byte("oooooh that's definetly a key hash!"),
		serverSalt:    0,
		addr:          "1337.228.1488.0",
		tokensStorage: storePath,
	}

	os.Remove(storePath)
	err := m.SaveSession()
	assert.NoError(t, err)

	data, err := ioutil.ReadFile(storePath)
	dry.PanicIfErr(err)

	assert.Equal(t, `{"key":"c29tZSBhdXRoIGtleQ==","hash":"b29vb29oIHRoYXQncyBkZWZpbmV0bHkgYSBrZXkgaGFzaCE="`+
		`,"salt":"AAAAAAAAAAA=","hostname":"1337.228.1488.0"}`, string(data))

}

func TestMTProto_LoadSession(t *testing.T) {
	storePath := filepath.Join(os.TempDir(), "session.json")
	tmpData := `{"key":"c29tZSBhdXRoIGtleQ==","hash":"b29vb29oIHRoYXQncyBkZWZpbmV0bHkgYSBrZXkgaGFzaCE="` +
		`,"salt":"AAAAAAAAAAA=","hostname":"1337.228.1488.0"}`
	ioutil.WriteFile(storePath, []byte(tmpData), 0666)
	defer os.Remove(storePath)

	m := &MTProto{
		tokensStorage: storePath,
	}

	err := m.LoadSession()
	assert.NoError(t, err)

	assert.Equal(t, &MTProto{
		authKey:       []byte("some auth key"),
		authKeyHash:   []byte("oooooh that's definetly a key hash!"),
		serverSalt:    0,
		addr:          "1337.228.1488.0",
		tokensStorage: storePath,
	}, m)
}
