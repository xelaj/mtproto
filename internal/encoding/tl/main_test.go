// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package tl_test

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/xelaj/go-dry"

	"github.com/xelaj/mtproto/internal/encoding/tl"
)

func tearup() {
	tl.RegisterObjects(
		&MultipleChats{},
		&Chat{},
		&AuthSentCode{},
		&SomeNullStruct{},
		&AuthSentCodeTypeApp{},
		&Rights{},
		&PollResults{},
		&PollAnswerVoters{},
		&AccountInstallThemeParams{},
		&InputThemeObj{},
		&AccountUnregisterDeviceParams{},
		&InvokeWithLayerParams{},
		&InitConnectionParams{},
		&ResPQ{},
		&AnyStructWithAnyType{},
		&AnyStructWithAnyObject{},
	)

	tl.RegisterEnums(
		AuthCodeTypeSms,
		AuthCodeTypeCall,
		AuthCodeTypeFlashCall,
	)
}

func teardown() {

}

func TestMain(m *testing.M) {
	tearup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func Hexed(in string) []byte {
	res, err := hex.DecodeString(in)
	dry.PanicIfErr(err)
	return res
}
