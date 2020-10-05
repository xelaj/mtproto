package telegram

import mtproto "github.com/xelaj/mtproto"

type BaseTheme uint32

const (
	BaseThemeClassic BaseTheme = 3282117730
	BaseThemeDay     BaseTheme = 4225242760
	BaseThemeNight   BaseTheme = 3081969320
	BaseThemeTinted  BaseTheme = 1834973166
	BaseThemeArctic  BaseTheme = 1527845466
)

func (e BaseTheme) String() string {
	switch e {
	case BaseTheme(3282117730):
		return "baseThemeClassic"
	case BaseTheme(4225242760):
		return "baseThemeDay"
	case BaseTheme(3081969320):
		return "baseThemeNight"
	case BaseTheme(1834973166):
		return "baseThemeTinted"
	case BaseTheme(1527845466):
		return "baseThemeArctic"
	default:
		return "<UNKNOWN BaseTheme>"
	}
}
func (e BaseTheme) CRC() uint32 {
	return uint32(e)
}
func (e BaseTheme) Encode() []byte {
	buf := mtproto.NewEncoder()
	buf.PutCRC(uint32(e))

	return buf.Result()
}

type SecureValueType uint32

const (
	SecureValueTypePersonalDetails       SecureValueType = 2636808675
	SecureValueTypePassport              SecureValueType = 1034709504
	SecureValueTypeDriverLicense         SecureValueType = 115615172
	SecureValueTypeIdentityCard          SecureValueType = 2698015819
	SecureValueTypeInternalPassport      SecureValueType = 2577698595
	SecureValueTypeAddress               SecureValueType = 3420659238
	SecureValueTypeUtilityBill           SecureValueType = 4231435598
	SecureValueTypeBankStatement         SecureValueType = 2299755533
	SecureValueTypeRentalAgreement       SecureValueType = 2340959368
	SecureValueTypePassportRegistration  SecureValueType = 2581823594
	SecureValueTypeTemporaryRegistration SecureValueType = 3926060083
	SecureValueTypePhone                 SecureValueType = 3005262555
	SecureValueTypeEmail                 SecureValueType = 2386339822
)

func (e SecureValueType) String() string {
	switch e {
	case SecureValueType(2636808675):
		return "secureValueTypePersonalDetails"
	case SecureValueType(1034709504):
		return "secureValueTypePassport"
	case SecureValueType(115615172):
		return "secureValueTypeDriverLicense"
	case SecureValueType(2698015819):
		return "secureValueTypeIdentityCard"
	case SecureValueType(2577698595):
		return "secureValueTypeInternalPassport"
	case SecureValueType(3420659238):
		return "secureValueTypeAddress"
	case SecureValueType(4231435598):
		return "secureValueTypeUtilityBill"
	case SecureValueType(2299755533):
		return "secureValueTypeBankStatement"
	case SecureValueType(2340959368):
		return "secureValueTypeRentalAgreement"
	case SecureValueType(2581823594):
		return "secureValueTypePassportRegistration"
	case SecureValueType(3926060083):
		return "secureValueTypeTemporaryRegistration"
	case SecureValueType(3005262555):
		return "secureValueTypePhone"
	case SecureValueType(2386339822):
		return "secureValueTypeEmail"
	default:
		return "<UNKNOWN SecureValueType>"
	}
}
func (e SecureValueType) CRC() uint32 {
	return uint32(e)
}
func (e SecureValueType) Encode() []byte {
	buf := mtproto.NewEncoder()
	buf.PutCRC(uint32(e))

	return buf.Result()
}

type PrivacyKey uint32

const (
	PrivacyKeyStatusTimestamp PrivacyKey = 3157175088
	PrivacyKeyChatInvite      PrivacyKey = 1343122938
	PrivacyKeyPhoneCall       PrivacyKey = 1030105979
	PrivacyKeyPhoneP2P        PrivacyKey = 961092808
	PrivacyKeyForwards        PrivacyKey = 1777096355
	PrivacyKeyProfilePhoto    PrivacyKey = 2517966829
	PrivacyKeyPhoneNumber     PrivacyKey = 3516589165
	PrivacyKeyAddedByPhone    PrivacyKey = 1124062251
)

