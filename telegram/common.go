package telegram

import (
	"reflect"
	"runtime"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/xelaj/errs"
	dry "github.com/xelaj/go-dry"

	"github.com/xelaj/mtproto"
	"github.com/xelaj/mtproto/keys"
	"github.com/xelaj/mtproto/serialize"
)

type Client struct {
	*mtproto.MTProto
	config       *ClientConfig
	serverConfig *Config
}

type ClientConfig struct {
	SessionFile    string
	ServerHost     string
	PublicKeysFile string
	DeviceModel    string
	SystemVersion  string
	AppVersion     string
	AppID          int
	AppHash        string
}

func NewClient(c ClientConfig) (*Client, error) { //nolint: gocritic arg is not ptr cause we call
	//                                                               it only once, don't care
	//                                                               about copying big args.
	if !dry.FileExists(c.PublicKeysFile) {
		return nil, errs.NotFound("file", c.PublicKeysFile)
	}

	if !dry.PathIsWirtable(c.SessionFile) {
		return nil, errs.Permission(c.SessionFile).Scope("write")
	}

	if c.DeviceModel == "" {
		c.DeviceModel = "Unknown"
	}

	if c.SystemVersion == "" {
		c.SystemVersion = runtime.GOOS + "/" + runtime.GOARCH
	}

	if c.AppVersion == "" {
		c.AppVersion = "v0.0.0"
	}

	publicKeys, err := keys.ReadFromFile(c.PublicKeysFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading public keys")
	}

	m, err := mtproto.NewMTProto(mtproto.Config{
		AuthKeyFile: c.SessionFile,
		ServerHost:  c.ServerHost,
		PublicKey:   publicKeys[0],
	})
	if err != nil {
		return nil, errors.Wrap(err, "setup common MTProto client")
	}

	err = m.CreateConnection()
	if err != nil {
		return nil, errors.Wrap(err, "creating connection")
	}

	client := &Client{
		MTProto: m,
		config:  &c,
	}

	client.AddCustomServerRequestHandler(client.handleSpecialRequests())

	resp, err := client.InvokeWithLayer(ApiVersion, &InitConnectionParams{
		ApiID:          int32(c.AppID),
		DeviceModel:    c.DeviceModel,
		SystemVersion:  c.SystemVersion,
		AppVersion:     c.AppVersion,
		SystemLangCode: "en", // can't be edited, cause docs says that a single possible parameter
		LangCode:       "en",
		Query:          &HelpGetConfigParams{},
	})

	if err != nil {
		return nil, errors.Wrap(err, "getting server configs")
	}

	config, ok := resp.(*Config)
	if !ok {
		return nil, errors.New("got wrong response: " + reflect.TypeOf(resp).String())
	}

	client.serverConfig = config

	dcList := make(map[int]string)
	for _, dc := range config.DcOptions {
		if dc.Cdn {
			continue
		}

		dcList[int(dc.Id)] = dc.IpAddress
	}
	client.SetDCStorages(dcList)

	return client, nil
}

func (c *Client) handleSpecialRequests() func(interface{}) bool {
	return func(i interface{}) bool {
		switch msg := i.(type) {
		case *UpdatesObj:
			pp.Println(msg, "UPDATE")
			return true
		case *UpdateShort:
			pp.Println(msg, "SHORT UPDATE")
			return true
		}

		return false
	}
}

//----------------------------------------------------------------------------

type InvokeWithLayerParams struct {
	Layer int32
	Query serialize.TLEncoder
}

func (*InvokeWithLayerParams) CRC() uint32 {
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
		return nil, errors.Wrap(err, "sending InvokeWithLayer")
	}

	return data, nil
}

type InvokeWithTakeoutParams struct {
	TakeoutID int64
	Query     serialize.TLEncoder
}

func (*InvokeWithTakeoutParams) CRC() uint32 {
	return 0xaca9fd2e
}

func (t *InvokeWithTakeoutParams) Encode() []byte {
	buf := serialize.NewEncoder()
	buf.PutUint(t.CRC())
	buf.PutLong(t.TakeoutID)
	buf.PutRawBytes(t.Query.Encode())
	return buf.Result()
}

func (t *InvokeWithTakeoutParams) DecodeFrom(d *serialize.Decoder) {
	panic("makes no sense")
}

func (m *Client) InvokeWithTakeout(takeoutID int, query serialize.TLEncoder) (serialize.TL, error) {
	data, err := m.MakeRequest(&InvokeWithTakeoutParams{
		TakeoutID: int64(takeoutID),
		Query:     query,
	})
	if err != nil {
		return nil, errors.Wrap(err, "sending InvokeWithLayer")
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

func (c *Client) InitConnection(params *InitConnectionParams) (serialize.TL, error) {
	data, err := c.MakeRequest(params)
	if err != nil {
		return nil, errors.Wrap(err, "sending InitConnection")
	}

	return data, nil
}
