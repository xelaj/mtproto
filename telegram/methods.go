package telegram

import (
	"reflect"

	"github.com/pkg/errors"
)

type HelpGetConfigParams struct{}

func (e *HelpGetConfigParams) CRC() uint32 {
	return uint32(0xc4f9186b)
}

func (c *Client) HelpGetConfig() (*Config, error) {
	data, err := c.MakeRequest(&HelpGetConfigParams{})
	if err != nil {
		return nil, errors.Wrap(err, "sedning HelpGetConfig")
	}

	resp, ok := data.(*Config)
	if !ok {
		panic("got invalid response type: " + reflect.TypeOf(data).String())
	}

	return resp, nil
}
