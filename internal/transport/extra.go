// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package transport

type any = interface{}
type null = struct{}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
