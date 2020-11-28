// Copyright (c) 2020 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package tl_test

// это набор разнообразных фикстур, которые хранятся исключительно для тестов
// фикстуры являются полноценно сгенерированными объектами из пакета telegram

import (
	"github.com/xelaj/mtproto/encoding/tl"
)

type MultipleChats struct {
	Chats []any
}

func (*MultipleChats) CRC() uint32 {
	return uint32(0xff1144cc)
}

type Chat struct {
	Creator           bool `tl:"flag:0,encoded_in_bitflags"`
	Kicked            bool `tl:"flag:1,encoded_in_bitflags"`
	Left              bool `tl:"flag:2,encoded_in_bitflags"`
	Deactivated       bool `tl:"flag:5,encoded_in_bitflags"`
	ID                int32
	Title             string
	Photo             string
	ParticipantsCount int32
	Date              int32
	Version           int32
	AdminRights       *Rights `tl:"flag:14"`
	BannedRights      *Rights `tl:"flag:18"`
}

func (*Chat) CRC() uint32 {
	return uint32(0x3bda1bde)
}

func (*Chat) FlagIndex() int {
	return 0
}

type AuthSentCode struct {
	Type          AuthSentCodeType
	PhoneCodeHash string
	NextType      AuthCodeType `tl:"flag:1"`
	Timeout       int32        `tl:"flag:2"`
}

func (*AuthSentCode) CRC() uint32 {
	return uint32(0x5e002502)
}

func (*AuthSentCode) FlagIndex() int {
	return 0
}

type SomeNullStruct struct{}

func (*SomeNullStruct) CRC() uint32 {
	return uint32(0xc4f9186b)
}

type AuthSentCodeType interface {
	tl.Object
	ImplementsAuthSentCodeType()
}

type AuthSentCodeTypeApp struct {
	Length int32
}

func (*AuthSentCodeTypeApp) CRC() uint32 {
	return uint32(0x3dbb5986)
}

func (*AuthSentCodeTypeApp) ImplementsAuthSentCodeType() {}

type Rights struct {
	DeleteMessages bool `tl:"flag:3,encoded_in_bitflags"`
	BanUsers       bool `tl:"flag:4,encoded_in_bitflags"`
}

func (*Rights) CRC() uint32 {
	return uint32(0x5fb224d5)
}

func (*Rights) FlagIndex() int {
	return 0
}

type AuthCodeType uint32

const (
	AuthCodeTypeSms       AuthCodeType = 1923290508
	AuthCodeTypeCall      AuthCodeType = 1948046307
	AuthCodeTypeFlashCall AuthCodeType = 577556219
)

func (e AuthCodeType) String() string {
	switch e {
	case AuthCodeTypeSms:
		return "auth.codeTypeSms"
	case AuthCodeTypeCall:
		return "auth.codeTypeCall"
	case AuthCodeTypeFlashCall:
		return "auth.codeTypeFlashCall"
	default:
		return "<UNKNOWN auth.CodeType>"
	}
}
func (e AuthCodeType) CRC() uint32 {
	return uint32(e)
}

type PollResults struct { //nolint:maligned required ordering
	Min              bool                `tl:"flag:0,encoded_in_bitflags"`
	Results          []*PollAnswerVoters `tl:"flag:1"`
	TotalVoters      int32               `tl:"flag:2"`
	RecentVoters     []int32             `tl:"flag:3"`
	Solution         string              `tl:"flag:4"`
	SolutionEntities []MessageEntity     `tl:"flag:4"`
}

func (*PollResults) CRC() uint32 {
	return uint32(0xbadcc1a3)
}
func (*PollResults) FlagIndex() int {
	return 0
}

type PollAnswerVoters struct { //nolint:maligned required ordering
	Chosen  bool `tl:"flag:0,encoded_in_bitflags"`
	Correct bool `tl:"flag:1,encoded_in_bitflags"`
	Option  []byte
	Voters  int32
}

func (*PollAnswerVoters) CRC() uint32 {
	return uint32(0x3b6ddad2)
}
func (*PollAnswerVoters) FlagIndex() int {
	return 0
}

type MessageEntity interface {
	tl.Object
	ImplementsMessageEntity()
}

// type MessageEntityUnknown struct {
// 	Offset int32
// 	Length int32
// }
//
// func (*MessageEntityUnknown) CRC() uint32 {
// 	return uint32(0xbb92ba95)
// }
//
// func (*MessageEntityUnknown) ImplementsMessageEntity() {}

type AccountInstallThemeParams struct {
	Dark   bool       `tl:"flag:0,encoded_in_bitflags"`
	Format string     `tl:"flag:1"`
	Theme  InputTheme `tl:"flag:1"`
}

func (e *AccountInstallThemeParams) CRC() uint32 {
	return uint32(0x7ae43737)
}

func (*AccountInstallThemeParams) FlagIndex() int {
	return 0
}

type InputTheme interface {
	tl.Object
	ImplementsInputTheme()
}

type InputThemeObj struct {
	ID         int64
	AccessHash int64
}

func (*InputThemeObj) CRC() uint32 {
	return uint32(0x3c5693e9)
}

func (*InputThemeObj) ImplementsInputTheme() {}

type AccountUnregisterDeviceParams struct {
	TokenType int32
	Token     string
	OtherUids []int32
}

func (e *AccountUnregisterDeviceParams) CRC() uint32 {
	return uint32(0x3076c4bf)
}

type InvokeWithLayerParams struct {
	Layer int32
	Query any
}

func (*InvokeWithLayerParams) CRC() uint32 {
	return 0xda9b0d0d
}

type InitConnectionParams struct {
	APIID          int32
	DeviceModel    string
	SystemVersion  string
	AppVersion     string
	SystemLangCode string
	LangPack       string
	LangCode       string
	Query          any
}

func (*InitConnectionParams) CRC() uint32 {
	return 0xc1cd5ea9
}

func (*InitConnectionParams) FlagIndex() int {
	return 0
}

type ResPQ struct {
	Nonce        *tl.Int128
	ServerNonce  *tl.Int128
	Pq           []byte
	Fingerprints []int64
}

func (*ResPQ) CRC() uint32 {
	return 0x05162463
}
