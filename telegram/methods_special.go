// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package telegram

// "Special" methods, cause they are wrappers over other queries. Cannot be
// generated, cause the the code gen does not understand what `!X` is (and for
// sure should not — 100%)

import (
	"github.com/pkg/errors"
	"github.com/xelaj/tl"
)

//invokeAfterMsg#cb9f372d {X:Type} msg_id:long query:!X = X;
//invokeAfterMsgs#3dc4b4f0 {X:Type} msg_ids:Vector<long> query:!X = X;

type InitConnectionParams[T tl.Object] struct {
	ApiID          int32             // Application identifier (see. App configuration)
	DeviceModel    string            // Device model
	SystemVersion  string            // Operation system version
	AppVersion     string            // Application version
	SystemLangCode string            // Code for the language used on the device's OS, ISO 639-1 standard
	LangPack       string            // Language pack to use
	LangCode       string            // Code for the language used on the client, ISO 639-1 standard
	Proxy          *InputClientProxy `tl:"flag:0"` // Info about an MTProto proxy
	Params         JsonValue         `tl:"flag:1"` // Additional initConnection parameters. For now, only the tz_offset field is supported, for specifying timezone offset in seconds.
	Query          T                 // The query itself
}

func (*InitConnectionParams[T]) CRC() uint32 { return 0xc1cd5ea9 }

type InvokeWithLayerParams[T tl.Object] struct {
	Layer int32
	Query T
}

func (*InvokeWithLayerParams[T]) CRC() uint32 { return 0xda9b0d0d }

func (m *Client) InvokeWithLayer(layer int, query tl.Object) (tl.Object, error) {
	data, err := m.MakeRequest(&InvokeWithLayerParams{
		Layer: int32(layer),
		Query: query,
	})
	if err != nil {
		return nil, errors.Wrap(err, "sending InvokeWithLayer")
	}

	return data.(tl.Object), nil
}

//invokeWithoutUpdates#bf9459b7 {X:Type} query:!X = X;
//invokeWithMessagesRange#365275f2 {X:Type} range:MessageRange query:!X = X;

func (m *Client) InvokeWithTakeout(takeoutID int, query tl.Object) (tl.Object, error) {
	data, err := m.MakeRequest(&InvokeWithTakeoutParams{
		TakeoutID: int64(takeoutID),
		Query:     query,
	})
	if err != nil {
		return nil, errors.Wrap(err, "sending InvokeWithLayer")
	}

	return data.(tl.Object), nil
}
