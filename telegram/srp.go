// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package telegram

import (
	"github.com/pkg/errors"

	"github.com/xelaj/mtproto/telegram/internal/srp"
)

func GetInputCheckPassword(password string, accountPassword *AccountPassword) (InputCheckPasswordSRP, error) {
	// У CurrentAlgo должен быть этот самый тип, с длинным названием алгоритма
	// https://github.com/tdlib/td/blob/f9009cbc01e9c4c77d31120a61feb9c639c6aeda/td/telegram/AuthManager.cpp#L537
	alg := accountPassword.CurrentAlgo
	current, ok := alg.(*PasswordKdfAlgoSHA256SHA256PBKDF2HMACSHA512iter100000SHA256ModPow)
	if !ok {
		return nil, errors.New("invalid CurrentAlgo type")
	}

	mp := &srp.ModPow{
		Salt1: current.Salt1,
		Salt2: current.Salt2,
		G:     current.G,
		P:     current.P,
	}

	res, err := srp.GetInputCheckPassword(password, accountPassword.SRPB, mp)
	if err != nil {
		return nil, errors.Wrap(err, "processing password")
	}

	if res == nil {
		return &InputCheckPasswordEmpty{}, nil
	}

	return &InputCheckPasswordSRPObj{
		SRPID: accountPassword.SRPID,
		A:     res.GA,
		M1:    res.M1,
	}, nil
}
