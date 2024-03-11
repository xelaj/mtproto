// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package mtproto

import (
	"context"
	"errors"
	"io"
)

type null = struct{}

// this is unofficial information, but it is suspected that the list of data
// centers is ABSOLUTELY identical for all the applications. Nevertheless, any
// client MUST explicitly specify the list of DCs, for the sake of reliability.
// this list is only experimental and is not part of the protocol.
func defaultDCList() map[int]string {
	return map[int]string{
		1: "149.154.175.58:443",
		2: "149.154.167.50:443",
		3: "149.154.175.100:443",
		4: "149.154.167.91:443",
		5: "91.108.56.151:443",
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func omitEOFErr(err error) error {
	if errors.Is(err, io.EOF) {
		return nil
	}

	return err
}

func omitContextErr(err error) error {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil
	}

	return err
}
