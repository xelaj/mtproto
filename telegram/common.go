package telegram

import (
	"crypto/rsa"
	"fmt"

	dry "github.com/xelaj/go-dry"

	"github.com/xelaj/mtproto"
	"github.com/xelaj/mtproto/keys"
	"github.com/xelaj/mtproto/serialize"
)

func init() {
	keyfile := "~/go/src/github.com/xelaj/mtproto/keys/keys.pem"
	TelegramPublicKeys, err := keys.ReadFromFile(keyfile)
	dry.PanicIfErr(err)
	choosedKey := dry.RandomChoose(dry.InterfaceSlice(TelegramPublicKeys)...).(*rsa.PublicKey)

	_ = choosedKey
}

type Client struct {
	*mtproto.MTProto
}

//----------------------------------------------------------------------------

type InvokeWithLayerParams struct {
	Layer int32
	Query serialize.TLEncoder
}

func (_ *InvokeWithLayerParams) CRC() uint32 {
	return 0xda9b0d0d
}

func (t *InvokeWithLayerParams) Encode() []byte {
	buf := serialize.NewEncoder()
	buf.PutUint(t.CRC())
	buf.PutInt(t.Layer)
	buf.PutRawBytes(t.Query.Encode())
	return buf.Result()
}

func (t *InvokeWithLayerParams) DecodeFrom(d *serialize.Decoder) {
	panic("makes no sense")
}

func (m *Client) InvokeWithLayer(layer int, query serialize.TLEncoder) (serialize.TL, error) {
	data, err := m.MakeRequest(&InvokeWithLayerParams{
		Layer: int32(layer),
		Query: query,
	})
	if err != nil {
		return nil, fmt.Errorf("sending InvokeWithLayer: %w", err)
	}

	return data, nil
}

type InitConnectionParams struct {
	__flagsPosition struct{}
	ApiID           int32
	DeviceModel     string
	SystemVersion   string
	AppVersion      string
	SystemLangCode  string
	LangPack        string
	LangCode        string
	Proxy           *InputClientProxy `flag:"0"`
	Params          JSONValue         `flag:"1"`
	Query           serialize.TLEncoder
}

func (_ *InitConnectionParams) CRC() uint32 {
	return 0xc1cd5ea9
}

func (t *InitConnectionParams) Encode() []byte {
	var flag uint32
	if t.Proxy != nil {
		flag |= 1 << 0
	}
	if t.Params != nil {
		flag |= 1 << 1
	}

	buf := serialize.NewEncoder()
	buf.PutUint(t.CRC())
	buf.PutUint(flag)
	buf.PutInt(t.ApiID)
	buf.PutString(t.DeviceModel)
	buf.PutString(t.SystemVersion)
	buf.PutString(t.AppVersion)
	buf.PutString(t.SystemLangCode)
	buf.PutString(t.LangPack)
	buf.PutString(t.LangCode)
	if t.Proxy != nil {
		buf.PutRawBytes(t.Proxy.Encode())
	}
	if t.Params != nil {
		buf.PutRawBytes(t.Params.Encode())
	}
	buf.PutRawBytes(t.Query.Encode())
	return buf.Result()
}

func (t *InitConnectionParams) DecodeFrom(d *serialize.Decoder) {
	panic("makes no sense")
}

func (m *Client) InitConnection(params InitConnectionParams) (serialize.TL, error) {
	data, err := m.MakeRequest(&params)
	if err != nil {
		return nil, fmt.Errorf("sending InitConnection: %w", err)
	}

	return data, nil
}
