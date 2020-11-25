package telegram

//! WARNING this file is temporary and will be use as test fixtures.

import "github.com/xelaj/mtproto/encoding/tl"

type AccountPassword struct {
	HasRecovery             bool            `tl:"flag:0,encoded_in_bitflags"`
	HasSecureValues         bool            `tl:"flag:1,encoded_in_bitflags"`
	HasPassword             bool            `tl:"flag:2,encoded_in_bitflags"`
	CurrentAlgo             PasswordKdfAlgo `tl:"flag:2"`
	SrpB                    []byte          `tl:"flag:2"`
	SrpId                   int64           `tl:"flag:2"`
	Hint                    string          `tl:"flag:3"`
	EmailUnconfirmedPattern string          `tl:"flag:4"`
	NewAlgo                 PasswordKdfAlgo
	NewSecureAlgo           SecurePasswordKdfAlgo
	SecureRandom            []byte
}

func (*AccountPassword) CRC() uint32 {
	return uint32(0xad2641f8)
}

func (*AccountPassword) FlagIndex() int {
	return 0
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////

type PasswordKdfAlgo interface {
	tl.Object
	ImplementsPasswordKdfAlgo()
}

type PasswordKdfAlgoUnknown struct{}

func (*PasswordKdfAlgoUnknown) CRC() uint32 {
	return uint32(0xd45ab096)
}

func (*PasswordKdfAlgoUnknown) ImplementsPasswordKdfAlgo() {}

type PasswordKdfAlgoSHA256SHA256PBKDF2HMACSHA512iter100000SHA256ModPow struct {
	Salt1 []byte
	Salt2 []byte
	G     int32
	P     []byte
}

func (*PasswordKdfAlgoSHA256SHA256PBKDF2HMACSHA512iter100000SHA256ModPow) CRC() uint32 {
	return uint32(0x3a912d4a)
}

func (*PasswordKdfAlgoSHA256SHA256PBKDF2HMACSHA512iter100000SHA256ModPow) ImplementsPasswordKdfAlgo() {
}

type SecurePasswordKdfAlgo interface {
	tl.Object
	ImplementsSecurePasswordKdfAlgo()
}

type SecurePasswordKdfAlgoUnknown struct{}

func (*SecurePasswordKdfAlgoUnknown) CRC() uint32 {
	return uint32(0x4a8537)
}

func (*SecurePasswordKdfAlgoUnknown) ImplementsSecurePasswordKdfAlgo() {}

type SecurePasswordKdfAlgoPBKDF2HMACSHA512iter100000 struct {
	Salt []byte
}

func (*SecurePasswordKdfAlgoPBKDF2HMACSHA512iter100000) CRC() uint32 {
	return uint32(0xbbf2dda0)
}

func (*SecurePasswordKdfAlgoPBKDF2HMACSHA512iter100000) ImplementsSecurePasswordKdfAlgo() {}

type SecurePasswordKdfAlgoSHA512 struct {
	Salt []byte
}

func (*SecurePasswordKdfAlgoSHA512) CRC() uint32 {
	return uint32(0x86471d92)
}

func (*SecurePasswordKdfAlgoSHA512) ImplementsSecurePasswordKdfAlgo() {}

type InputCheckPasswordSRP interface {
	tl.Object
	ImplementsInputCheckPasswordSRP()
}

type InputCheckPasswordEmpty struct{}

func (*InputCheckPasswordEmpty) CRC() uint32 {
	return uint32(0x9880f658)
}

func (*InputCheckPasswordEmpty) ImplementsInputCheckPasswordSRP() {}

type InputCheckPasswordSRPObj struct {
	SrpId int64
	A     []byte
	M1    []byte
}

func (*InputCheckPasswordSRPObj) CRC() uint32 {
	return uint32(0xd27ff082)
}

func (*InputCheckPasswordSRPObj) ImplementsInputCheckPasswordSRP() {}

type Config struct {
	PhonecallsEnabled       bool `tl:"flag:1,encoded_in_bitflags"`
	DefaultP2PContacts      bool `tl:"flag:3,encoded_in_bitflags"`
	PreloadFeaturedStickers bool `tl:"flag:4,encoded_in_bitflags"`
	IgnorePhoneEntities     bool `tl:"flag:5,encoded_in_bitflags"`
	RevokePmInbox           bool `tl:"flag:6,encoded_in_bitflags"`
	BlockedMode             bool `tl:"flag:8,encoded_in_bitflags"`
	PfsEnabled              bool `tl:"flag:13,encoded_in_bitflags"`
	Date                    int32
	Expires                 int32
	TestMode                bool
	ThisDc                  int32
	DcOptions               []*DcOption
	DcTxtDomainName         string
	ChatSizeMax             int32
	MegagroupSizeMax        int32
	ForwardedCountMax       int32
	OnlineUpdatePeriodMs    int32
	OfflineBlurTimeoutMs    int32
	OfflineIdleTimeoutMs    int32
	OnlineCloudTimeoutMs    int32
	NotifyCloudDelayMs      int32
	NotifyDefaultDelayMs    int32
	PushChatPeriodMs        int32
	PushChatLimit           int32
	SavedGifsLimit          int32
	EditTimeLimit           int32
	RevokeTimeLimit         int32
	RevokePmTimeLimit       int32
	RatingEDecay            int32
	StickersRecentLimit     int32
	StickersFavedLimit      int32
	ChannelsReadMediaPeriod int32
	TmpSessions             int32 `tl:"flag:0"`
	PinnedDialogsCountMax   int32
	PinnedInfolderCountMax  int32
	CallReceiveTimeoutMs    int32
	CallRingTimeoutMs       int32
	CallConnectTimeoutMs    int32
	CallPacketTimeoutMs     int32
	MeUrlPrefix             string
	AutoupdateUrlPrefix     string `tl:"flag:7"`
	GifSearchUsername       string `tl:"flag:9"`
	VenueSearchUsername     string `tl:"flag:10"`
	ImgSearchUsername       string `tl:"flag:11"`
	StaticMapsProvider      string `tl:"flag:12"`
	CaptionLengthMax        int32
	MessageLengthMax        int32
	WebfileDcId             int32
	SuggestedLangCode       string `tl:"flag:2"`
	LangPackVersion         int32  `tl:"flag:2"`
	BaseLangPackVersion     int32  `tl:"flag:2"`
}

func (e *Config) CRC() uint32 {
	return uint32(0x330b4067)
}

func (*Config) FlagIndex() int {
	return 0
}

type DcOption struct {
	Ipv6      bool `tl:"flag:0,encoded_in_bitflags"`
	MediaOnly bool `tl:"flag:1,encoded_in_bitflags"`
	TcpoOnly  bool `tl:"flag:2,encoded_in_bitflags"`
	Cdn       bool `tl:"flag:3,encoded_in_bitflags"`
	Static    bool `tl:"flag:4,encoded_in_bitflags"`
	Id        int32
	IpAddress string
	Port      int32
	Secret    []byte `tl:"flag:10"`
}

func (e *DcOption) CRC() uint32 {
	return uint32(0x18b7a10d)
}

func (*DcOption) FlagIndex() int {
	return 0
}

type InputClientProxy struct {
	Address string
	Port    int32
}

func (*InputClientProxy) CRC() uint32 {
	return uint32(0x75588b3f)
}

type JSONValue interface {
	tl.Object
	ImplementsJSONValue()
}

type JsonNull struct{}

func (*JsonNull) CRC() uint32 {
	return uint32(0x3f6d7b68)
}

func (*JsonNull) ImplementsJSONValue() {}

type JsonBool struct {
	Value bool
}

func (*JsonBool) CRC() uint32 {
	return uint32(0xc7345e6a)
}

func (*JsonBool) ImplementsJSONValue() {}

type JsonNumber struct {
	Value float64
}

func (*JsonNumber) CRC() uint32 {
	return uint32(0x2be0dfa4)
}

func (*JsonNumber) ImplementsJSONValue() {}

type JsonString struct {
	Value string
}

func (*JsonString) CRC() uint32 {
	return uint32(0xb71e767a)
}

func (*JsonString) ImplementsJSONValue() {}

type JsonArray struct {
	Value []JSONValue
}

func (*JsonArray) CRC() uint32 {
	return uint32(0xf7444763)
}

func (*JsonArray) ImplementsJSONValue() {}

type JsonObject struct {
	Value []*JSONObjectValue
}

func (*JsonObject) CRC() uint32 {
	return uint32(0x99c1d49d)
}

func (*JsonObject) ImplementsJSONValue() {}

type JSONObjectValue struct {
	Key   string
	Value JSONValue
}

func (e *JSONObjectValue) CRC() uint32 {
	return uint32(0xc0de1bd9)
}