func (e PrivacyKey) String() string {
	switch e {
	case PrivacyKey(3157175088):
		return "privacyKeyStatusTimestamp"
	case PrivacyKey(1343122938):
		return "privacyKeyChatInvite"
	case PrivacyKey(1030105979):
		return "privacyKeyPhoneCall"
	case PrivacyKey(961092808):
		return "privacyKeyPhoneP2P"
	case PrivacyKey(1777096355):
		return "privacyKeyForwards"
	case PrivacyKey(2517966829):
		return "privacyKeyProfilePhoto"
	case PrivacyKey(3516589165):
		return "privacyKeyPhoneNumber"
	case PrivacyKey(1124062251):
		return "privacyKeyAddedByPhone"
	default:
		return "<UNKNOWN PrivacyKey>"
	}
}
func (e PrivacyKey) CRC() uint32 {
	return uint32(e)
}
func (e PrivacyKey) Encode() []byte {
	buf := mtproto.NewEncoder()
	buf.PutCRC(uint32(e))

	return buf.Result()
}

type StorageFileType uint32

const (
	StorageFileUnknown StorageFileType = 2861972229
	StorageFilePartial StorageFileType = 1086091090
	StorageFileJpeg    StorageFileType = 8322574
	StorageFileGif     StorageFileType = 3403786975
	StorageFilePng     StorageFileType = 172975040
	StorageFilePdf     StorageFileType = 2921222285
	StorageFileMp3     StorageFileType = 1384777335
	StorageFileMov     StorageFileType = 1258941372
	StorageFileMp4     StorageFileType = 3016663268
	StorageFileWebp    StorageFileType = 276907596
)

func (e StorageFileType) String() string {
	switch e {
	case StorageFileType(2861972229):
		return "storage.fileUnknown"
	case StorageFileType(1086091090):
		return "storage.filePartial"
	case StorageFileType(8322574):
		return "storage.fileJpeg"
	case StorageFileType(3403786975):
		return "storage.fileGif"
	case StorageFileType(172975040):
		return "storage.filePng"
	case StorageFileType(2921222285):
		return "storage.filePdf"
	case StorageFileType(1384777335):
		return "storage.fileMp3"
	case StorageFileType(1258941372):
		return "storage.fileMov"
	case StorageFileType(3016663268):
		return "storage.fileMp4"
	case StorageFileType(276907596):
		return "storage.fileWebp"
	default:
		return "<UNKNOWN storage.FileType>"
	}
}
func (e StorageFileType) CRC() uint32 {
	return uint32(e)
}
func (e StorageFileType) Encode() []byte {
	buf := mtproto.NewEncoder()
	buf.PutCRC(uint32(e))

	return buf.Result()
}

type AuthCodeType uint32

const (
	AuthCodeTypeSms       AuthCodeType = 1923290508
	AuthCodeTypeCall      AuthCodeType = 1948046307
	AuthCodeTypeFlashCall AuthCodeType = 577556219
)

func (e AuthCodeType) String() string {
	switch e {
	case AuthCodeType(1923290508):
		return "auth.codeTypeSms"
	case AuthCodeType(1948046307):
		return "auth.codeTypeCall"
	case AuthCodeType(577556219):
		return "auth.codeTypeFlashCall"
	default:
		return "<UNKNOWN auth.CodeType>"
	}
}
func (e AuthCodeType) CRC() uint32 {
	return uint32(e)
}
func (e AuthCodeType) Encode() []byte {
	buf := mtproto.NewEncoder()
	buf.PutCRC(uint32(e))

	return buf.Result()
}

type TopPeerCategory uint32

const (
	TopPeerCategoryBotsPM         TopPeerCategory = 2875595611
	TopPeerCategoryBotsInline     TopPeerCategory = 344356834
	TopPeerCategoryCorrespondents TopPeerCategory = 104314861
	TopPeerCategoryGroups         TopPeerCategory = 3172442442
	TopPeerCategoryChannels       TopPeerCategory = 371037736
	TopPeerCategoryPhoneCalls     TopPeerCategory = 511092620
	TopPeerCategoryForwardUsers   TopPeerCategory = 2822794409
	TopPeerCategoryForwardChats   TopPeerCategory = 4226728176
)

