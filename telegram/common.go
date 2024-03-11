// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package telegram

import (
	"context"
	"net"
	"reflect"
	"runtime"
	"strconv"

	"github.com/pkg/errors"

	"github.com/xelaj/mtproto"
	"github.com/xelaj/mtproto/internal/session"
)

type Client struct {
	*mtproto.MTProto
	config       *ClientConfig
	serverConfig *Config
}

type ClientConfig struct {
	SessionFile     string
	ServerHost      string
	PublicKeysFile  string
	DeviceModel     string
	SystemVersion   string
	AppVersion      string
	AppID           int
	AppHash         string
	InitWarnChannel bool
}

const (
	warnChannelDefaultCapacity = 100
)

func NewClient(c ClientConfig) (*Client, error) { //nolint: gocritic arg is not ptr cause we call
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

	m := mtproto.New(mtproto.Config{
		SessionStorage: session.NewCached(session.NewFromFile(c.SessionFile)),
		PublicKey:      publicKeys[0],
	})

	if c.InitWarnChannel {
		m.Warnings = make(chan error, warnChannelDefaultCapacity)
	}

	go func() {
		if err := m.Connect(context.Background(), c.ServerHost); err != nil {
			panic("connecting")
		}
	}()

	client := &Client{
		MTProto: m,
		config:  &c,
	}

	//client.AddCustomServerRequestHandler(client.handleSpecialRequests())

	resp, err := client.InvokeWithLayer(ApiVersion, &InitConnectionParams[*HelpGetConfigParams]{
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

		dcList[int(dc.ID)] = net.JoinHostPort(dc.IpAddress, strconv.Itoa(int(dc.Port)))
	}
	client.SetDCList(dcList)
	return client, nil
}

func (m *Client) IsSessionRegistred() (bool, error) {
	_, err := m.UsersGetFullUser(&InputUserSelf{})
	if err == nil {
		return true, nil
	}

	if e := new(mtproto.ErrResponseCode); errors.As(err, &errCode) {
		if errCode.Message == "AUTH_KEY_UNREGISTERED" {
			return false, nil
		}
		return false, err
	} else {
		return false, err
	}
}

/*
func (c *Client) handleSpecialRequests() func(any) bool {
	return func(i any) bool {
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
*/
//----------------------------------------------------------------------------
