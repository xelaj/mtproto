package telegram

import (
	"fmt"
	"runtime"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/xelaj/errs"
	dry "github.com/xelaj/go-dry"

	"github.com/xelaj/mtproto"
	"github.com/xelaj/mtproto/encoding/tl"
	"github.com/xelaj/mtproto/keys"
)

const ApiVersion = 117

type Client struct {
	*mtproto.MTProto
	config *ClientConfig
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

	// if !dry.PathIsWirtable(c.SessionFile) {
	// 	return nil, errs.Permission(c.SessionFile).Scope("write")
	// }

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
	fmt.Println("mtproto created")
	err = m.CreateConnection()
	if err != nil {
		return nil, errors.Wrap(err, "creating connection")
	}
	fmt.Println("connection created")
	client := &Client{
		MTProto: m,
		config:  &c,
	}

	client.AddCustomServerRequestHandler(client.handleSpecialRequests())
	fmt.Println("HelpGetCfgParams invoking...")
	config := new(Config)
	err = client.InvokeWithLayer(ApiVersion, &InitConnectionParams{
		ApiID:          int32(c.AppID),
		DeviceModel:    c.DeviceModel,
		SystemVersion:  c.SystemVersion,
		AppVersion:     c.AppVersion,
		SystemLangCode: "en", // can't be edited, cause docs says that a single possible parameter
		LangCode:       "en",
		Query:          &HelpGetConfigParams{},
	}, config)
	fmt.Println("HelpGetCfgParams done...")
	if err != nil {
		return nil, errors.Wrap(err, "getting server configs")
	}

	pp.Println(config)

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
	Query tl.Object
}

func (_ *InvokeWithLayerParams) CRC() uint32 { return 0xda9b0d0d }

func (m *Client) InvokeWithLayer(layer int, query tl.Object, resp interface{}) error {
	return m.MakeRequest(&InvokeWithLayerParams{
		Layer: int32(layer),
		Query: query,
	}, resp)
}

type InvokeWithTakeoutParams struct {
	TakeoutID int64
	Query     tl.Object
}

func (*InvokeWithTakeoutParams) CRC() uint32 { return 0xaca9fd2e }

func (m *Client) InvokeWithTakeout(takeoutID int, query tl.Object, resp interface{}) error {
	return m.MakeRequest(&InvokeWithTakeoutParams{
		TakeoutID: int64(takeoutID),
		Query:     query,
	}, resp)
}

type InitConnectionParams struct {
	ApiID          int32
	DeviceModel    string
	SystemVersion  string
	AppVersion     string
	SystemLangCode string
	LangPack       string
	LangCode       string
	Proxy          *InputClientProxy `tl:"flag:0"`
	Params         JSONValue         `tl:"flag:1"`
	Query          tl.Object
}

func (_ *InitConnectionParams) CRC() uint32 { return 0xc1cd5ea9 }