func (e TopPeerCategory) String() string {
	switch e {
	case TopPeerCategory(2875595611):
		return "topPeerCategoryBotsPM"
	case TopPeerCategory(344356834):
		return "topPeerCategoryBotsInline"
	case TopPeerCategory(104314861):
		return "topPeerCategoryCorrespondents"
	case TopPeerCategory(3172442442):
		return "topPeerCategoryGroups"
	case TopPeerCategory(371037736):
		return "topPeerCategoryChannels"
	case TopPeerCategory(511092620):
		return "topPeerCategoryPhoneCalls"
	case TopPeerCategory(2822794409):
		return "topPeerCategoryForwardUsers"
	case TopPeerCategory(4226728176):
		return "topPeerCategoryForwardChats"
	default:
		return "<UNKNOWN TopPeerCategory>"
	}
}
func (e TopPeerCategory) CRC() uint32 {
	return uint32(e)
}
func (e TopPeerCategory) Encode() []byte {
	buf := mtproto.NewEncoder()
	buf.PutCRC(uint32(e))

	return buf.Result()
}

type InputPrivacyKey uint32

const (
	InputPrivacyKeyStatusTimestamp InputPrivacyKey = 1335282456
	InputPrivacyKeyChatInvite      InputPrivacyKey = 3187344422
	InputPrivacyKeyPhoneCall       InputPrivacyKey = 4206550111
	InputPrivacyKeyPhoneP2P        InputPrivacyKey = 3684593874
	InputPrivacyKeyForwards        InputPrivacyKey = 2765966344
	InputPrivacyKeyProfilePhoto    InputPrivacyKey = 1461304012
	InputPrivacyKeyPhoneNumber     InputPrivacyKey = 55761658
	InputPrivacyKeyAddedByPhone    InputPrivacyKey = 3508640733
)

func (e InputPrivacyKey) String() string {
	switch e {
	case InputPrivacyKey(1335282456):
		return "inputPrivacyKeyStatusTimestamp"
	case InputPrivacyKey(3187344422):
		return "inputPrivacyKeyChatInvite"
	case InputPrivacyKey(4206550111):
		return "inputPrivacyKeyPhoneCall"
	case InputPrivacyKey(3684593874):
		return "inputPrivacyKeyPhoneP2P"
	case InputPrivacyKey(2765966344):
		return "inputPrivacyKeyForwards"
	case InputPrivacyKey(1461304012):
		return "inputPrivacyKeyProfilePhoto"
	case InputPrivacyKey(55761658):
		return "inputPrivacyKeyPhoneNumber"
	case InputPrivacyKey(3508640733):
		return "inputPrivacyKeyAddedByPhone"
	default:
		return "<UNKNOWN InputPrivacyKey>"
	}
}
func (e InputPrivacyKey) CRC() uint32 {
	return uint32(e)
}
func (e InputPrivacyKey) Encode() []byte {
	buf := mtproto.NewEncoder()
	buf.PutCRC(uint32(e))

	return buf.Result()
}

type PhoneCallDiscardReason uint32

const (
	PhoneCallDiscardReasonMissed     PhoneCallDiscardReason = 2246320897
	PhoneCallDiscardReasonDisconnect PhoneCallDiscardReason = 3767910816
	PhoneCallDiscardReasonHangup     PhoneCallDiscardReason = 1471006352
	PhoneCallDiscardReasonBusy       PhoneCallDiscardReason = 4210550985
)

func (e PhoneCallDiscardReason) String() string {
	switch e {
	case PhoneCallDiscardReason(2246320897):
		return "phoneCallDiscardReasonMissed"
	case PhoneCallDiscardReason(3767910816):
		return "phoneCallDiscardReasonDisconnect"
	case PhoneCallDiscardReason(1471006352):
		return "phoneCallDiscardReasonHangup"
	case PhoneCallDiscardReason(4210550985):
		return "phoneCallDiscardReasonBusy"
	default:
		return "<UNKNOWN PhoneCallDiscardReason>"
	}
}
func (e PhoneCallDiscardReason) CRC() uint32 {
	return uint32(e)
}
func (e PhoneCallDiscardReason) Encode() []byte {
	buf := mtproto.NewEncoder()
	buf.PutCRC(uint32(e))

	return buf.Result()
}
