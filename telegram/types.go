package telegram

import (
	"fmt"
	validator "github.com/go-playground/validator"
	zero "github.com/vikyd/zero"
	dry "github.com/xelaj/go-dry"
	mtproto "github.com/xelaj/mtproto"
	"reflect"
)

type PopularContact struct {
	ClientId  int64 `validate:"required"`
	Importers int32 `validate:"required"`
}

func (e *PopularContact) CRC() uint32 {
	return uint32(0x5ce14175)
}
func (e *PopularContact) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutLong(e.ClientId)
	buf.PutInt(e.Importers)
	return buf.Result()
}

func (e *PopularContact) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.ClientId = buf.PopLong()
	e.Importers = buf.PopInt()
}

type BotCommand struct {
	Command     string `validate:"required"`
	Description string `validate:"required"`
}

func (e *BotCommand) CRC() uint32 {
	return uint32(0xc27ac8c7)
}
func (e *BotCommand) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Command)
	buf.PutString(e.Description)
	return buf.Result()
}

func (e *BotCommand) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Command = buf.PopString()
	e.Description = buf.PopString()
}

type UploadWebFile struct {
	Size     int32           `validate:"required"`
	MimeType string          `validate:"required"`
	FileType StorageFileType `validate:"required"`
	Mtime    int32           `validate:"required"`
	Bytes    []byte          `validate:"required"`
}

func (e *UploadWebFile) CRC() uint32 {
	return uint32(0x21e753bc)
}
func (e *UploadWebFile) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.Size)
	buf.PutString(e.MimeType)
	buf.PutRawBytes(e.FileType.Encode())
	buf.PutInt(e.Mtime)
	buf.PutMessage(e.Bytes)
	return buf.Result()
}

func (e *UploadWebFile) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Size = buf.PopInt()
	e.MimeType = buf.PopString()
	e.FileType = *(buf.PopObj().(*StorageFileType))
	e.Mtime = buf.PopInt()
	e.Bytes = buf.PopMessage()
}

type StatsAbsValueAndPrev struct {
	Current  float64 `validate:"required"`
	Previous float64 `validate:"required"`
}

func (e *StatsAbsValueAndPrev) CRC() uint32 {
	return uint32(0xcb43acde)
}
func (e *StatsAbsValueAndPrev) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutDouble(e.Current)
	buf.PutDouble(e.Previous)
	return buf.Result()
}

func (e *StatsAbsValueAndPrev) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Current = buf.PopDouble()
	e.Previous = buf.PopDouble()
}

type MessagesMessageEditData struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	Caption         bool     `flag:"0,encoded_in_bitflags"`
}

func (e *MessagesMessageEditData) CRC() uint32 {
	return uint32(0x26b5dde6)
}
func (e *MessagesMessageEditData) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Caption) {
		flag |= 1 << 0
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Caption) {
	}
	return buf.Result()
}

func (e *MessagesMessageEditData) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Caption = buf.PopBool()
	}
}

type PostAddress struct {
	StreetLine1 string `validate:"required"`
	StreetLine2 string `validate:"required"`
	City        string `validate:"required"`
	State       string `validate:"required"`
	CountryIso2 string `validate:"required"`
	PostCode    string `validate:"required"`
}

func (e *PostAddress) CRC() uint32 {
	return uint32(0x1e8caaeb)
}
func (e *PostAddress) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.StreetLine1)
	buf.PutString(e.StreetLine2)
	buf.PutString(e.City)
	buf.PutString(e.State)
	buf.PutString(e.CountryIso2)
	buf.PutString(e.PostCode)
	return buf.Result()
}

func (e *PostAddress) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.StreetLine1 = buf.PopString()
	e.StreetLine2 = buf.PopString()
	e.City = buf.PopString()
	e.State = buf.PopString()
	e.CountryIso2 = buf.PopString()
	e.PostCode = buf.PopString()
}

type CodeSettings struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	AllowFlashcall  bool     `flag:"0,encoded_in_bitflags"`
	CurrentNumber   bool     `flag:"1,encoded_in_bitflags"`
	AllowAppHash    bool     `flag:"4,encoded_in_bitflags"`
}

func (e *CodeSettings) CRC() uint32 {
	return uint32(0xdebebe83)
}
func (e *CodeSettings) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.AllowFlashcall) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.CurrentNumber) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.AllowAppHash) {
		flag |= 1 << 4
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.AllowFlashcall) {
	}
	if !zero.IsZeroVal(e.CurrentNumber) {
	}
	if !zero.IsZeroVal(e.AllowAppHash) {
	}
	return buf.Result()
}

func (e *CodeSettings) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.AllowFlashcall = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.CurrentNumber = buf.PopBool()
	}
	if flags&1<<4 > 0 {
		e.AllowAppHash = buf.PopBool()
	}
}

type PhotosPhoto struct {
	Photo Photo  `validate:"required"`
	Users []User `validate:"required"`
}

func (e *PhotosPhoto) CRC() uint32 {
	return uint32(0x20212ca8)
}
func (e *PhotosPhoto) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Photo.Encode())
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *PhotosPhoto) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Photo = Photo(buf.PopObj())
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type Theme struct {
	__flagsPosition struct{}       // flags param position `validate:"required"`
	Creator         bool           `flag:"0,encoded_in_bitflags"`
	Default         bool           `flag:"1,encoded_in_bitflags"`
	Id              int64          `validate:"required"`
	AccessHash      int64          `validate:"required"`
	Slug            string         `validate:"required"`
	Title           string         `validate:"required"`
	Document        Document       `flag:"2"`
	Settings        *ThemeSettings `flag:"3"`
	InstallsCount   int32          `validate:"required"`
}

func (e *Theme) CRC() uint32 {
	return uint32(0x28f1114)
}
func (e *Theme) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Creator) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Default) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.Document) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.Settings) {
		flag |= 1 << 3
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Creator) {
	}
	if !zero.IsZeroVal(e.Default) {
	}
	buf.PutLong(e.Id)
	buf.PutLong(e.AccessHash)
	buf.PutString(e.Slug)
	buf.PutString(e.Title)
	if !zero.IsZeroVal(e.Document) {
		buf.PutRawBytes(e.Document.Encode())
	}
	if !zero.IsZeroVal(e.Settings) {
		buf.PutRawBytes(e.Settings.Encode())
	}
	buf.PutInt(e.InstallsCount)
	return buf.Result()
}

func (e *Theme) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Creator = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.Default = buf.PopBool()
	}
	e.Id = buf.PopLong()
	e.AccessHash = buf.PopLong()
	e.Slug = buf.PopString()
	e.Title = buf.PopString()
	if flags&1<<2 > 0 {
		e.Document = Document(buf.PopObj())
	}
	if flags&1<<3 > 0 {
		e.Settings = buf.PopObj().(*ThemeSettings)
	}
	e.InstallsCount = buf.PopInt()
}

type Error struct {
	Code int32  `validate:"required"`
	Text string `validate:"required"`
}

func (e *Error) CRC() uint32 {
	return uint32(0xc4b9f9bb)
}
func (e *Error) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.Code)
	buf.PutString(e.Text)
	return buf.Result()
}

func (e *Error) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Code = buf.PopInt()
	e.Text = buf.PopString()
}

type InputPhoneCall struct {
	Id         int64 `validate:"required"`
	AccessHash int64 `validate:"required"`
}

func (e *InputPhoneCall) CRC() uint32 {
	return uint32(0x1e36fded)
}
func (e *InputPhoneCall) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutLong(e.Id)
	buf.PutLong(e.AccessHash)
	return buf.Result()
}

func (e *InputPhoneCall) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Id = buf.PopLong()
	e.AccessHash = buf.PopLong()
}

type StatsMegagroupStats struct {
	Period                  *StatsDateRangeDays     `validate:"required"`
	Members                 *StatsAbsValueAndPrev   `validate:"required"`
	Messages                *StatsAbsValueAndPrev   `validate:"required"`
	Viewers                 *StatsAbsValueAndPrev   `validate:"required"`
	Posters                 *StatsAbsValueAndPrev   `validate:"required"`
	GrowthGraph             StatsGraph              `validate:"required"`
	MembersGraph            StatsGraph              `validate:"required"`
	NewMembersBySourceGraph StatsGraph              `validate:"required"`
	LanguagesGraph          StatsGraph              `validate:"required"`
	MessagesGraph           StatsGraph              `validate:"required"`
	ActionsGraph            StatsGraph              `validate:"required"`
	TopHoursGraph           StatsGraph              `validate:"required"`
	WeekdaysGraph           StatsGraph              `validate:"required"`
	TopPosters              []*StatsGroupTopPoster  `validate:"required"`
	TopAdmins               []*StatsGroupTopAdmin   `validate:"required"`
	TopInviters             []*StatsGroupTopInviter `validate:"required"`
	Users                   []User                  `validate:"required"`
}

func (e *StatsMegagroupStats) CRC() uint32 {
	return uint32(0xef7ff916)
}
func (e *StatsMegagroupStats) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Period.Encode())
	buf.PutRawBytes(e.Members.Encode())
	buf.PutRawBytes(e.Messages.Encode())
	buf.PutRawBytes(e.Viewers.Encode())
	buf.PutRawBytes(e.Posters.Encode())
	buf.PutRawBytes(e.GrowthGraph.Encode())
	buf.PutRawBytes(e.MembersGraph.Encode())
	buf.PutRawBytes(e.NewMembersBySourceGraph.Encode())
	buf.PutRawBytes(e.LanguagesGraph.Encode())
	buf.PutRawBytes(e.MessagesGraph.Encode())
	buf.PutRawBytes(e.ActionsGraph.Encode())
	buf.PutRawBytes(e.TopHoursGraph.Encode())
	buf.PutRawBytes(e.WeekdaysGraph.Encode())
	buf.PutVector(e.TopPosters)
	buf.PutVector(e.TopAdmins)
	buf.PutVector(e.TopInviters)
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *StatsMegagroupStats) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Period = buf.PopObj().(*StatsDateRangeDays)
	e.Members = buf.PopObj().(*StatsAbsValueAndPrev)
	e.Messages = buf.PopObj().(*StatsAbsValueAndPrev)
	e.Viewers = buf.PopObj().(*StatsAbsValueAndPrev)
	e.Posters = buf.PopObj().(*StatsAbsValueAndPrev)
	e.GrowthGraph = StatsGraph(buf.PopObj())
	e.MembersGraph = StatsGraph(buf.PopObj())
	e.NewMembersBySourceGraph = StatsGraph(buf.PopObj())
	e.LanguagesGraph = StatsGraph(buf.PopObj())
	e.MessagesGraph = StatsGraph(buf.PopObj())
	e.ActionsGraph = StatsGraph(buf.PopObj())
	e.TopHoursGraph = StatsGraph(buf.PopObj())
	e.WeekdaysGraph = StatsGraph(buf.PopObj())
	e.TopPosters = buf.PopVector(reflect.TypeOf(*StatsGroupTopPoster{})).([]*StatsGroupTopPoster)
	e.TopAdmins = buf.PopVector(reflect.TypeOf(*StatsGroupTopAdmin{})).([]*StatsGroupTopAdmin)
	e.TopInviters = buf.PopVector(reflect.TypeOf(*StatsGroupTopInviter{})).([]*StatsGroupTopInviter)
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type GlobalPrivacySettings struct {
	__flagsPosition                  struct{} // flags param position `validate:"required"`
	ArchiveAndMuteNewNoncontactPeers bool     `flag:"0"`
}

func (e *GlobalPrivacySettings) CRC() uint32 {
	return uint32(0xbea2f424)
}
func (e *GlobalPrivacySettings) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.ArchiveAndMuteNewNoncontactPeers) {
		flag |= 1 << 0
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.ArchiveAndMuteNewNoncontactPeers) {
		buf.PutBool(e.ArchiveAndMuteNewNoncontactPeers)
	}
	return buf.Result()
}

func (e *GlobalPrivacySettings) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.ArchiveAndMuteNewNoncontactPeers = buf.PopBool()
	}
}

type PaymentCharge struct {
	Id               string `validate:"required"`
	ProviderChargeId string `validate:"required"`
}

func (e *PaymentCharge) CRC() uint32 {
	return uint32(0xea02c27e)
}
func (e *PaymentCharge) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Id)
	buf.PutString(e.ProviderChargeId)
	return buf.Result()
}

func (e *PaymentCharge) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Id = buf.PopString()
	e.ProviderChargeId = buf.PopString()
}

type LangPackDifference struct {
	LangCode    string           `validate:"required"`
	FromVersion int32            `validate:"required"`
	Version     int32            `validate:"required"`
	Strings     []LangPackString `validate:"required"`
}

func (e *LangPackDifference) CRC() uint32 {
	return uint32(0xf385c1f6)
}
func (e *LangPackDifference) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.LangCode)
	buf.PutInt(e.FromVersion)
	buf.PutInt(e.Version)
	buf.PutVector(e.Strings)
	return buf.Result()
}

func (e *LangPackDifference) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.LangCode = buf.PopString()
	e.FromVersion = buf.PopInt()
	e.Version = buf.PopInt()
	e.Strings = buf.PopVector(reflect.TypeOf(LangPackString{})).([]LangPackString)
}

type InputEncryptedChat struct {
	ChatId     int32 `validate:"required"`
	AccessHash int64 `validate:"required"`
}

func (e *InputEncryptedChat) CRC() uint32 {
	return uint32(0xf141b5e1)
}
func (e *InputEncryptedChat) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.ChatId)
	buf.PutLong(e.AccessHash)
	return buf.Result()
}

func (e *InputEncryptedChat) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.ChatId = buf.PopInt()
	e.AccessHash = buf.PopLong()
}

type VideoSize struct {
	__flagsPosition struct{}      // flags param position `validate:"required"`
	Type            string        `validate:"required"`
	Location        *FileLocation `validate:"required"`
	W               int32         `validate:"required"`
	H               int32         `validate:"required"`
	Size            int32         `validate:"required"`
	VideoStartTs    float64       `flag:"0"`
}

func (e *VideoSize) CRC() uint32 {
	return uint32(0xe831c556)
}
func (e *VideoSize) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.VideoStartTs) {
		flag |= 1 << 0
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Type)
	buf.PutRawBytes(e.Location.Encode())
	buf.PutInt(e.W)
	buf.PutInt(e.H)
	buf.PutInt(e.Size)
	if !zero.IsZeroVal(e.VideoStartTs) {
		buf.PutDouble(e.VideoStartTs)
	}
	return buf.Result()
}

func (e *VideoSize) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	e.Type = buf.PopString()
	e.Location = buf.PopObj().(*FileLocation)
	e.W = buf.PopInt()
	e.H = buf.PopInt()
	e.Size = buf.PopInt()
	if flags&1<<0 > 0 {
		e.VideoStartTs = buf.PopDouble()
	}
}

type Poll struct {
	Id              int64         `validate:"required"`
	__flagsPosition struct{}      // flags param position `validate:"required"`
	Closed          bool          `flag:"0,encoded_in_bitflags"`
	PublicVoters    bool          `flag:"1,encoded_in_bitflags"`
	MultipleChoice  bool          `flag:"2,encoded_in_bitflags"`
	Quiz            bool          `flag:"3,encoded_in_bitflags"`
	Question        string        `validate:"required"`
	Answers         []*PollAnswer `validate:"required"`
	ClosePeriod     int32         `flag:"4"`
	CloseDate       int32         `flag:"5"`
}

func (e *Poll) CRC() uint32 {
	return uint32(0x86e18161)
}
func (e *Poll) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Closed) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.PublicVoters) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.MultipleChoice) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.Quiz) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.ClosePeriod) {
		flag |= 1 << 4
	}
	if !zero.IsZeroVal(e.CloseDate) {
		flag |= 1 << 5
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutLong(e.Id)
	if !zero.IsZeroVal(e.Closed) {
	}
	if !zero.IsZeroVal(e.PublicVoters) {
	}
	if !zero.IsZeroVal(e.MultipleChoice) {
	}
	if !zero.IsZeroVal(e.Quiz) {
	}
	buf.PutString(e.Question)
	buf.PutVector(e.Answers)
	if !zero.IsZeroVal(e.ClosePeriod) {
		buf.PutInt(e.ClosePeriod)
	}
	if !zero.IsZeroVal(e.CloseDate) {
		buf.PutInt(e.CloseDate)
	}
	return buf.Result()
}

func (e *Poll) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Id = buf.PopLong()
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Closed = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.PublicVoters = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.MultipleChoice = buf.PopBool()
	}
	if flags&1<<3 > 0 {
		e.Quiz = buf.PopBool()
	}
	e.Question = buf.PopString()
	e.Answers = buf.PopVector(reflect.TypeOf(*PollAnswer{})).([]*PollAnswer)
	if flags&1<<4 > 0 {
		e.ClosePeriod = buf.PopInt()
	}
	if flags&1<<5 > 0 {
		e.CloseDate = buf.PopInt()
	}
}

type ShippingOption struct {
	Id     string          `validate:"required"`
	Title  string          `validate:"required"`
	Prices []*LabeledPrice `validate:"required"`
}

func (e *ShippingOption) CRC() uint32 {
	return uint32(0xb6213cdf)
}
func (e *ShippingOption) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Id)
	buf.PutString(e.Title)
	buf.PutVector(e.Prices)
	return buf.Result()
}

func (e *ShippingOption) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Id = buf.PopString()
	e.Title = buf.PopString()
	e.Prices = buf.PopVector(reflect.TypeOf(*LabeledPrice{})).([]*LabeledPrice)
}

type StatsGroupTopPoster struct {
	UserId   int32 `validate:"required"`
	Messages int32 `validate:"required"`
	AvgChars int32 `validate:"required"`
}

func (e *StatsGroupTopPoster) CRC() uint32 {
	return uint32(0x18f3d0f7)
}
func (e *StatsGroupTopPoster) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.UserId)
	buf.PutInt(e.Messages)
	buf.PutInt(e.AvgChars)
	return buf.Result()
}

func (e *StatsGroupTopPoster) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.UserId = buf.PopInt()
	e.Messages = buf.PopInt()
	e.AvgChars = buf.PopInt()
}

type Authorization struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	Current         bool     `flag:"0,encoded_in_bitflags"`
	OfficialApp     bool     `flag:"1,encoded_in_bitflags"`
	PasswordPending bool     `flag:"2,encoded_in_bitflags"`
	Hash            int64    `validate:"required"`
	DeviceModel     string   `validate:"required"`
	Platform        string   `validate:"required"`
	SystemVersion   string   `validate:"required"`
	ApiId           int32    `validate:"required"`
	AppName         string   `validate:"required"`
	AppVersion      string   `validate:"required"`
	DateCreated     int32    `validate:"required"`
	DateActive      int32    `validate:"required"`
	Ip              string   `validate:"required"`
	Country         string   `validate:"required"`
	Region          string   `validate:"required"`
}

func (e *Authorization) CRC() uint32 {
	return uint32(0xad01d61d)
}
func (e *Authorization) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Current) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.OfficialApp) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.PasswordPending) {
		flag |= 1 << 2
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Current) {
	}
	if !zero.IsZeroVal(e.OfficialApp) {
	}
	if !zero.IsZeroVal(e.PasswordPending) {
	}
	buf.PutLong(e.Hash)
	buf.PutString(e.DeviceModel)
	buf.PutString(e.Platform)
	buf.PutString(e.SystemVersion)
	buf.PutInt(e.ApiId)
	buf.PutString(e.AppName)
	buf.PutString(e.AppVersion)
	buf.PutInt(e.DateCreated)
	buf.PutInt(e.DateActive)
	buf.PutString(e.Ip)
	buf.PutString(e.Country)
	buf.PutString(e.Region)
	return buf.Result()
}

func (e *Authorization) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Current = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.OfficialApp = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.PasswordPending = buf.PopBool()
	}
	e.Hash = buf.PopLong()
	e.DeviceModel = buf.PopString()
	e.Platform = buf.PopString()
	e.SystemVersion = buf.PopString()
	e.ApiId = buf.PopInt()
	e.AppName = buf.PopString()
	e.AppVersion = buf.PopString()
	e.DateCreated = buf.PopInt()
	e.DateActive = buf.PopInt()
	e.Ip = buf.PopString()
	e.Country = buf.PopString()
	e.Region = buf.PopString()
}

type PaymentsSavedInfo struct {
	__flagsPosition     struct{}              // flags param position `validate:"required"`
	HasSavedCredentials bool                  `flag:"1,encoded_in_bitflags"`
	SavedInfo           *PaymentRequestedInfo `flag:"0"`
}

func (e *PaymentsSavedInfo) CRC() uint32 {
	return uint32(0xfb8fe43c)
}
func (e *PaymentsSavedInfo) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.SavedInfo) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.HasSavedCredentials) {
		flag |= 1 << 1
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.HasSavedCredentials) {
	}
	if !zero.IsZeroVal(e.SavedInfo) {
		buf.PutRawBytes(e.SavedInfo.Encode())
	}
	return buf.Result()
}

func (e *PaymentsSavedInfo) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<1 > 0 {
		e.HasSavedCredentials = buf.PopBool()
	}
	if flags&1<<0 > 0 {
		e.SavedInfo = buf.PopObj().(*PaymentRequestedInfo)
	}
}

type CdnConfig struct {
	PublicKeys []*CdnPublicKey `validate:"required"`
}

func (e *CdnConfig) CRC() uint32 {
	return uint32(0x5725e40a)
}
func (e *CdnConfig) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutVector(e.PublicKeys)
	return buf.Result()
}

func (e *CdnConfig) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.PublicKeys = buf.PopVector(reflect.TypeOf(*CdnPublicKey{})).([]*CdnPublicKey)
}

type MaskCoords struct {
	N    int32   `validate:"required"`
	X    float64 `validate:"required"`
	Y    float64 `validate:"required"`
	Zoom float64 `validate:"required"`
}

func (e *MaskCoords) CRC() uint32 {
	return uint32(0xaed6dbb2)
}
func (e *MaskCoords) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.N)
	buf.PutDouble(e.X)
	buf.PutDouble(e.Y)
	buf.PutDouble(e.Zoom)
	return buf.Result()
}

func (e *MaskCoords) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.N = buf.PopInt()
	e.X = buf.PopDouble()
	e.Y = buf.PopDouble()
	e.Zoom = buf.PopDouble()
}

type MessagesChatFull struct {
	FullChat ChatFull `validate:"required"`
	Chats    []Chat   `validate:"required"`
	Users    []User   `validate:"required"`
}

func (e *MessagesChatFull) CRC() uint32 {
	return uint32(0xe5d7d19c)
}
func (e *MessagesChatFull) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.FullChat.Encode())
	buf.PutVector(e.Chats)
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *MessagesChatFull) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.FullChat = ChatFull(buf.PopObj())
	e.Chats = buf.PopVector(reflect.TypeOf(Chat{})).([]Chat)
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type InputBotInlineMessageID struct {
	DcId       int32 `validate:"required"`
	Id         int64 `validate:"required"`
	AccessHash int64 `validate:"required"`
}

func (e *InputBotInlineMessageID) CRC() uint32 {
	return uint32(0x890c3d89)
}
func (e *InputBotInlineMessageID) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.DcId)
	buf.PutLong(e.Id)
	buf.PutLong(e.AccessHash)
	return buf.Result()
}

func (e *InputBotInlineMessageID) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.DcId = buf.PopInt()
	e.Id = buf.PopLong()
	e.AccessHash = buf.PopLong()
}

type HelpRecentMeUrls struct {
	Urls  []RecentMeUrl `validate:"required"`
	Chats []Chat        `validate:"required"`
	Users []User        `validate:"required"`
}

func (e *HelpRecentMeUrls) CRC() uint32 {
	return uint32(0xe0310d7)
}
func (e *HelpRecentMeUrls) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutVector(e.Urls)
	buf.PutVector(e.Chats)
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *HelpRecentMeUrls) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Urls = buf.PopVector(reflect.TypeOf(RecentMeUrl{})).([]RecentMeUrl)
	e.Chats = buf.PopVector(reflect.TypeOf(Chat{})).([]Chat)
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type PeerNotifySettings struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	ShowPreviews    bool     `flag:"0"`
	Silent          bool     `flag:"1"`
	MuteUntil       int32    `flag:"2"`
	Sound           string   `flag:"3"`
}

func (e *PeerNotifySettings) CRC() uint32 {
	return uint32(0xaf509d20)
}
func (e *PeerNotifySettings) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.ShowPreviews) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Silent) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.MuteUntil) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.Sound) {
		flag |= 1 << 3
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.ShowPreviews) {
		buf.PutBool(e.ShowPreviews)
	}
	if !zero.IsZeroVal(e.Silent) {
		buf.PutBool(e.Silent)
	}
	if !zero.IsZeroVal(e.MuteUntil) {
		buf.PutInt(e.MuteUntil)
	}
	if !zero.IsZeroVal(e.Sound) {
		buf.PutString(e.Sound)
	}
	return buf.Result()
}

func (e *PeerNotifySettings) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.ShowPreviews = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.Silent = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.MuteUntil = buf.PopInt()
	}
	if flags&1<<3 > 0 {
		e.Sound = buf.PopString()
	}
}

type SecureCredentialsEncrypted struct {
	Data   []byte `validate:"required"`
	Hash   []byte `validate:"required"`
	Secret []byte `validate:"required"`
}

func (e *SecureCredentialsEncrypted) CRC() uint32 {
	return uint32(0x33f0ea47)
}
func (e *SecureCredentialsEncrypted) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutMessage(e.Data)
	buf.PutMessage(e.Hash)
	buf.PutMessage(e.Secret)
	return buf.Result()
}

func (e *SecureCredentialsEncrypted) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Data = buf.PopMessage()
	e.Hash = buf.PopMessage()
	e.Secret = buf.PopMessage()
}

type ImportedContact struct {
	UserId   int32 `validate:"required"`
	ClientId int64 `validate:"required"`
}

func (e *ImportedContact) CRC() uint32 {
	return uint32(0xd0028438)
}
func (e *ImportedContact) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.UserId)
	buf.PutLong(e.ClientId)
	return buf.Result()
}

func (e *ImportedContact) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.UserId = buf.PopInt()
	e.ClientId = buf.PopLong()
}

type EmojiKeywordsDifference struct {
	LangCode    string         `validate:"required"`
	FromVersion int32          `validate:"required"`
	Version     int32          `validate:"required"`
	Keywords    []EmojiKeyword `validate:"required"`
}

func (e *EmojiKeywordsDifference) CRC() uint32 {
	return uint32(0x5cc761bd)
}
func (e *EmojiKeywordsDifference) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.LangCode)
	buf.PutInt(e.FromVersion)
	buf.PutInt(e.Version)
	buf.PutVector(e.Keywords)
	return buf.Result()
}

func (e *EmojiKeywordsDifference) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.LangCode = buf.PopString()
	e.FromVersion = buf.PopInt()
	e.Version = buf.PopInt()
	e.Keywords = buf.PopVector(reflect.TypeOf(EmojiKeyword{})).([]EmojiKeyword)
}

type FileLocation struct {
	VolumeId int64 `validate:"required"`
	LocalId  int32 `validate:"required"`
}

func (e *FileLocation) CRC() uint32 {
	return uint32(0xbc7fc6cd)
}
func (e *FileLocation) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutLong(e.VolumeId)
	buf.PutInt(e.LocalId)
	return buf.Result()
}

func (e *FileLocation) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.VolumeId = buf.PopLong()
	e.LocalId = buf.PopInt()
}

type Game struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	Id              int64    `validate:"required"`
	AccessHash      int64    `validate:"required"`
	ShortName       string   `validate:"required"`
	Title           string   `validate:"required"`
	Description     string   `validate:"required"`
	Photo           Photo    `validate:"required"`
	Document        Document `flag:"0"`
}

func (e *Game) CRC() uint32 {
	return uint32(0xbdf9653b)
}
func (e *Game) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Document) {
		flag |= 1 << 0
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutLong(e.Id)
	buf.PutLong(e.AccessHash)
	buf.PutString(e.ShortName)
	buf.PutString(e.Title)
	buf.PutString(e.Description)
	buf.PutRawBytes(e.Photo.Encode())
	if !zero.IsZeroVal(e.Document) {
		buf.PutRawBytes(e.Document.Encode())
	}
	return buf.Result()
}

func (e *Game) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	e.Id = buf.PopLong()
	e.AccessHash = buf.PopLong()
	e.ShortName = buf.PopString()
	e.Title = buf.PopString()
	e.Description = buf.PopString()
	e.Photo = Photo(buf.PopObj())
	if flags&1<<0 > 0 {
		e.Document = Document(buf.PopObj())
	}
}

type InputStickerSetItem struct {
	__flagsPosition struct{}      // flags param position `validate:"required"`
	Document        InputDocument `validate:"required"`
	Emoji           string        `validate:"required"`
	MaskCoords      *MaskCoords   `flag:"0"`
}

func (e *InputStickerSetItem) CRC() uint32 {
	return uint32(0xffa0a496)
}
func (e *InputStickerSetItem) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.MaskCoords) {
		flag |= 1 << 0
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Document.Encode())
	buf.PutString(e.Emoji)
	if !zero.IsZeroVal(e.MaskCoords) {
		buf.PutRawBytes(e.MaskCoords.Encode())
	}
	return buf.Result()
}

func (e *InputStickerSetItem) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	e.Document = InputDocument(buf.PopObj())
	e.Emoji = buf.PopString()
	if flags&1<<0 > 0 {
		e.MaskCoords = buf.PopObj().(*MaskCoords)
	}
}

type ChannelsChannelParticipant struct {
	Participant ChannelParticipant `validate:"required"`
	Users       []User             `validate:"required"`
}

func (e *ChannelsChannelParticipant) CRC() uint32 {
	return uint32(0xd0d9b163)
}
func (e *ChannelsChannelParticipant) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Participant.Encode())
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *ChannelsChannelParticipant) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Participant = ChannelParticipant(buf.PopObj())
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type PaymentSavedCredentials struct {
	Id    string `validate:"required"`
	Title string `validate:"required"`
}

func (e *PaymentSavedCredentials) CRC() uint32 {
	return uint32(0xcdc27a1f)
}
func (e *PaymentSavedCredentials) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Id)
	buf.PutString(e.Title)
	return buf.Result()
}

func (e *PaymentSavedCredentials) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Id = buf.PopString()
	e.Title = buf.PopString()
}

type PaymentsPaymentForm struct {
	__flagsPosition    struct{}                 // flags param position `validate:"required"`
	CanSaveCredentials bool                     `flag:"2,encoded_in_bitflags"`
	PasswordMissing    bool                     `flag:"3,encoded_in_bitflags"`
	BotId              int32                    `validate:"required"`
	Invoice            *Invoice                 `validate:"required"`
	ProviderId         int32                    `validate:"required"`
	Url                string                   `validate:"required"`
	NativeProvider     string                   `flag:"4"`
	NativeParams       *DataJSON                `flag:"4"`
	SavedInfo          *PaymentRequestedInfo    `flag:"0"`
	SavedCredentials   *PaymentSavedCredentials `flag:"1"`
	Users              []User                   `validate:"required"`
}

func (e *PaymentsPaymentForm) CRC() uint32 {
	return uint32(0x3f56aea3)
}
func (e *PaymentsPaymentForm) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.SavedInfo) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.SavedCredentials) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.CanSaveCredentials) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.PasswordMissing) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.NativeProvider) || !zero.IsZeroVal(e.NativeParams) {
		flag |= 1 << 4
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.CanSaveCredentials) {
	}
	if !zero.IsZeroVal(e.PasswordMissing) {
	}
	buf.PutInt(e.BotId)
	buf.PutRawBytes(e.Invoice.Encode())
	buf.PutInt(e.ProviderId)
	buf.PutString(e.Url)
	if !zero.IsZeroVal(e.NativeProvider) {
		buf.PutString(e.NativeProvider)
	}
	if !zero.IsZeroVal(e.NativeParams) {
		buf.PutRawBytes(e.NativeParams.Encode())
	}
	if !zero.IsZeroVal(e.SavedInfo) {
		buf.PutRawBytes(e.SavedInfo.Encode())
	}
	if !zero.IsZeroVal(e.SavedCredentials) {
		buf.PutRawBytes(e.SavedCredentials.Encode())
	}
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *PaymentsPaymentForm) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<2 > 0 {
		e.CanSaveCredentials = buf.PopBool()
	}
	if flags&1<<3 > 0 {
		e.PasswordMissing = buf.PopBool()
	}
	e.BotId = buf.PopInt()
	e.Invoice = buf.PopObj().(*Invoice)
	e.ProviderId = buf.PopInt()
	e.Url = buf.PopString()
	if flags&1<<4 > 0 {
		e.NativeProvider = buf.PopString()
	}
	if flags&1<<4 > 0 {
		e.NativeParams = buf.PopObj().(*DataJSON)
	}
	if flags&1<<0 > 0 {
		e.SavedInfo = buf.PopObj().(*PaymentRequestedInfo)
	}
	if flags&1<<1 > 0 {
		e.SavedCredentials = buf.PopObj().(*PaymentSavedCredentials)
	}
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type ChannelsAdminLogResults struct {
	Events []*ChannelAdminLogEvent `validate:"required"`
	Chats  []Chat                  `validate:"required"`
	Users  []User                  `validate:"required"`
}

func (e *ChannelsAdminLogResults) CRC() uint32 {
	return uint32(0xed8af74d)
}
func (e *ChannelsAdminLogResults) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutVector(e.Events)
	buf.PutVector(e.Chats)
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *ChannelsAdminLogResults) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Events = buf.PopVector(reflect.TypeOf(*ChannelAdminLogEvent{})).([]*ChannelAdminLogEvent)
	e.Chats = buf.PopVector(reflect.TypeOf(Chat{})).([]Chat)
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type SavedContact struct {
	Phone     string `validate:"required"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Date      int32  `validate:"required"`
}

func (e *SavedContact) CRC() uint32 {
	return uint32(0x1142bd56)
}
func (e *SavedContact) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Phone)
	buf.PutString(e.FirstName)
	buf.PutString(e.LastName)
	buf.PutInt(e.Date)
	return buf.Result()
}

func (e *SavedContact) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Phone = buf.PopString()
	e.FirstName = buf.PopString()
	e.LastName = buf.PopString()
	e.Date = buf.PopInt()
}

type MessagesInactiveChats struct {
	Dates []int32 `validate:"required"`
	Chats []Chat  `validate:"required"`
	Users []User  `validate:"required"`
}

func (e *MessagesInactiveChats) CRC() uint32 {
	return uint32(0xa927fec5)
}
func (e *MessagesInactiveChats) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutVector(e.Dates)
	buf.PutVector(e.Chats)
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *MessagesInactiveChats) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Dates = buf.PopVector(reflect.TypeOf(int32{})).([]int32)
	e.Chats = buf.PopVector(reflect.TypeOf(Chat{})).([]Chat)
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type ThemeSettings struct {
	__flagsPosition    struct{}  // flags param position `validate:"required"`
	BaseTheme          BaseTheme `validate:"required"`
	AccentColor        int32     `validate:"required"`
	MessageTopColor    int32     `flag:"0"`
	MessageBottomColor int32     `flag:"0"`
	Wallpaper          WallPaper `flag:"1"`
}

func (e *ThemeSettings) CRC() uint32 {
	return uint32(0x9c14984a)
}
func (e *ThemeSettings) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.MessageTopColor) || !zero.IsZeroVal(e.MessageBottomColor) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Wallpaper) {
		flag |= 1 << 1
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.BaseTheme.Encode())
	buf.PutInt(e.AccentColor)
	if !zero.IsZeroVal(e.MessageTopColor) {
		buf.PutInt(e.MessageTopColor)
	}
	if !zero.IsZeroVal(e.MessageBottomColor) {
		buf.PutInt(e.MessageBottomColor)
	}
	if !zero.IsZeroVal(e.Wallpaper) {
		buf.PutRawBytes(e.Wallpaper.Encode())
	}
	return buf.Result()
}

func (e *ThemeSettings) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	e.BaseTheme = *(buf.PopObj().(*BaseTheme))
	e.AccentColor = buf.PopInt()
	if flags&1<<0 > 0 {
		e.MessageTopColor = buf.PopInt()
	}
	if flags&1<<0 > 0 {
		e.MessageBottomColor = buf.PopInt()
	}
	if flags&1<<1 > 0 {
		e.Wallpaper = WallPaper(buf.PopObj())
	}
}

type LabeledPrice struct {
	Label  string `validate:"required"`
	Amount int64  `validate:"required"`
}

func (e *LabeledPrice) CRC() uint32 {
	return uint32(0xcb296bf8)
}
func (e *LabeledPrice) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Label)
	buf.PutLong(e.Amount)
	return buf.Result()
}

func (e *LabeledPrice) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Label = buf.PopString()
	e.Amount = buf.PopLong()
}

type StickerPack struct {
	Emoticon  string  `validate:"required"`
	Documents []int64 `validate:"required"`
}

func (e *StickerPack) CRC() uint32 {
	return uint32(0x12b299d4)
}
func (e *StickerPack) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Emoticon)
	buf.PutVector(e.Documents)
	return buf.Result()
}

func (e *StickerPack) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Emoticon = buf.PopString()
	e.Documents = buf.PopVector(reflect.TypeOf(int64{})).([]int64)
}

type PageTableRow struct {
	Cells []*PageTableCell `validate:"required"`
}

func (e *PageTableRow) CRC() uint32 {
	return uint32(0xe0c0c5e5)
}
func (e *PageTableRow) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutVector(e.Cells)
	return buf.Result()
}

func (e *PageTableRow) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Cells = buf.PopVector(reflect.TypeOf(*PageTableCell{})).([]*PageTableCell)
}

type FolderPeer struct {
	Peer     Peer  `validate:"required"`
	FolderId int32 `validate:"required"`
}

func (e *FolderPeer) CRC() uint32 {
	return uint32(0xe9baa668)
}
func (e *FolderPeer) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Peer.Encode())
	buf.PutInt(e.FolderId)
	return buf.Result()
}

func (e *FolderPeer) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Peer = Peer(buf.PopObj())
	e.FolderId = buf.PopInt()
}

type AccountContentSettings struct {
	__flagsPosition    struct{} // flags param position `validate:"required"`
	SensitiveEnabled   bool     `flag:"0,encoded_in_bitflags"`
	SensitiveCanChange bool     `flag:"1,encoded_in_bitflags"`
}

func (e *AccountContentSettings) CRC() uint32 {
	return uint32(0x57e28221)
}
func (e *AccountContentSettings) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.SensitiveEnabled) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.SensitiveCanChange) {
		flag |= 1 << 1
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.SensitiveEnabled) {
	}
	if !zero.IsZeroVal(e.SensitiveCanChange) {
	}
	return buf.Result()
}

func (e *AccountContentSettings) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.SensitiveEnabled = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.SensitiveCanChange = buf.PopBool()
	}
}

type AccountTakeout struct {
	Id int64 `validate:"required"`
}

func (e *AccountTakeout) CRC() uint32 {
	return uint32(0x4dba4501)
}
func (e *AccountTakeout) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutLong(e.Id)
	return buf.Result()
}

func (e *AccountTakeout) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Id = buf.PopLong()
}

type StatsGroupTopAdmin struct {
	UserId  int32 `validate:"required"`
	Deleted int32 `validate:"required"`
	Kicked  int32 `validate:"required"`
	Banned  int32 `validate:"required"`
}

func (e *StatsGroupTopAdmin) CRC() uint32 {
	return uint32(0x6014f412)
}
func (e *StatsGroupTopAdmin) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.UserId)
	buf.PutInt(e.Deleted)
	buf.PutInt(e.Kicked)
	buf.PutInt(e.Banned)
	return buf.Result()
}

func (e *StatsGroupTopAdmin) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.UserId = buf.PopInt()
	e.Deleted = buf.PopInt()
	e.Kicked = buf.PopInt()
	e.Banned = buf.PopInt()
}

type MessagesBotCallbackAnswer struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	Alert           bool     `flag:"1,encoded_in_bitflags"`
	HasUrl          bool     `flag:"3,encoded_in_bitflags"`
	NativeUi        bool     `flag:"4,encoded_in_bitflags"`
	Message         string   `flag:"0"`
	Url             string   `flag:"2"`
	CacheTime       int32    `validate:"required"`
}

func (e *MessagesBotCallbackAnswer) CRC() uint32 {
	return uint32(0x36585ea4)
}
func (e *MessagesBotCallbackAnswer) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Message) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Alert) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.Url) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.HasUrl) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.NativeUi) {
		flag |= 1 << 4
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Alert) {
	}
	if !zero.IsZeroVal(e.HasUrl) {
	}
	if !zero.IsZeroVal(e.NativeUi) {
	}
	if !zero.IsZeroVal(e.Message) {
		buf.PutString(e.Message)
	}
	if !zero.IsZeroVal(e.Url) {
		buf.PutString(e.Url)
	}
	buf.PutInt(e.CacheTime)
	return buf.Result()
}

func (e *MessagesBotCallbackAnswer) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<1 > 0 {
		e.Alert = buf.PopBool()
	}
	if flags&1<<3 > 0 {
		e.HasUrl = buf.PopBool()
	}
	if flags&1<<4 > 0 {
		e.NativeUi = buf.PopBool()
	}
	if flags&1<<0 > 0 {
		e.Message = buf.PopString()
	}
	if flags&1<<2 > 0 {
		e.Url = buf.PopString()
	}
	e.CacheTime = buf.PopInt()
}

type StatsPercentValue struct {
	Part  float64 `validate:"required"`
	Total float64 `validate:"required"`
}

func (e *StatsPercentValue) CRC() uint32 {
	return uint32(0xcbce2fe0)
}
func (e *StatsPercentValue) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutDouble(e.Part)
	buf.PutDouble(e.Total)
	return buf.Result()
}

func (e *StatsPercentValue) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Part = buf.PopDouble()
	e.Total = buf.PopDouble()
}

type AccountPrivacyRules struct {
	Rules []PrivacyRule `validate:"required"`
	Chats []Chat        `validate:"required"`
	Users []User        `validate:"required"`
}

func (e *AccountPrivacyRules) CRC() uint32 {
	return uint32(0x50a04e45)
}
func (e *AccountPrivacyRules) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutVector(e.Rules)
	buf.PutVector(e.Chats)
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *AccountPrivacyRules) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Rules = buf.PopVector(reflect.TypeOf(PrivacyRule{})).([]PrivacyRule)
	e.Chats = buf.PopVector(reflect.TypeOf(Chat{})).([]Chat)
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type MessageFwdHeader struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	FromId          int32    `flag:"0"`
	FromName        string   `flag:"5"`
	Date            int32    `validate:"required"`
	ChannelId       int32    `flag:"1"`
	ChannelPost     int32    `flag:"2"`
	PostAuthor      string   `flag:"3"`
	SavedFromPeer   Peer     `flag:"4"`
	SavedFromMsgId  int32    `flag:"4"`
	PsaType         string   `flag:"6"`
}

func (e *MessageFwdHeader) CRC() uint32 {
	return uint32(0x353a686b)
}
func (e *MessageFwdHeader) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.FromId) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.ChannelId) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.ChannelPost) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.PostAuthor) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.SavedFromPeer) || !zero.IsZeroVal(e.SavedFromMsgId) {
		flag |= 1 << 4
	}
	if !zero.IsZeroVal(e.FromName) {
		flag |= 1 << 5
	}
	if !zero.IsZeroVal(e.PsaType) {
		flag |= 1 << 6
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.FromId) {
		buf.PutInt(e.FromId)
	}
	if !zero.IsZeroVal(e.FromName) {
		buf.PutString(e.FromName)
	}
	buf.PutInt(e.Date)
	if !zero.IsZeroVal(e.ChannelId) {
		buf.PutInt(e.ChannelId)
	}
	if !zero.IsZeroVal(e.ChannelPost) {
		buf.PutInt(e.ChannelPost)
	}
	if !zero.IsZeroVal(e.PostAuthor) {
		buf.PutString(e.PostAuthor)
	}
	if !zero.IsZeroVal(e.SavedFromPeer) {
		buf.PutRawBytes(e.SavedFromPeer.Encode())
	}
	if !zero.IsZeroVal(e.SavedFromMsgId) {
		buf.PutInt(e.SavedFromMsgId)
	}
	if !zero.IsZeroVal(e.PsaType) {
		buf.PutString(e.PsaType)
	}
	return buf.Result()
}

func (e *MessageFwdHeader) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.FromId = buf.PopInt()
	}
	if flags&1<<5 > 0 {
		e.FromName = buf.PopString()
	}
	e.Date = buf.PopInt()
	if flags&1<<1 > 0 {
		e.ChannelId = buf.PopInt()
	}
	if flags&1<<2 > 0 {
		e.ChannelPost = buf.PopInt()
	}
	if flags&1<<3 > 0 {
		e.PostAuthor = buf.PopString()
	}
	if flags&1<<4 > 0 {
		e.SavedFromPeer = Peer(buf.PopObj())
	}
	if flags&1<<4 > 0 {
		e.SavedFromMsgId = buf.PopInt()
	}
	if flags&1<<6 > 0 {
		e.PsaType = buf.PopString()
	}
}

type InputAppEvent struct {
	Time float64   `validate:"required"`
	Type string    `validate:"required"`
	Peer int64     `validate:"required"`
	Data JSONValue `validate:"required"`
}

func (e *InputAppEvent) CRC() uint32 {
	return uint32(0x1d1b1245)
}
func (e *InputAppEvent) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutDouble(e.Time)
	buf.PutString(e.Type)
	buf.PutLong(e.Peer)
	buf.PutRawBytes(e.Data.Encode())
	return buf.Result()
}

func (e *InputAppEvent) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Time = buf.PopDouble()
	e.Type = buf.PopString()
	e.Peer = buf.PopLong()
	e.Data = JSONValue(buf.PopObj())
}

type KeyboardButtonRow struct {
	Buttons []KeyboardButton `validate:"required"`
}

func (e *KeyboardButtonRow) CRC() uint32 {
	return uint32(0x77608b83)
}
func (e *KeyboardButtonRow) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutVector(e.Buttons)
	return buf.Result()
}

func (e *KeyboardButtonRow) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Buttons = buf.PopVector(reflect.TypeOf(KeyboardButton{})).([]KeyboardButton)
}

type MessagesBotResults struct {
	__flagsPosition struct{}           // flags param position `validate:"required"`
	Gallery         bool               `flag:"0,encoded_in_bitflags"`
	QueryId         int64              `validate:"required"`
	NextOffset      string             `flag:"1"`
	SwitchPm        *InlineBotSwitchPM `flag:"2"`
	Results         []BotInlineResult  `validate:"required"`
	CacheTime       int32              `validate:"required"`
	Users           []User             `validate:"required"`
}

func (e *MessagesBotResults) CRC() uint32 {
	return uint32(0x947ca848)
}
func (e *MessagesBotResults) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Gallery) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.NextOffset) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.SwitchPm) {
		flag |= 1 << 2
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Gallery) {
	}
	buf.PutLong(e.QueryId)
	if !zero.IsZeroVal(e.NextOffset) {
		buf.PutString(e.NextOffset)
	}
	if !zero.IsZeroVal(e.SwitchPm) {
		buf.PutRawBytes(e.SwitchPm.Encode())
	}
	buf.PutVector(e.Results)
	buf.PutInt(e.CacheTime)
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *MessagesBotResults) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Gallery = buf.PopBool()
	}
	e.QueryId = buf.PopLong()
	if flags&1<<1 > 0 {
		e.NextOffset = buf.PopString()
	}
	if flags&1<<2 > 0 {
		e.SwitchPm = buf.PopObj().(*InlineBotSwitchPM)
	}
	e.Results = buf.PopVector(reflect.TypeOf(BotInlineResult{})).([]BotInlineResult)
	e.CacheTime = buf.PopInt()
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type EmojiURL struct {
	Url string `validate:"required"`
}

func (e *EmojiURL) CRC() uint32 {
	return uint32(0xa575739d)
}
func (e *EmojiURL) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Url)
	return buf.Result()
}

func (e *EmojiURL) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Url = buf.PopString()
}

type BankCardOpenUrl struct {
	Url  string `validate:"required"`
	Name string `validate:"required"`
}

func (e *BankCardOpenUrl) CRC() uint32 {
	return uint32(0xf568028a)
}
func (e *BankCardOpenUrl) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Url)
	buf.PutString(e.Name)
	return buf.Result()
}

func (e *BankCardOpenUrl) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Url = buf.PopString()
	e.Name = buf.PopString()
}

type AccessPointRule struct {
	PhonePrefixRules string   `validate:"required"`
	DcId             int32    `validate:"required"`
	Ips              []IpPort `validate:"required"`
}

func (e *AccessPointRule) CRC() uint32 {
	return uint32(0x4679b65f)
}
func (e *AccessPointRule) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.PhonePrefixRules)
	buf.PutInt(e.DcId)
	buf.PutVector(e.Ips)
	return buf.Result()
}

func (e *AccessPointRule) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.PhonePrefixRules = buf.PopString()
	e.DcId = buf.PopInt()
	e.Ips = buf.PopVector(reflect.TypeOf(IpPort{})).([]IpPort)
}

type PeerSettings struct {
	__flagsPosition       struct{} // flags param position `validate:"required"`
	ReportSpam            bool     `flag:"0,encoded_in_bitflags"`
	AddContact            bool     `flag:"1,encoded_in_bitflags"`
	BlockContact          bool     `flag:"2,encoded_in_bitflags"`
	ShareContact          bool     `flag:"3,encoded_in_bitflags"`
	NeedContactsException bool     `flag:"4,encoded_in_bitflags"`
	ReportGeo             bool     `flag:"5,encoded_in_bitflags"`
	Autoarchived          bool     `flag:"7,encoded_in_bitflags"`
	GeoDistance           int32    `flag:"6"`
}

func (e *PeerSettings) CRC() uint32 {
	return uint32(0x733f2961)
}
func (e *PeerSettings) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.ReportSpam) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.AddContact) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.BlockContact) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.ShareContact) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.NeedContactsException) {
		flag |= 1 << 4
	}
	if !zero.IsZeroVal(e.ReportGeo) {
		flag |= 1 << 5
	}
	if !zero.IsZeroVal(e.GeoDistance) {
		flag |= 1 << 6
	}
	if !zero.IsZeroVal(e.Autoarchived) {
		flag |= 1 << 7
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.ReportSpam) {
	}
	if !zero.IsZeroVal(e.AddContact) {
	}
	if !zero.IsZeroVal(e.BlockContact) {
	}
	if !zero.IsZeroVal(e.ShareContact) {
	}
	if !zero.IsZeroVal(e.NeedContactsException) {
	}
	if !zero.IsZeroVal(e.ReportGeo) {
	}
	if !zero.IsZeroVal(e.Autoarchived) {
	}
	if !zero.IsZeroVal(e.GeoDistance) {
		buf.PutInt(e.GeoDistance)
	}
	return buf.Result()
}

func (e *PeerSettings) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.ReportSpam = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.AddContact = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.BlockContact = buf.PopBool()
	}
	if flags&1<<3 > 0 {
		e.ShareContact = buf.PopBool()
	}
	if flags&1<<4 > 0 {
		e.NeedContactsException = buf.PopBool()
	}
	if flags&1<<5 > 0 {
		e.ReportGeo = buf.PopBool()
	}
	if flags&1<<7 > 0 {
		e.Autoarchived = buf.PopBool()
	}
	if flags&1<<6 > 0 {
		e.GeoDistance = buf.PopInt()
	}
}

type Folder struct {
	__flagsPosition           struct{}  // flags param position `validate:"required"`
	AutofillNewBroadcasts     bool      `flag:"0,encoded_in_bitflags"`
	AutofillPublicGroups      bool      `flag:"1,encoded_in_bitflags"`
	AutofillNewCorrespondents bool      `flag:"2,encoded_in_bitflags"`
	Id                        int32     `validate:"required"`
	Title                     string    `validate:"required"`
	Photo                     ChatPhoto `flag:"3"`
}

func (e *Folder) CRC() uint32 {
	return uint32(0xff544e65)
}
func (e *Folder) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.AutofillNewBroadcasts) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.AutofillPublicGroups) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.AutofillNewCorrespondents) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.Photo) {
		flag |= 1 << 3
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.AutofillNewBroadcasts) {
	}
	if !zero.IsZeroVal(e.AutofillPublicGroups) {
	}
	if !zero.IsZeroVal(e.AutofillNewCorrespondents) {
	}
	buf.PutInt(e.Id)
	buf.PutString(e.Title)
	if !zero.IsZeroVal(e.Photo) {
		buf.PutRawBytes(e.Photo.Encode())
	}
	return buf.Result()
}

func (e *Folder) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.AutofillNewBroadcasts = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.AutofillPublicGroups = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.AutofillNewCorrespondents = buf.PopBool()
	}
	e.Id = buf.PopInt()
	e.Title = buf.PopString()
	if flags&1<<3 > 0 {
		e.Photo = ChatPhoto(buf.PopObj())
	}
}

type FileHash struct {
	Offset int32  `validate:"required"`
	Limit  int32  `validate:"required"`
	Hash   []byte `validate:"required"`
}

func (e *FileHash) CRC() uint32 {
	return uint32(0x6242c773)
}
func (e *FileHash) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.Offset)
	buf.PutInt(e.Limit)
	buf.PutMessage(e.Hash)
	return buf.Result()
}

func (e *FileHash) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Offset = buf.PopInt()
	e.Limit = buf.PopInt()
	e.Hash = buf.PopMessage()
}

type ChatAdminRights struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	ChangeInfo      bool     `flag:"0,encoded_in_bitflags"`
	PostMessages    bool     `flag:"1,encoded_in_bitflags"`
	EditMessages    bool     `flag:"2,encoded_in_bitflags"`
	DeleteMessages  bool     `flag:"3,encoded_in_bitflags"`
	BanUsers        bool     `flag:"4,encoded_in_bitflags"`
	InviteUsers     bool     `flag:"5,encoded_in_bitflags"`
	PinMessages     bool     `flag:"7,encoded_in_bitflags"`
	AddAdmins       bool     `flag:"9,encoded_in_bitflags"`
}

func (e *ChatAdminRights) CRC() uint32 {
	return uint32(0x5fb224d5)
}
func (e *ChatAdminRights) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.ChangeInfo) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.PostMessages) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.EditMessages) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.DeleteMessages) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.BanUsers) {
		flag |= 1 << 4
	}
	if !zero.IsZeroVal(e.InviteUsers) {
		flag |= 1 << 5
	}
	if !zero.IsZeroVal(e.PinMessages) {
		flag |= 1 << 7
	}
	if !zero.IsZeroVal(e.AddAdmins) {
		flag |= 1 << 9
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.ChangeInfo) {
	}
	if !zero.IsZeroVal(e.PostMessages) {
	}
	if !zero.IsZeroVal(e.EditMessages) {
	}
	if !zero.IsZeroVal(e.DeleteMessages) {
	}
	if !zero.IsZeroVal(e.BanUsers) {
	}
	if !zero.IsZeroVal(e.InviteUsers) {
	}
	if !zero.IsZeroVal(e.PinMessages) {
	}
	if !zero.IsZeroVal(e.AddAdmins) {
	}
	return buf.Result()
}

func (e *ChatAdminRights) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.ChangeInfo = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.PostMessages = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.EditMessages = buf.PopBool()
	}
	if flags&1<<3 > 0 {
		e.DeleteMessages = buf.PopBool()
	}
	if flags&1<<4 > 0 {
		e.BanUsers = buf.PopBool()
	}
	if flags&1<<5 > 0 {
		e.InviteUsers = buf.PopBool()
	}
	if flags&1<<7 > 0 {
		e.PinMessages = buf.PopBool()
	}
	if flags&1<<9 > 0 {
		e.AddAdmins = buf.PopBool()
	}
}

type WebPageAttribute struct {
	__flagsPosition struct{}       // flags param position `validate:"required"`
	Documents       []Document     `flag:"0"`
	Settings        *ThemeSettings `flag:"1"`
}

func (e *WebPageAttribute) CRC() uint32 {
	return uint32(0x54b56617)
}
func (e *WebPageAttribute) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Documents) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Settings) {
		flag |= 1 << 1
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Documents) {
		buf.PutVector(e.Documents)
	}
	if !zero.IsZeroVal(e.Settings) {
		buf.PutRawBytes(e.Settings.Encode())
	}
	return buf.Result()
}

func (e *WebPageAttribute) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Documents = buf.PopVector(reflect.TypeOf(Document{})).([]Document)
	}
	if flags&1<<1 > 0 {
		e.Settings = buf.PopObj().(*ThemeSettings)
	}
}

type PaymentsBankCardData struct {
	Title    string             `validate:"required"`
	OpenUrls []*BankCardOpenUrl `validate:"required"`
}

func (e *PaymentsBankCardData) CRC() uint32 {
	return uint32(0x3e24e573)
}
func (e *PaymentsBankCardData) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Title)
	buf.PutVector(e.OpenUrls)
	return buf.Result()
}

func (e *PaymentsBankCardData) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Title = buf.PopString()
	e.OpenUrls = buf.PopVector(reflect.TypeOf(*BankCardOpenUrl{})).([]*BankCardOpenUrl)
}

type Invoice struct {
	__flagsPosition          struct{}        // flags param position `validate:"required"`
	Test                     bool            `flag:"0,encoded_in_bitflags"`
	NameRequested            bool            `flag:"1,encoded_in_bitflags"`
	PhoneRequested           bool            `flag:"2,encoded_in_bitflags"`
	EmailRequested           bool            `flag:"3,encoded_in_bitflags"`
	ShippingAddressRequested bool            `flag:"4,encoded_in_bitflags"`
	Flexible                 bool            `flag:"5,encoded_in_bitflags"`
	PhoneToProvider          bool            `flag:"6,encoded_in_bitflags"`
	EmailToProvider          bool            `flag:"7,encoded_in_bitflags"`
	Currency                 string          `validate:"required"`
	Prices                   []*LabeledPrice `validate:"required"`
}

func (e *Invoice) CRC() uint32 {
	return uint32(0xc30aa358)
}
func (e *Invoice) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Test) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.NameRequested) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.PhoneRequested) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.EmailRequested) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.ShippingAddressRequested) {
		flag |= 1 << 4
	}
	if !zero.IsZeroVal(e.Flexible) {
		flag |= 1 << 5
	}
	if !zero.IsZeroVal(e.PhoneToProvider) {
		flag |= 1 << 6
	}
	if !zero.IsZeroVal(e.EmailToProvider) {
		flag |= 1 << 7
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Test) {
	}
	if !zero.IsZeroVal(e.NameRequested) {
	}
	if !zero.IsZeroVal(e.PhoneRequested) {
	}
	if !zero.IsZeroVal(e.EmailRequested) {
	}
	if !zero.IsZeroVal(e.ShippingAddressRequested) {
	}
	if !zero.IsZeroVal(e.Flexible) {
	}
	if !zero.IsZeroVal(e.PhoneToProvider) {
	}
	if !zero.IsZeroVal(e.EmailToProvider) {
	}
	buf.PutString(e.Currency)
	buf.PutVector(e.Prices)
	return buf.Result()
}

func (e *Invoice) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Test = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.NameRequested = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.PhoneRequested = buf.PopBool()
	}
	if flags&1<<3 > 0 {
		e.EmailRequested = buf.PopBool()
	}
	if flags&1<<4 > 0 {
		e.ShippingAddressRequested = buf.PopBool()
	}
	if flags&1<<5 > 0 {
		e.Flexible = buf.PopBool()
	}
	if flags&1<<6 > 0 {
		e.PhoneToProvider = buf.PopBool()
	}
	if flags&1<<7 > 0 {
		e.EmailToProvider = buf.PopBool()
	}
	e.Currency = buf.PopString()
	e.Prices = buf.PopVector(reflect.TypeOf(*LabeledPrice{})).([]*LabeledPrice)
}

type InputFolderPeer struct {
	Peer     InputPeer `validate:"required"`
	FolderId int32     `validate:"required"`
}

func (e *InputFolderPeer) CRC() uint32 {
	return uint32(0xfbd2c296)
}
func (e *InputFolderPeer) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Peer.Encode())
	buf.PutInt(e.FolderId)
	return buf.Result()
}

func (e *InputFolderPeer) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Peer = InputPeer(buf.PopObj())
	e.FolderId = buf.PopInt()
}

type AccountAuthorizations struct {
	Authorizations []*Authorization `validate:"required"`
}

func (e *AccountAuthorizations) CRC() uint32 {
	return uint32(0x1250abde)
}
func (e *AccountAuthorizations) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutVector(e.Authorizations)
	return buf.Result()
}

func (e *AccountAuthorizations) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Authorizations = buf.PopVector(reflect.TypeOf(*Authorization{})).([]*Authorization)
}

type StatsURL struct {
	Url string `validate:"required"`
}

func (e *StatsURL) CRC() uint32 {
	return uint32(0x47a971e0)
}
func (e *StatsURL) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Url)
	return buf.Result()
}

func (e *StatsURL) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Url = buf.PopString()
}

type Config struct {
	__flagsPosition         struct{}    // flags param position `validate:"required"`
	PhonecallsEnabled       bool        `flag:"1,encoded_in_bitflags"`
	DefaultP2PContacts      bool        `flag:"3,encoded_in_bitflags"`
	PreloadFeaturedStickers bool        `flag:"4,encoded_in_bitflags"`
	IgnorePhoneEntities     bool        `flag:"5,encoded_in_bitflags"`
	RevokePmInbox           bool        `flag:"6,encoded_in_bitflags"`
	BlockedMode             bool        `flag:"8,encoded_in_bitflags"`
	PfsEnabled              bool        `flag:"13,encoded_in_bitflags"`
	Date                    int32       `validate:"required"`
	Expires                 int32       `validate:"required"`
	TestMode                bool        `validate:"required"`
	ThisDc                  int32       `validate:"required"`
	DcOptions               []*DcOption `validate:"required"`
	DcTxtDomainName         string      `validate:"required"`
	ChatSizeMax             int32       `validate:"required"`
	MegagroupSizeMax        int32       `validate:"required"`
	ForwardedCountMax       int32       `validate:"required"`
	OnlineUpdatePeriodMs    int32       `validate:"required"`
	OfflineBlurTimeoutMs    int32       `validate:"required"`
	OfflineIdleTimeoutMs    int32       `validate:"required"`
	OnlineCloudTimeoutMs    int32       `validate:"required"`
	NotifyCloudDelayMs      int32       `validate:"required"`
	NotifyDefaultDelayMs    int32       `validate:"required"`
	PushChatPeriodMs        int32       `validate:"required"`
	PushChatLimit           int32       `validate:"required"`
	SavedGifsLimit          int32       `validate:"required"`
	EditTimeLimit           int32       `validate:"required"`
	RevokeTimeLimit         int32       `validate:"required"`
	RevokePmTimeLimit       int32       `validate:"required"`
	RatingEDecay            int32       `validate:"required"`
	StickersRecentLimit     int32       `validate:"required"`
	StickersFavedLimit      int32       `validate:"required"`
	ChannelsReadMediaPeriod int32       `validate:"required"`
	TmpSessions             int32       `flag:"0"`
	PinnedDialogsCountMax   int32       `validate:"required"`
	PinnedInfolderCountMax  int32       `validate:"required"`
	CallReceiveTimeoutMs    int32       `validate:"required"`
	CallRingTimeoutMs       int32       `validate:"required"`
	CallConnectTimeoutMs    int32       `validate:"required"`
	CallPacketTimeoutMs     int32       `validate:"required"`
	MeUrlPrefix             string      `validate:"required"`
	AutoupdateUrlPrefix     string      `flag:"7"`
	GifSearchUsername       string      `flag:"9"`
	VenueSearchUsername     string      `flag:"10"`
	ImgSearchUsername       string      `flag:"11"`
	StaticMapsProvider      string      `flag:"12"`
	CaptionLengthMax        int32       `validate:"required"`
	MessageLengthMax        int32       `validate:"required"`
	WebfileDcId             int32       `validate:"required"`
	SuggestedLangCode       string      `flag:"2"`
	LangPackVersion         int32       `flag:"2"`
	BaseLangPackVersion     int32       `flag:"2"`
}

func (e *Config) CRC() uint32 {
	return uint32(0x330b4067)
}
func (e *Config) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.TmpSessions) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.PhonecallsEnabled) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.SuggestedLangCode) || !zero.IsZeroVal(e.LangPackVersion) || !zero.IsZeroVal(e.BaseLangPackVersion) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.DefaultP2PContacts) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.PreloadFeaturedStickers) {
		flag |= 1 << 4
	}
	if !zero.IsZeroVal(e.IgnorePhoneEntities) {
		flag |= 1 << 5
	}
	if !zero.IsZeroVal(e.RevokePmInbox) {
		flag |= 1 << 6
	}
	if !zero.IsZeroVal(e.AutoupdateUrlPrefix) {
		flag |= 1 << 7
	}
	if !zero.IsZeroVal(e.BlockedMode) {
		flag |= 1 << 8
	}
	if !zero.IsZeroVal(e.GifSearchUsername) {
		flag |= 1 << 9
	}
	if !zero.IsZeroVal(e.VenueSearchUsername) {
		flag |= 1 << 10
	}
	if !zero.IsZeroVal(e.ImgSearchUsername) {
		flag |= 1 << 11
	}
	if !zero.IsZeroVal(e.StaticMapsProvider) {
		flag |= 1 << 12
	}
	if !zero.IsZeroVal(e.PfsEnabled) {
		flag |= 1 << 13
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.PhonecallsEnabled) {
	}
	if !zero.IsZeroVal(e.DefaultP2PContacts) {
	}
	if !zero.IsZeroVal(e.PreloadFeaturedStickers) {
	}
	if !zero.IsZeroVal(e.IgnorePhoneEntities) {
	}
	if !zero.IsZeroVal(e.RevokePmInbox) {
	}
	if !zero.IsZeroVal(e.BlockedMode) {
	}
	if !zero.IsZeroVal(e.PfsEnabled) {
	}
	buf.PutInt(e.Date)
	buf.PutInt(e.Expires)
	buf.PutBool(e.TestMode)
	buf.PutInt(e.ThisDc)
	buf.PutVector(e.DcOptions)
	buf.PutString(e.DcTxtDomainName)
	buf.PutInt(e.ChatSizeMax)
	buf.PutInt(e.MegagroupSizeMax)
	buf.PutInt(e.ForwardedCountMax)
	buf.PutInt(e.OnlineUpdatePeriodMs)
	buf.PutInt(e.OfflineBlurTimeoutMs)
	buf.PutInt(e.OfflineIdleTimeoutMs)
	buf.PutInt(e.OnlineCloudTimeoutMs)
	buf.PutInt(e.NotifyCloudDelayMs)
	buf.PutInt(e.NotifyDefaultDelayMs)
	buf.PutInt(e.PushChatPeriodMs)
	buf.PutInt(e.PushChatLimit)
	buf.PutInt(e.SavedGifsLimit)
	buf.PutInt(e.EditTimeLimit)
	buf.PutInt(e.RevokeTimeLimit)
	buf.PutInt(e.RevokePmTimeLimit)
	buf.PutInt(e.RatingEDecay)
	buf.PutInt(e.StickersRecentLimit)
	buf.PutInt(e.StickersFavedLimit)
	buf.PutInt(e.ChannelsReadMediaPeriod)
	if !zero.IsZeroVal(e.TmpSessions) {
		buf.PutInt(e.TmpSessions)
	}
	buf.PutInt(e.PinnedDialogsCountMax)
	buf.PutInt(e.PinnedInfolderCountMax)
	buf.PutInt(e.CallReceiveTimeoutMs)
	buf.PutInt(e.CallRingTimeoutMs)
	buf.PutInt(e.CallConnectTimeoutMs)
	buf.PutInt(e.CallPacketTimeoutMs)
	buf.PutString(e.MeUrlPrefix)
	if !zero.IsZeroVal(e.AutoupdateUrlPrefix) {
		buf.PutString(e.AutoupdateUrlPrefix)
	}
	if !zero.IsZeroVal(e.GifSearchUsername) {
		buf.PutString(e.GifSearchUsername)
	}
	if !zero.IsZeroVal(e.VenueSearchUsername) {
		buf.PutString(e.VenueSearchUsername)
	}
	if !zero.IsZeroVal(e.ImgSearchUsername) {
		buf.PutString(e.ImgSearchUsername)
	}
	if !zero.IsZeroVal(e.StaticMapsProvider) {
		buf.PutString(e.StaticMapsProvider)
	}
	buf.PutInt(e.CaptionLengthMax)
	buf.PutInt(e.MessageLengthMax)
	buf.PutInt(e.WebfileDcId)
	if !zero.IsZeroVal(e.SuggestedLangCode) {
		buf.PutString(e.SuggestedLangCode)
	}
	if !zero.IsZeroVal(e.LangPackVersion) {
		buf.PutInt(e.LangPackVersion)
	}
	if !zero.IsZeroVal(e.BaseLangPackVersion) {
		buf.PutInt(e.BaseLangPackVersion)
	}
	return buf.Result()
}

func (e *Config) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<1 > 0 {
		e.PhonecallsEnabled = buf.PopBool()
	}
	if flags&1<<3 > 0 {
		e.DefaultP2PContacts = buf.PopBool()
	}
	if flags&1<<4 > 0 {
		e.PreloadFeaturedStickers = buf.PopBool()
	}
	if flags&1<<5 > 0 {
		e.IgnorePhoneEntities = buf.PopBool()
	}
	if flags&1<<6 > 0 {
		e.RevokePmInbox = buf.PopBool()
	}
	if flags&1<<8 > 0 {
		e.BlockedMode = buf.PopBool()
	}
	if flags&1<<13 > 0 {
		e.PfsEnabled = buf.PopBool()
	}
	e.Date = buf.PopInt()
	e.Expires = buf.PopInt()
	e.TestMode = buf.PopBool()
	e.ThisDc = buf.PopInt()
	e.DcOptions = buf.PopVector(reflect.TypeOf(*DcOption{})).([]*DcOption)
	e.DcTxtDomainName = buf.PopString()
	e.ChatSizeMax = buf.PopInt()
	e.MegagroupSizeMax = buf.PopInt()
	e.ForwardedCountMax = buf.PopInt()
	e.OnlineUpdatePeriodMs = buf.PopInt()
	e.OfflineBlurTimeoutMs = buf.PopInt()
	e.OfflineIdleTimeoutMs = buf.PopInt()
	e.OnlineCloudTimeoutMs = buf.PopInt()
	e.NotifyCloudDelayMs = buf.PopInt()
	e.NotifyDefaultDelayMs = buf.PopInt()
	e.PushChatPeriodMs = buf.PopInt()
	e.PushChatLimit = buf.PopInt()
	e.SavedGifsLimit = buf.PopInt()
	e.EditTimeLimit = buf.PopInt()
	e.RevokeTimeLimit = buf.PopInt()
	e.RevokePmTimeLimit = buf.PopInt()
	e.RatingEDecay = buf.PopInt()
	e.StickersRecentLimit = buf.PopInt()
	e.StickersFavedLimit = buf.PopInt()
	e.ChannelsReadMediaPeriod = buf.PopInt()
	if flags&1<<0 > 0 {
		e.TmpSessions = buf.PopInt()
	}
	e.PinnedDialogsCountMax = buf.PopInt()
	e.PinnedInfolderCountMax = buf.PopInt()
	e.CallReceiveTimeoutMs = buf.PopInt()
	e.CallRingTimeoutMs = buf.PopInt()
	e.CallConnectTimeoutMs = buf.PopInt()
	e.CallPacketTimeoutMs = buf.PopInt()
	e.MeUrlPrefix = buf.PopString()
	if flags&1<<7 > 0 {
		e.AutoupdateUrlPrefix = buf.PopString()
	}
	if flags&1<<9 > 0 {
		e.GifSearchUsername = buf.PopString()
	}
	if flags&1<<10 > 0 {
		e.VenueSearchUsername = buf.PopString()
	}
	if flags&1<<11 > 0 {
		e.ImgSearchUsername = buf.PopString()
	}
	if flags&1<<12 > 0 {
		e.StaticMapsProvider = buf.PopString()
	}
	e.CaptionLengthMax = buf.PopInt()
	e.MessageLengthMax = buf.PopInt()
	e.WebfileDcId = buf.PopInt()
	if flags&1<<2 > 0 {
		e.SuggestedLangCode = buf.PopString()
	}
	if flags&1<<2 > 0 {
		e.LangPackVersion = buf.PopInt()
	}
	if flags&1<<2 > 0 {
		e.BaseLangPackVersion = buf.PopInt()
	}
}

type SecureData struct {
	Data     []byte `validate:"required"`
	DataHash []byte `validate:"required"`
	Secret   []byte `validate:"required"`
}

func (e *SecureData) CRC() uint32 {
	return uint32(0x8aeabec3)
}
func (e *SecureData) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutMessage(e.Data)
	buf.PutMessage(e.DataHash)
	buf.PutMessage(e.Secret)
	return buf.Result()
}

func (e *SecureData) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Data = buf.PopMessage()
	e.DataHash = buf.PopMessage()
	e.Secret = buf.PopMessage()
}

type MessagesAffectedHistory struct {
	Pts      int32 `validate:"required"`
	PtsCount int32 `validate:"required"`
	Offset   int32 `validate:"required"`
}

func (e *MessagesAffectedHistory) CRC() uint32 {
	return uint32(0xb45c69d1)
}
func (e *MessagesAffectedHistory) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.Pts)
	buf.PutInt(e.PtsCount)
	buf.PutInt(e.Offset)
	return buf.Result()
}

func (e *MessagesAffectedHistory) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Pts = buf.PopInt()
	e.PtsCount = buf.PopInt()
	e.Offset = buf.PopInt()
}

type ContactBlocked struct {
	UserId int32 `validate:"required"`
	Date   int32 `validate:"required"`
}

func (e *ContactBlocked) CRC() uint32 {
	return uint32(0x561bc879)
}
func (e *ContactBlocked) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.UserId)
	buf.PutInt(e.Date)
	return buf.Result()
}

func (e *ContactBlocked) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.UserId = buf.PopInt()
	e.Date = buf.PopInt()
}

type InputSecureValue struct {
	__flagsPosition struct{}          // flags param position `validate:"required"`
	Type            SecureValueType   `validate:"required"`
	Data            *SecureData       `flag:"0"`
	FrontSide       InputSecureFile   `flag:"1"`
	ReverseSide     InputSecureFile   `flag:"2"`
	Selfie          InputSecureFile   `flag:"3"`
	Translation     []InputSecureFile `flag:"6"`
	Files           []InputSecureFile `flag:"4"`
	PlainData       SecurePlainData   `flag:"5"`
}

func (e *InputSecureValue) CRC() uint32 {
	return uint32(0xdb21d0a7)
}
func (e *InputSecureValue) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Data) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.FrontSide) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.ReverseSide) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.Selfie) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.Files) {
		flag |= 1 << 4
	}
	if !zero.IsZeroVal(e.PlainData) {
		flag |= 1 << 5
	}
	if !zero.IsZeroVal(e.Translation) {
		flag |= 1 << 6
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Type.Encode())
	if !zero.IsZeroVal(e.Data) {
		buf.PutRawBytes(e.Data.Encode())
	}
	if !zero.IsZeroVal(e.FrontSide) {
		buf.PutRawBytes(e.FrontSide.Encode())
	}
	if !zero.IsZeroVal(e.ReverseSide) {
		buf.PutRawBytes(e.ReverseSide.Encode())
	}
	if !zero.IsZeroVal(e.Selfie) {
		buf.PutRawBytes(e.Selfie.Encode())
	}
	if !zero.IsZeroVal(e.Translation) {
		buf.PutVector(e.Translation)
	}
	if !zero.IsZeroVal(e.Files) {
		buf.PutVector(e.Files)
	}
	if !zero.IsZeroVal(e.PlainData) {
		buf.PutRawBytes(e.PlainData.Encode())
	}
	return buf.Result()
}

func (e *InputSecureValue) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	e.Type = *(buf.PopObj().(*SecureValueType))
	if flags&1<<0 > 0 {
		e.Data = buf.PopObj().(*SecureData)
	}
	if flags&1<<1 > 0 {
		e.FrontSide = InputSecureFile(buf.PopObj())
	}
	if flags&1<<2 > 0 {
		e.ReverseSide = InputSecureFile(buf.PopObj())
	}
	if flags&1<<3 > 0 {
		e.Selfie = InputSecureFile(buf.PopObj())
	}
	if flags&1<<6 > 0 {
		e.Translation = buf.PopVector(reflect.TypeOf(InputSecureFile{})).([]InputSecureFile)
	}
	if flags&1<<4 > 0 {
		e.Files = buf.PopVector(reflect.TypeOf(InputSecureFile{})).([]InputSecureFile)
	}
	if flags&1<<5 > 0 {
		e.PlainData = SecurePlainData(buf.PopObj())
	}
}

type AuthSentCode struct {
	__flagsPosition struct{}         // flags param position `validate:"required"`
	Type            AuthSentCodeType `validate:"required"`
	PhoneCodeHash   string           `validate:"required"`
	NextType        AuthCodeType     `flag:"1"`
	Timeout         int32            `flag:"2"`
}

func (e *AuthSentCode) CRC() uint32 {
	return uint32(0x5e002502)
}
func (e *AuthSentCode) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.NextType) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.Timeout) {
		flag |= 1 << 2
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Type.Encode())
	buf.PutString(e.PhoneCodeHash)
	if !zero.IsZeroVal(e.NextType) {
		buf.PutRawBytes(e.NextType.Encode())
	}
	if !zero.IsZeroVal(e.Timeout) {
		buf.PutInt(e.Timeout)
	}
	return buf.Result()
}

func (e *AuthSentCode) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	e.Type = AuthSentCodeType(buf.PopObj())
	e.PhoneCodeHash = buf.PopString()
	if flags&1<<1 > 0 {
		e.NextType = *(buf.PopObj().(*AuthCodeType))
	}
	if flags&1<<2 > 0 {
		e.Timeout = buf.PopInt()
	}
}

type HelpInviteText struct {
	Message string `validate:"required"`
}

func (e *HelpInviteText) CRC() uint32 {
	return uint32(0x18cb9f78)
}
func (e *HelpInviteText) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Message)
	return buf.Result()
}

func (e *HelpInviteText) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Message = buf.PopString()
}

type HighScore struct {
	Pos    int32 `validate:"required"`
	UserId int32 `validate:"required"`
	Score  int32 `validate:"required"`
}

func (e *HighScore) CRC() uint32 {
	return uint32(0x58fffcd0)
}
func (e *HighScore) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.Pos)
	buf.PutInt(e.UserId)
	buf.PutInt(e.Score)
	return buf.Result()
}

func (e *HighScore) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Pos = buf.PopInt()
	e.UserId = buf.PopInt()
	e.Score = buf.PopInt()
}

type PollAnswer struct {
	Text   string `validate:"required"`
	Option []byte `validate:"required"`
}

func (e *PollAnswer) CRC() uint32 {
	return uint32(0x6ca9c2e9)
}
func (e *PollAnswer) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Text)
	buf.PutMessage(e.Option)
	return buf.Result()
}

func (e *PollAnswer) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Text = buf.PopString()
	e.Option = buf.PopMessage()
}

type DialogFilter struct {
	__flagsPosition struct{}    // flags param position `validate:"required"`
	Contacts        bool        `flag:"0,encoded_in_bitflags"`
	NonContacts     bool        `flag:"1,encoded_in_bitflags"`
	Groups          bool        `flag:"2,encoded_in_bitflags"`
	Broadcasts      bool        `flag:"3,encoded_in_bitflags"`
	Bots            bool        `flag:"4,encoded_in_bitflags"`
	ExcludeMuted    bool        `flag:"11,encoded_in_bitflags"`
	ExcludeRead     bool        `flag:"12,encoded_in_bitflags"`
	ExcludeArchived bool        `flag:"13,encoded_in_bitflags"`
	Id              int32       `validate:"required"`
	Title           string      `validate:"required"`
	Emoticon        string      `flag:"25"`
	PinnedPeers     []InputPeer `validate:"required"`
	IncludePeers    []InputPeer `validate:"required"`
	ExcludePeers    []InputPeer `validate:"required"`
}

func (e *DialogFilter) CRC() uint32 {
	return uint32(0x7438f7e8)
}
func (e *DialogFilter) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Contacts) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.NonContacts) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.Groups) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.Broadcasts) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.Bots) {
		flag |= 1 << 4
	}
	if !zero.IsZeroVal(e.ExcludeMuted) {
		flag |= 1 << 11
	}
	if !zero.IsZeroVal(e.ExcludeRead) {
		flag |= 1 << 12
	}
	if !zero.IsZeroVal(e.ExcludeArchived) {
		flag |= 1 << 13
	}
	if !zero.IsZeroVal(e.Emoticon) {
		flag |= 1 << 25
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Contacts) {
	}
	if !zero.IsZeroVal(e.NonContacts) {
	}
	if !zero.IsZeroVal(e.Groups) {
	}
	if !zero.IsZeroVal(e.Broadcasts) {
	}
	if !zero.IsZeroVal(e.Bots) {
	}
	if !zero.IsZeroVal(e.ExcludeMuted) {
	}
	if !zero.IsZeroVal(e.ExcludeRead) {
	}
	if !zero.IsZeroVal(e.ExcludeArchived) {
	}
	buf.PutInt(e.Id)
	buf.PutString(e.Title)
	if !zero.IsZeroVal(e.Emoticon) {
		buf.PutString(e.Emoticon)
	}
	buf.PutVector(e.PinnedPeers)
	buf.PutVector(e.IncludePeers)
	buf.PutVector(e.ExcludePeers)
	return buf.Result()
}

func (e *DialogFilter) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Contacts = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.NonContacts = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.Groups = buf.PopBool()
	}
	if flags&1<<3 > 0 {
		e.Broadcasts = buf.PopBool()
	}
	if flags&1<<4 > 0 {
		e.Bots = buf.PopBool()
	}
	if flags&1<<11 > 0 {
		e.ExcludeMuted = buf.PopBool()
	}
	if flags&1<<12 > 0 {
		e.ExcludeRead = buf.PopBool()
	}
	if flags&1<<13 > 0 {
		e.ExcludeArchived = buf.PopBool()
	}
	e.Id = buf.PopInt()
	e.Title = buf.PopString()
	if flags&1<<25 > 0 {
		e.Emoticon = buf.PopString()
	}
	e.PinnedPeers = buf.PopVector(reflect.TypeOf(InputPeer{})).([]InputPeer)
	e.IncludePeers = buf.PopVector(reflect.TypeOf(InputPeer{})).([]InputPeer)
	e.ExcludePeers = buf.PopVector(reflect.TypeOf(InputPeer{})).([]InputPeer)
}

type ExportedMessageLink struct {
	Link string `validate:"required"`
	Html string `validate:"required"`
}

func (e *ExportedMessageLink) CRC() uint32 {
	return uint32(0x5dab1af4)
}
func (e *ExportedMessageLink) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Link)
	buf.PutString(e.Html)
	return buf.Result()
}

func (e *ExportedMessageLink) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Link = buf.PopString()
	e.Html = buf.PopString()
}

type PageRelatedArticle struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	Url             string   `validate:"required"`
	WebpageId       int64    `validate:"required"`
	Title           string   `flag:"0"`
	Description     string   `flag:"1"`
	PhotoId         int64    `flag:"2"`
	Author          string   `flag:"3"`
	PublishedDate   int32    `flag:"4"`
}

func (e *PageRelatedArticle) CRC() uint32 {
	return uint32(0xb390dc08)
}
func (e *PageRelatedArticle) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Title) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Description) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.PhotoId) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.Author) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.PublishedDate) {
		flag |= 1 << 4
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Url)
	buf.PutLong(e.WebpageId)
	if !zero.IsZeroVal(e.Title) {
		buf.PutString(e.Title)
	}
	if !zero.IsZeroVal(e.Description) {
		buf.PutString(e.Description)
	}
	if !zero.IsZeroVal(e.PhotoId) {
		buf.PutLong(e.PhotoId)
	}
	if !zero.IsZeroVal(e.Author) {
		buf.PutString(e.Author)
	}
	if !zero.IsZeroVal(e.PublishedDate) {
		buf.PutInt(e.PublishedDate)
	}
	return buf.Result()
}

func (e *PageRelatedArticle) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	e.Url = buf.PopString()
	e.WebpageId = buf.PopLong()
	if flags&1<<0 > 0 {
		e.Title = buf.PopString()
	}
	if flags&1<<1 > 0 {
		e.Description = buf.PopString()
	}
	if flags&1<<2 > 0 {
		e.PhotoId = buf.PopLong()
	}
	if flags&1<<3 > 0 {
		e.Author = buf.PopString()
	}
	if flags&1<<4 > 0 {
		e.PublishedDate = buf.PopInt()
	}
}

type PollAnswerVoters struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	Chosen          bool     `flag:"0,encoded_in_bitflags"`
	Correct         bool     `flag:"1,encoded_in_bitflags"`
	Option          []byte   `validate:"required"`
	Voters          int32    `validate:"required"`
}

func (e *PollAnswerVoters) CRC() uint32 {
	return uint32(0x3b6ddad2)
}
func (e *PollAnswerVoters) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Chosen) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Correct) {
		flag |= 1 << 1
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Chosen) {
	}
	if !zero.IsZeroVal(e.Correct) {
	}
	buf.PutMessage(e.Option)
	buf.PutInt(e.Voters)
	return buf.Result()
}

func (e *PollAnswerVoters) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Chosen = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.Correct = buf.PopBool()
	}
	e.Option = buf.PopMessage()
	e.Voters = buf.PopInt()
}

type AccountTmpPassword struct {
	TmpPassword []byte `validate:"required"`
	ValidUntil  int32  `validate:"required"`
}

func (e *AccountTmpPassword) CRC() uint32 {
	return uint32(0xdb64fd34)
}
func (e *AccountTmpPassword) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutMessage(e.TmpPassword)
	buf.PutInt(e.ValidUntil)
	return buf.Result()
}

func (e *AccountTmpPassword) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.TmpPassword = buf.PopMessage()
	e.ValidUntil = buf.PopInt()
}

type WallPaperSettings struct {
	__flagsPosition       struct{} // flags param position `validate:"required"`
	Blur                  bool     `flag:"1,encoded_in_bitflags"`
	Motion                bool     `flag:"2,encoded_in_bitflags"`
	BackgroundColor       int32    `flag:"0"`
	SecondBackgroundColor int32    `flag:"4"`
	Intensity             int32    `flag:"3"`
	Rotation              int32    `flag:"4"`
}

func (e *WallPaperSettings) CRC() uint32 {
	return uint32(0x5086cf8)
}
func (e *WallPaperSettings) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.BackgroundColor) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Blur) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.Motion) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.Intensity) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.SecondBackgroundColor) || !zero.IsZeroVal(e.Rotation) {
		flag |= 1 << 4
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Blur) {
	}
	if !zero.IsZeroVal(e.Motion) {
	}
	if !zero.IsZeroVal(e.BackgroundColor) {
		buf.PutInt(e.BackgroundColor)
	}
	if !zero.IsZeroVal(e.SecondBackgroundColor) {
		buf.PutInt(e.SecondBackgroundColor)
	}
	if !zero.IsZeroVal(e.Intensity) {
		buf.PutInt(e.Intensity)
	}
	if !zero.IsZeroVal(e.Rotation) {
		buf.PutInt(e.Rotation)
	}
	return buf.Result()
}

func (e *WallPaperSettings) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<1 > 0 {
		e.Blur = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.Motion = buf.PopBool()
	}
	if flags&1<<0 > 0 {
		e.BackgroundColor = buf.PopInt()
	}
	if flags&1<<4 > 0 {
		e.SecondBackgroundColor = buf.PopInt()
	}
	if flags&1<<3 > 0 {
		e.Intensity = buf.PopInt()
	}
	if flags&1<<4 > 0 {
		e.Rotation = buf.PopInt()
	}
}

type HelpTermsOfService struct {
	__flagsPosition struct{}        // flags param position `validate:"required"`
	Popup           bool            `flag:"0,encoded_in_bitflags"`
	Id              *DataJSON       `validate:"required"`
	Text            string          `validate:"required"`
	Entities        []MessageEntity `validate:"required"`
	MinAgeConfirm   int32           `flag:"1"`
}

func (e *HelpTermsOfService) CRC() uint32 {
	return uint32(0x780a0310)
}
func (e *HelpTermsOfService) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Popup) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.MinAgeConfirm) {
		flag |= 1 << 1
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Popup) {
	}
	buf.PutRawBytes(e.Id.Encode())
	buf.PutString(e.Text)
	buf.PutVector(e.Entities)
	if !zero.IsZeroVal(e.MinAgeConfirm) {
		buf.PutInt(e.MinAgeConfirm)
	}
	return buf.Result()
}

func (e *HelpTermsOfService) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Popup = buf.PopBool()
	}
	e.Id = buf.PopObj().(*DataJSON)
	e.Text = buf.PopString()
	e.Entities = buf.PopVector(reflect.TypeOf(MessageEntity{})).([]MessageEntity)
	if flags&1<<1 > 0 {
		e.MinAgeConfirm = buf.PopInt()
	}
}

type ChatOnlines struct {
	Onlines int32 `validate:"required"`
}

func (e *ChatOnlines) CRC() uint32 {
	return uint32(0xf041e250)
}
func (e *ChatOnlines) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.Onlines)
	return buf.Result()
}

func (e *ChatOnlines) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Onlines = buf.PopInt()
}

type DataJSON struct {
	Data string `validate:"required"`
}

func (e *DataJSON) CRC() uint32 {
	return uint32(0x7d748d04)
}
func (e *DataJSON) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Data)
	return buf.Result()
}

func (e *DataJSON) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Data = buf.PopString()
}

type AutoDownloadSettings struct {
	__flagsPosition       struct{} // flags param position `validate:"required"`
	Disabled              bool     `flag:"0,encoded_in_bitflags"`
	VideoPreloadLarge     bool     `flag:"1,encoded_in_bitflags"`
	AudioPreloadNext      bool     `flag:"2,encoded_in_bitflags"`
	PhonecallsLessData    bool     `flag:"3,encoded_in_bitflags"`
	PhotoSizeMax          int32    `validate:"required"`
	VideoSizeMax          int32    `validate:"required"`
	FileSizeMax           int32    `validate:"required"`
	VideoUploadMaxbitrate int32    `validate:"required"`
}

func (e *AutoDownloadSettings) CRC() uint32 {
	return uint32(0xe04232f3)
}
func (e *AutoDownloadSettings) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Disabled) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.VideoPreloadLarge) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.AudioPreloadNext) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.PhonecallsLessData) {
		flag |= 1 << 3
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Disabled) {
	}
	if !zero.IsZeroVal(e.VideoPreloadLarge) {
	}
	if !zero.IsZeroVal(e.AudioPreloadNext) {
	}
	if !zero.IsZeroVal(e.PhonecallsLessData) {
	}
	buf.PutInt(e.PhotoSizeMax)
	buf.PutInt(e.VideoSizeMax)
	buf.PutInt(e.FileSizeMax)
	buf.PutInt(e.VideoUploadMaxbitrate)
	return buf.Result()
}

func (e *AutoDownloadSettings) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Disabled = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.VideoPreloadLarge = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.AudioPreloadNext = buf.PopBool()
	}
	if flags&1<<3 > 0 {
		e.PhonecallsLessData = buf.PopBool()
	}
	e.PhotoSizeMax = buf.PopInt()
	e.VideoSizeMax = buf.PopInt()
	e.FileSizeMax = buf.PopInt()
	e.VideoUploadMaxbitrate = buf.PopInt()
}

type EmojiLanguage struct {
	LangCode string `validate:"required"`
}

func (e *EmojiLanguage) CRC() uint32 {
	return uint32(0xb3fb5361)
}
func (e *EmojiLanguage) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.LangCode)
	return buf.Result()
}

func (e *EmojiLanguage) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.LangCode = buf.PopString()
}

type DcOption struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	Ipv6            bool     `flag:"0,encoded_in_bitflags"`
	MediaOnly       bool     `flag:"1,encoded_in_bitflags"`
	TcpoOnly        bool     `flag:"2,encoded_in_bitflags"`
	Cdn             bool     `flag:"3,encoded_in_bitflags"`
	Static          bool     `flag:"4,encoded_in_bitflags"`
	Id              int32    `validate:"required"`
	IpAddress       string   `validate:"required"`
	Port            int32    `validate:"required"`
	Secret          []byte   `flag:"10"`
}

func (e *DcOption) CRC() uint32 {
	return uint32(0x18b7a10d)
}
func (e *DcOption) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Ipv6) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.MediaOnly) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.TcpoOnly) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.Cdn) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.Static) {
		flag |= 1 << 4
	}
	if !zero.IsZeroVal(e.Secret) {
		flag |= 1 << 10
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Ipv6) {
	}
	if !zero.IsZeroVal(e.MediaOnly) {
	}
	if !zero.IsZeroVal(e.TcpoOnly) {
	}
	if !zero.IsZeroVal(e.Cdn) {
	}
	if !zero.IsZeroVal(e.Static) {
	}
	buf.PutInt(e.Id)
	buf.PutString(e.IpAddress)
	buf.PutInt(e.Port)
	if !zero.IsZeroVal(e.Secret) {
		buf.PutMessage(e.Secret)
	}
	return buf.Result()
}

func (e *DcOption) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Ipv6 = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.MediaOnly = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.TcpoOnly = buf.PopBool()
	}
	if flags&1<<3 > 0 {
		e.Cdn = buf.PopBool()
	}
	if flags&1<<4 > 0 {
		e.Static = buf.PopBool()
	}
	e.Id = buf.PopInt()
	e.IpAddress = buf.PopString()
	e.Port = buf.PopInt()
	if flags&1<<10 > 0 {
		e.Secret = buf.PopMessage()
	}
}

type MessagesStickerSet struct {
	Set       *StickerSet    `validate:"required"`
	Packs     []*StickerPack `validate:"required"`
	Documents []Document     `validate:"required"`
}

func (e *MessagesStickerSet) CRC() uint32 {
	return uint32(0xb60a24a6)
}
func (e *MessagesStickerSet) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Set.Encode())
	buf.PutVector(e.Packs)
	buf.PutVector(e.Documents)
	return buf.Result()
}

func (e *MessagesStickerSet) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Set = buf.PopObj().(*StickerSet)
	e.Packs = buf.PopVector(reflect.TypeOf(*StickerPack{})).([]*StickerPack)
	e.Documents = buf.PopVector(reflect.TypeOf(Document{})).([]Document)
}

type TopPeerCategoryPeers struct {
	Category TopPeerCategory `validate:"required"`
	Count    int32           `validate:"required"`
	Peers    []*TopPeer      `validate:"required"`
}

func (e *TopPeerCategoryPeers) CRC() uint32 {
	return uint32(0xfb834291)
}
func (e *TopPeerCategoryPeers) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Category.Encode())
	buf.PutInt(e.Count)
	buf.PutVector(e.Peers)
	return buf.Result()
}

func (e *TopPeerCategoryPeers) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Category = *(buf.PopObj().(*TopPeerCategory))
	e.Count = buf.PopInt()
	e.Peers = buf.PopVector(reflect.TypeOf(*TopPeer{})).([]*TopPeer)
}

type AccountPassword struct {
	__flagsPosition         struct{}              // flags param position `validate:"required"`
	HasRecovery             bool                  `flag:"0,encoded_in_bitflags"`
	HasSecureValues         bool                  `flag:"1,encoded_in_bitflags"`
	HasPassword             bool                  `flag:"2,encoded_in_bitflags"`
	CurrentAlgo             PasswordKdfAlgo       `flag:"2"`
	SrpB                    []byte                `flag:"2"`
	SrpId                   int64                 `flag:"2"`
	Hint                    string                `flag:"3"`
	EmailUnconfirmedPattern string                `flag:"4"`
	NewAlgo                 PasswordKdfAlgo       `validate:"required"`
	NewSecureAlgo           SecurePasswordKdfAlgo `validate:"required"`
	SecureRandom            []byte                `validate:"required"`
}

func (e *AccountPassword) CRC() uint32 {
	return uint32(0xad2641f8)
}
func (e *AccountPassword) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.HasRecovery) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.HasSecureValues) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.HasPassword) || !zero.IsZeroVal(e.CurrentAlgo) || !zero.IsZeroVal(e.SrpB) || !zero.IsZeroVal(e.SrpId) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.Hint) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.EmailUnconfirmedPattern) {
		flag |= 1 << 4
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.HasRecovery) {
	}
	if !zero.IsZeroVal(e.HasSecureValues) {
	}
	if !zero.IsZeroVal(e.HasPassword) {
	}
	if !zero.IsZeroVal(e.CurrentAlgo) {
		buf.PutRawBytes(e.CurrentAlgo.Encode())
	}
	if !zero.IsZeroVal(e.SrpB) {
		buf.PutMessage(e.SrpB)
	}
	if !zero.IsZeroVal(e.SrpId) {
		buf.PutLong(e.SrpId)
	}
	if !zero.IsZeroVal(e.Hint) {
		buf.PutString(e.Hint)
	}
	if !zero.IsZeroVal(e.EmailUnconfirmedPattern) {
		buf.PutString(e.EmailUnconfirmedPattern)
	}
	buf.PutRawBytes(e.NewAlgo.Encode())
	buf.PutRawBytes(e.NewSecureAlgo.Encode())
	buf.PutMessage(e.SecureRandom)
	return buf.Result()
}

func (e *AccountPassword) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.HasRecovery = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.HasSecureValues = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.HasPassword = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.CurrentAlgo = PasswordKdfAlgo(buf.PopObj())
	}
	if flags&1<<2 > 0 {
		e.SrpB = buf.PopMessage()
	}
	if flags&1<<2 > 0 {
		e.SrpId = buf.PopLong()
	}
	if flags&1<<3 > 0 {
		e.Hint = buf.PopString()
	}
	if flags&1<<4 > 0 {
		e.EmailUnconfirmedPattern = buf.PopString()
	}
	e.NewAlgo = PasswordKdfAlgo(buf.PopObj())
	e.NewSecureAlgo = SecurePasswordKdfAlgo(buf.PopObj())
	e.SecureRandom = buf.PopMessage()
}

type InputWebDocument struct {
	Url        string              `validate:"required"`
	Size       int32               `validate:"required"`
	MimeType   string              `validate:"required"`
	Attributes []DocumentAttribute `validate:"required"`
}

func (e *InputWebDocument) CRC() uint32 {
	return uint32(0x9bed434d)
}
func (e *InputWebDocument) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Url)
	buf.PutInt(e.Size)
	buf.PutString(e.MimeType)
	buf.PutVector(e.Attributes)
	return buf.Result()
}

func (e *InputWebDocument) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Url = buf.PopString()
	e.Size = buf.PopInt()
	e.MimeType = buf.PopString()
	e.Attributes = buf.PopVector(reflect.TypeOf(DocumentAttribute{})).([]DocumentAttribute)
}

type PageTableCell struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	Header          bool     `flag:"0,encoded_in_bitflags"`
	AlignCenter     bool     `flag:"3,encoded_in_bitflags"`
	AlignRight      bool     `flag:"4,encoded_in_bitflags"`
	ValignMiddle    bool     `flag:"5,encoded_in_bitflags"`
	ValignBottom    bool     `flag:"6,encoded_in_bitflags"`
	Text            RichText `flag:"7"`
	Colspan         int32    `flag:"1"`
	Rowspan         int32    `flag:"2"`
}

func (e *PageTableCell) CRC() uint32 {
	return uint32(0x34566b6a)
}
func (e *PageTableCell) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Header) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Colspan) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.Rowspan) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.AlignCenter) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.AlignRight) {
		flag |= 1 << 4
	}
	if !zero.IsZeroVal(e.ValignMiddle) {
		flag |= 1 << 5
	}
	if !zero.IsZeroVal(e.ValignBottom) {
		flag |= 1 << 6
	}
	if !zero.IsZeroVal(e.Text) {
		flag |= 1 << 7
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Header) {
	}
	if !zero.IsZeroVal(e.AlignCenter) {
	}
	if !zero.IsZeroVal(e.AlignRight) {
	}
	if !zero.IsZeroVal(e.ValignMiddle) {
	}
	if !zero.IsZeroVal(e.ValignBottom) {
	}
	if !zero.IsZeroVal(e.Text) {
		buf.PutRawBytes(e.Text.Encode())
	}
	if !zero.IsZeroVal(e.Colspan) {
		buf.PutInt(e.Colspan)
	}
	if !zero.IsZeroVal(e.Rowspan) {
		buf.PutInt(e.Rowspan)
	}
	return buf.Result()
}

func (e *PageTableCell) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Header = buf.PopBool()
	}
	if flags&1<<3 > 0 {
		e.AlignCenter = buf.PopBool()
	}
	if flags&1<<4 > 0 {
		e.AlignRight = buf.PopBool()
	}
	if flags&1<<5 > 0 {
		e.ValignMiddle = buf.PopBool()
	}
	if flags&1<<6 > 0 {
		e.ValignBottom = buf.PopBool()
	}
	if flags&1<<7 > 0 {
		e.Text = RichText(buf.PopObj())
	}
	if flags&1<<1 > 0 {
		e.Colspan = buf.PopInt()
	}
	if flags&1<<2 > 0 {
		e.Rowspan = buf.PopInt()
	}
}

type PageCaption struct {
	Text   RichText `validate:"required"`
	Credit RichText `validate:"required"`
}

func (e *PageCaption) CRC() uint32 {
	return uint32(0x6f747657)
}
func (e *PageCaption) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Text.Encode())
	buf.PutRawBytes(e.Credit.Encode())
	return buf.Result()
}

func (e *PageCaption) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Text = RichText(buf.PopObj())
	e.Credit = RichText(buf.PopObj())
}

type InputContact struct {
	ClientId  int64  `validate:"required"`
	Phone     string `validate:"required"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
}

func (e *InputContact) CRC() uint32 {
	return uint32(0xf392b7f4)
}
func (e *InputContact) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutLong(e.ClientId)
	buf.PutString(e.Phone)
	buf.PutString(e.FirstName)
	buf.PutString(e.LastName)
	return buf.Result()
}

func (e *InputContact) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.ClientId = buf.PopLong()
	e.Phone = buf.PopString()
	e.FirstName = buf.PopString()
	e.LastName = buf.PopString()
}

type MessagesHighScores struct {
	Scores []*HighScore `validate:"required"`
	Users  []User       `validate:"required"`
}

func (e *MessagesHighScores) CRC() uint32 {
	return uint32(0x9a3bfd99)
}
func (e *MessagesHighScores) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutVector(e.Scores)
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *MessagesHighScores) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Scores = buf.PopVector(reflect.TypeOf(*HighScore{})).([]*HighScore)
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type AccountPasswordSettings struct {
	__flagsPosition struct{}              // flags param position `validate:"required"`
	Email           string                `flag:"0"`
	SecureSettings  *SecureSecretSettings `flag:"1"`
}

func (e *AccountPasswordSettings) CRC() uint32 {
	return uint32(0x9a5c33e5)
}
func (e *AccountPasswordSettings) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Email) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.SecureSettings) {
		flag |= 1 << 1
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Email) {
		buf.PutString(e.Email)
	}
	if !zero.IsZeroVal(e.SecureSettings) {
		buf.PutRawBytes(e.SecureSettings.Encode())
	}
	return buf.Result()
}

func (e *AccountPasswordSettings) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Email = buf.PopString()
	}
	if flags&1<<1 > 0 {
		e.SecureSettings = buf.PopObj().(*SecureSecretSettings)
	}
}

type AuthPasswordRecovery struct {
	EmailPattern string `validate:"required"`
}

func (e *AuthPasswordRecovery) CRC() uint32 {
	return uint32(0x137948a5)
}
func (e *AuthPasswordRecovery) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.EmailPattern)
	return buf.Result()
}

func (e *AuthPasswordRecovery) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.EmailPattern = buf.PopString()
}

type MessagesPeerDialogs struct {
	Dialogs  []Dialog      `validate:"required"`
	Messages []Message     `validate:"required"`
	Chats    []Chat        `validate:"required"`
	Users    []User        `validate:"required"`
	State    *UpdatesState `validate:"required"`
}

func (e *MessagesPeerDialogs) CRC() uint32 {
	return uint32(0x3371c354)
}
func (e *MessagesPeerDialogs) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutVector(e.Dialogs)
	buf.PutVector(e.Messages)
	buf.PutVector(e.Chats)
	buf.PutVector(e.Users)
	buf.PutRawBytes(e.State.Encode())
	return buf.Result()
}

func (e *MessagesPeerDialogs) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Dialogs = buf.PopVector(reflect.TypeOf(Dialog{})).([]Dialog)
	e.Messages = buf.PopVector(reflect.TypeOf(Message{})).([]Message)
	e.Chats = buf.PopVector(reflect.TypeOf(Chat{})).([]Chat)
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
	e.State = buf.PopObj().(*UpdatesState)
}

type ContactsFound struct {
	MyResults []Peer `validate:"required"`
	Results   []Peer `validate:"required"`
	Chats     []Chat `validate:"required"`
	Users     []User `validate:"required"`
}

func (e *ContactsFound) CRC() uint32 {
	return uint32(0xb3134d9d)
}
func (e *ContactsFound) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutVector(e.MyResults)
	buf.PutVector(e.Results)
	buf.PutVector(e.Chats)
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *ContactsFound) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.MyResults = buf.PopVector(reflect.TypeOf(Peer{})).([]Peer)
	e.Results = buf.PopVector(reflect.TypeOf(Peer{})).([]Peer)
	e.Chats = buf.PopVector(reflect.TypeOf(Chat{})).([]Chat)
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type AccountDaysTTL struct {
	Days int32 `validate:"required"`
}

func (e *AccountDaysTTL) CRC() uint32 {
	return uint32(0xb8d0afdf)
}
func (e *AccountDaysTTL) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.Days)
	return buf.Result()
}

func (e *AccountDaysTTL) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Days = buf.PopInt()
}

type WebAuthorization struct {
	Hash        int64  `validate:"required"`
	BotId       int32  `validate:"required"`
	Domain      string `validate:"required"`
	Browser     string `validate:"required"`
	Platform    string `validate:"required"`
	DateCreated int32  `validate:"required"`
	DateActive  int32  `validate:"required"`
	Ip          string `validate:"required"`
	Region      string `validate:"required"`
}

func (e *WebAuthorization) CRC() uint32 {
	return uint32(0xcac943f2)
}
func (e *WebAuthorization) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutLong(e.Hash)
	buf.PutInt(e.BotId)
	buf.PutString(e.Domain)
	buf.PutString(e.Browser)
	buf.PutString(e.Platform)
	buf.PutInt(e.DateCreated)
	buf.PutInt(e.DateActive)
	buf.PutString(e.Ip)
	buf.PutString(e.Region)
	return buf.Result()
}

func (e *WebAuthorization) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Hash = buf.PopLong()
	e.BotId = buf.PopInt()
	e.Domain = buf.PopString()
	e.Browser = buf.PopString()
	e.Platform = buf.PopString()
	e.DateCreated = buf.PopInt()
	e.DateActive = buf.PopInt()
	e.Ip = buf.PopString()
	e.Region = buf.PopString()
}

type StatsGroupTopInviter struct {
	UserId      int32 `validate:"required"`
	Invitations int32 `validate:"required"`
}

func (e *StatsGroupTopInviter) CRC() uint32 {
	return uint32(0x31962a4c)
}
func (e *StatsGroupTopInviter) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.UserId)
	buf.PutInt(e.Invitations)
	return buf.Result()
}

func (e *StatsGroupTopInviter) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.UserId = buf.PopInt()
	e.Invitations = buf.PopInt()
}

type ContactStatus struct {
	UserId int32      `validate:"required"`
	Status UserStatus `validate:"required"`
}

func (e *ContactStatus) CRC() uint32 {
	return uint32(0xd3680c61)
}
func (e *ContactStatus) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.UserId)
	buf.PutRawBytes(e.Status.Encode())
	return buf.Result()
}

func (e *ContactStatus) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.UserId = buf.PopInt()
	e.Status = UserStatus(buf.PopObj())
}

type MessagesArchivedStickers struct {
	Count int32               `validate:"required"`
	Sets  []StickerSetCovered `validate:"required"`
}

func (e *MessagesArchivedStickers) CRC() uint32 {
	return uint32(0x4fcba9c8)
}
func (e *MessagesArchivedStickers) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.Count)
	buf.PutVector(e.Sets)
	return buf.Result()
}

func (e *MessagesArchivedStickers) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Count = buf.PopInt()
	e.Sets = buf.PopVector(reflect.TypeOf(StickerSetCovered{})).([]StickerSetCovered)
}

type InputThemeSettings struct {
	__flagsPosition    struct{}           // flags param position `validate:"required"`
	BaseTheme          BaseTheme          `validate:"required"`
	AccentColor        int32              `validate:"required"`
	MessageTopColor    int32              `flag:"0"`
	MessageBottomColor int32              `flag:"0"`
	Wallpaper          InputWallPaper     `flag:"1"`
	WallpaperSettings  *WallPaperSettings `flag:"1"`
}

func (e *InputThemeSettings) CRC() uint32 {
	return uint32(0xbd507cd1)
}
func (e *InputThemeSettings) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.MessageTopColor) || !zero.IsZeroVal(e.MessageBottomColor) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Wallpaper) || !zero.IsZeroVal(e.WallpaperSettings) {
		flag |= 1 << 1
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.BaseTheme.Encode())
	buf.PutInt(e.AccentColor)
	if !zero.IsZeroVal(e.MessageTopColor) {
		buf.PutInt(e.MessageTopColor)
	}
	if !zero.IsZeroVal(e.MessageBottomColor) {
		buf.PutInt(e.MessageBottomColor)
	}
	if !zero.IsZeroVal(e.Wallpaper) {
		buf.PutRawBytes(e.Wallpaper.Encode())
	}
	if !zero.IsZeroVal(e.WallpaperSettings) {
		buf.PutRawBytes(e.WallpaperSettings.Encode())
	}
	return buf.Result()
}

func (e *InputThemeSettings) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	e.BaseTheme = *(buf.PopObj().(*BaseTheme))
	e.AccentColor = buf.PopInt()
	if flags&1<<0 > 0 {
		e.MessageTopColor = buf.PopInt()
	}
	if flags&1<<0 > 0 {
		e.MessageBottomColor = buf.PopInt()
	}
	if flags&1<<1 > 0 {
		e.Wallpaper = InputWallPaper(buf.PopObj())
	}
	if flags&1<<1 > 0 {
		e.WallpaperSettings = buf.PopObj().(*WallPaperSettings)
	}
}

type MessagesVotesList struct {
	__flagsPosition struct{}          // flags param position `validate:"required"`
	Count           int32             `validate:"required"`
	Votes           []MessageUserVote `validate:"required"`
	Users           []User            `validate:"required"`
	NextOffset      string            `flag:"0"`
}

func (e *MessagesVotesList) CRC() uint32 {
	return uint32(0x823f649)
}
func (e *MessagesVotesList) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.NextOffset) {
		flag |= 1 << 0
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.Count)
	buf.PutVector(e.Votes)
	buf.PutVector(e.Users)
	if !zero.IsZeroVal(e.NextOffset) {
		buf.PutString(e.NextOffset)
	}
	return buf.Result()
}

func (e *MessagesVotesList) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	e.Count = buf.PopInt()
	e.Votes = buf.PopVector(reflect.TypeOf(MessageUserVote{})).([]MessageUserVote)
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
	if flags&1<<0 > 0 {
		e.NextOffset = buf.PopString()
	}
}

type SecureSecretSettings struct {
	SecureAlgo     SecurePasswordKdfAlgo `validate:"required"`
	SecureSecret   []byte                `validate:"required"`
	SecureSecretId int64                 `validate:"required"`
}

func (e *SecureSecretSettings) CRC() uint32 {
	return uint32(0x1527bcac)
}
func (e *SecureSecretSettings) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.SecureAlgo.Encode())
	buf.PutMessage(e.SecureSecret)
	buf.PutLong(e.SecureSecretId)
	return buf.Result()
}

func (e *SecureSecretSettings) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.SecureAlgo = SecurePasswordKdfAlgo(buf.PopObj())
	e.SecureSecret = buf.PopMessage()
	e.SecureSecretId = buf.PopLong()
}

type ChatBannedRights struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	ViewMessages    bool     `flag:"0,encoded_in_bitflags"`
	SendMessages    bool     `flag:"1,encoded_in_bitflags"`
	SendMedia       bool     `flag:"2,encoded_in_bitflags"`
	SendStickers    bool     `flag:"3,encoded_in_bitflags"`
	SendGifs        bool     `flag:"4,encoded_in_bitflags"`
	SendGames       bool     `flag:"5,encoded_in_bitflags"`
	SendInline      bool     `flag:"6,encoded_in_bitflags"`
	EmbedLinks      bool     `flag:"7,encoded_in_bitflags"`
	SendPolls       bool     `flag:"8,encoded_in_bitflags"`
	ChangeInfo      bool     `flag:"10,encoded_in_bitflags"`
	InviteUsers     bool     `flag:"15,encoded_in_bitflags"`
	PinMessages     bool     `flag:"17,encoded_in_bitflags"`
	UntilDate       int32    `validate:"required"`
}

func (e *ChatBannedRights) CRC() uint32 {
	return uint32(0x9f120418)
}
func (e *ChatBannedRights) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.ViewMessages) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.SendMessages) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.SendMedia) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.SendStickers) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.SendGifs) {
		flag |= 1 << 4
	}
	if !zero.IsZeroVal(e.SendGames) {
		flag |= 1 << 5
	}
	if !zero.IsZeroVal(e.SendInline) {
		flag |= 1 << 6
	}
	if !zero.IsZeroVal(e.EmbedLinks) {
		flag |= 1 << 7
	}
	if !zero.IsZeroVal(e.SendPolls) {
		flag |= 1 << 8
	}
	if !zero.IsZeroVal(e.ChangeInfo) {
		flag |= 1 << 10
	}
	if !zero.IsZeroVal(e.InviteUsers) {
		flag |= 1 << 15
	}
	if !zero.IsZeroVal(e.PinMessages) {
		flag |= 1 << 17
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.ViewMessages) {
	}
	if !zero.IsZeroVal(e.SendMessages) {
	}
	if !zero.IsZeroVal(e.SendMedia) {
	}
	if !zero.IsZeroVal(e.SendStickers) {
	}
	if !zero.IsZeroVal(e.SendGifs) {
	}
	if !zero.IsZeroVal(e.SendGames) {
	}
	if !zero.IsZeroVal(e.SendInline) {
	}
	if !zero.IsZeroVal(e.EmbedLinks) {
	}
	if !zero.IsZeroVal(e.SendPolls) {
	}
	if !zero.IsZeroVal(e.ChangeInfo) {
	}
	if !zero.IsZeroVal(e.InviteUsers) {
	}
	if !zero.IsZeroVal(e.PinMessages) {
	}
	buf.PutInt(e.UntilDate)
	return buf.Result()
}

func (e *ChatBannedRights) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.ViewMessages = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.SendMessages = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.SendMedia = buf.PopBool()
	}
	if flags&1<<3 > 0 {
		e.SendStickers = buf.PopBool()
	}
	if flags&1<<4 > 0 {
		e.SendGifs = buf.PopBool()
	}
	if flags&1<<5 > 0 {
		e.SendGames = buf.PopBool()
	}
	if flags&1<<6 > 0 {
		e.SendInline = buf.PopBool()
	}
	if flags&1<<7 > 0 {
		e.EmbedLinks = buf.PopBool()
	}
	if flags&1<<8 > 0 {
		e.SendPolls = buf.PopBool()
	}
	if flags&1<<10 > 0 {
		e.ChangeInfo = buf.PopBool()
	}
	if flags&1<<15 > 0 {
		e.InviteUsers = buf.PopBool()
	}
	if flags&1<<17 > 0 {
		e.PinMessages = buf.PopBool()
	}
	e.UntilDate = buf.PopInt()
}

type AccountPasswordInputSettings struct {
	__flagsPosition   struct{}              // flags param position `validate:"required"`
	NewAlgo           PasswordKdfAlgo       `flag:"0"`
	NewPasswordHash   []byte                `flag:"0"`
	Hint              string                `flag:"0"`
	Email             string                `flag:"1"`
	NewSecureSettings *SecureSecretSettings `flag:"2"`
}

func (e *AccountPasswordInputSettings) CRC() uint32 {
	return uint32(0xc23727c9)
}
func (e *AccountPasswordInputSettings) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.NewAlgo) || !zero.IsZeroVal(e.NewPasswordHash) || !zero.IsZeroVal(e.Hint) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Email) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.NewSecureSettings) {
		flag |= 1 << 2
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.NewAlgo) {
		buf.PutRawBytes(e.NewAlgo.Encode())
	}
	if !zero.IsZeroVal(e.NewPasswordHash) {
		buf.PutMessage(e.NewPasswordHash)
	}
	if !zero.IsZeroVal(e.Hint) {
		buf.PutString(e.Hint)
	}
	if !zero.IsZeroVal(e.Email) {
		buf.PutString(e.Email)
	}
	if !zero.IsZeroVal(e.NewSecureSettings) {
		buf.PutRawBytes(e.NewSecureSettings.Encode())
	}
	return buf.Result()
}

func (e *AccountPasswordInputSettings) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.NewAlgo = PasswordKdfAlgo(buf.PopObj())
	}
	if flags&1<<0 > 0 {
		e.NewPasswordHash = buf.PopMessage()
	}
	if flags&1<<0 > 0 {
		e.Hint = buf.PopString()
	}
	if flags&1<<1 > 0 {
		e.Email = buf.PopString()
	}
	if flags&1<<2 > 0 {
		e.NewSecureSettings = buf.PopObj().(*SecureSecretSettings)
	}
}

type MessageRange struct {
	MinId int32 `validate:"required"`
	MaxId int32 `validate:"required"`
}

func (e *MessageRange) CRC() uint32 {
	return uint32(0xae30253)
}
func (e *MessageRange) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.MinId)
	buf.PutInt(e.MaxId)
	return buf.Result()
}

func (e *MessageRange) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.MinId = buf.PopInt()
	e.MaxId = buf.PopInt()
}

type PollResults struct {
	__flagsPosition  struct{}            // flags param position `validate:"required"`
	Min              bool                `flag:"0,encoded_in_bitflags"`
	Results          []*PollAnswerVoters `flag:"1"`
	TotalVoters      int32               `flag:"2"`
	RecentVoters     []int32             `flag:"3"`
	Solution         string              `flag:"4"`
	SolutionEntities []MessageEntity     `flag:"4"`
}

func (e *PollResults) CRC() uint32 {
	return uint32(0xbadcc1a3)
}
func (e *PollResults) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Min) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Results) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.TotalVoters) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.RecentVoters) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.Solution) || !zero.IsZeroVal(e.SolutionEntities) {
		flag |= 1 << 4
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Min) {
	}
	if !zero.IsZeroVal(e.Results) {
		buf.PutVector(e.Results)
	}
	if !zero.IsZeroVal(e.TotalVoters) {
		buf.PutInt(e.TotalVoters)
	}
	if !zero.IsZeroVal(e.RecentVoters) {
		buf.PutVector(e.RecentVoters)
	}
	if !zero.IsZeroVal(e.Solution) {
		buf.PutString(e.Solution)
	}
	if !zero.IsZeroVal(e.SolutionEntities) {
		buf.PutVector(e.SolutionEntities)
	}
	return buf.Result()
}

func (e *PollResults) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Min = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.Results = buf.PopVector(reflect.TypeOf(*PollAnswerVoters{})).([]*PollAnswerVoters)
	}
	if flags&1<<2 > 0 {
		e.TotalVoters = buf.PopInt()
	}
	if flags&1<<3 > 0 {
		e.RecentVoters = buf.PopVector(reflect.TypeOf(int32{})).([]int32)
	}
	if flags&1<<4 > 0 {
		e.Solution = buf.PopString()
	}
	if flags&1<<4 > 0 {
		e.SolutionEntities = buf.PopVector(reflect.TypeOf(MessageEntity{})).([]MessageEntity)
	}
}

type MessagesAffectedMessages struct {
	Pts      int32 `validate:"required"`
	PtsCount int32 `validate:"required"`
}

func (e *MessagesAffectedMessages) CRC() uint32 {
	return uint32(0x84d19185)
}
func (e *MessagesAffectedMessages) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.Pts)
	buf.PutInt(e.PtsCount)
	return buf.Result()
}

func (e *MessagesAffectedMessages) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Pts = buf.PopInt()
	e.PtsCount = buf.PopInt()
}

type PaymentsPaymentReceipt struct {
	__flagsPosition  struct{}              // flags param position `validate:"required"`
	Date             int32                 `validate:"required"`
	BotId            int32                 `validate:"required"`
	Invoice          *Invoice              `validate:"required"`
	ProviderId       int32                 `validate:"required"`
	Info             *PaymentRequestedInfo `flag:"0"`
	Shipping         *ShippingOption       `flag:"1"`
	Currency         string                `validate:"required"`
	TotalAmount      int64                 `validate:"required"`
	CredentialsTitle string                `validate:"required"`
	Users            []User                `validate:"required"`
}

func (e *PaymentsPaymentReceipt) CRC() uint32 {
	return uint32(0x500911e1)
}
func (e *PaymentsPaymentReceipt) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Info) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Shipping) {
		flag |= 1 << 1
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.Date)
	buf.PutInt(e.BotId)
	buf.PutRawBytes(e.Invoice.Encode())
	buf.PutInt(e.ProviderId)
	if !zero.IsZeroVal(e.Info) {
		buf.PutRawBytes(e.Info.Encode())
	}
	if !zero.IsZeroVal(e.Shipping) {
		buf.PutRawBytes(e.Shipping.Encode())
	}
	buf.PutString(e.Currency)
	buf.PutLong(e.TotalAmount)
	buf.PutString(e.CredentialsTitle)
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *PaymentsPaymentReceipt) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	e.Date = buf.PopInt()
	e.BotId = buf.PopInt()
	e.Invoice = buf.PopObj().(*Invoice)
	e.ProviderId = buf.PopInt()
	if flags&1<<0 > 0 {
		e.Info = buf.PopObj().(*PaymentRequestedInfo)
	}
	if flags&1<<1 > 0 {
		e.Shipping = buf.PopObj().(*ShippingOption)
	}
	e.Currency = buf.PopString()
	e.TotalAmount = buf.PopLong()
	e.CredentialsTitle = buf.PopString()
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type ChannelAdminLogEventsFilter struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	Join            bool     `flag:"0,encoded_in_bitflags"`
	Leave           bool     `flag:"1,encoded_in_bitflags"`
	Invite          bool     `flag:"2,encoded_in_bitflags"`
	Ban             bool     `flag:"3,encoded_in_bitflags"`
	Unban           bool     `flag:"4,encoded_in_bitflags"`
	Kick            bool     `flag:"5,encoded_in_bitflags"`
	Unkick          bool     `flag:"6,encoded_in_bitflags"`
	Promote         bool     `flag:"7,encoded_in_bitflags"`
	Demote          bool     `flag:"8,encoded_in_bitflags"`
	Info            bool     `flag:"9,encoded_in_bitflags"`
	Settings        bool     `flag:"10,encoded_in_bitflags"`
	Pinned          bool     `flag:"11,encoded_in_bitflags"`
	Edit            bool     `flag:"12,encoded_in_bitflags"`
	Delete          bool     `flag:"13,encoded_in_bitflags"`
}

func (e *ChannelAdminLogEventsFilter) CRC() uint32 {
	return uint32(0xea107ae4)
}
func (e *ChannelAdminLogEventsFilter) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Join) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Leave) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.Invite) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.Ban) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.Unban) {
		flag |= 1 << 4
	}
	if !zero.IsZeroVal(e.Kick) {
		flag |= 1 << 5
	}
	if !zero.IsZeroVal(e.Unkick) {
		flag |= 1 << 6
	}
	if !zero.IsZeroVal(e.Promote) {
		flag |= 1 << 7
	}
	if !zero.IsZeroVal(e.Demote) {
		flag |= 1 << 8
	}
	if !zero.IsZeroVal(e.Info) {
		flag |= 1 << 9
	}
	if !zero.IsZeroVal(e.Settings) {
		flag |= 1 << 10
	}
	if !zero.IsZeroVal(e.Pinned) {
		flag |= 1 << 11
	}
	if !zero.IsZeroVal(e.Edit) {
		flag |= 1 << 12
	}
	if !zero.IsZeroVal(e.Delete) {
		flag |= 1 << 13
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Join) {
	}
	if !zero.IsZeroVal(e.Leave) {
	}
	if !zero.IsZeroVal(e.Invite) {
	}
	if !zero.IsZeroVal(e.Ban) {
	}
	if !zero.IsZeroVal(e.Unban) {
	}
	if !zero.IsZeroVal(e.Kick) {
	}
	if !zero.IsZeroVal(e.Unkick) {
	}
	if !zero.IsZeroVal(e.Promote) {
	}
	if !zero.IsZeroVal(e.Demote) {
	}
	if !zero.IsZeroVal(e.Info) {
	}
	if !zero.IsZeroVal(e.Settings) {
	}
	if !zero.IsZeroVal(e.Pinned) {
	}
	if !zero.IsZeroVal(e.Edit) {
	}
	if !zero.IsZeroVal(e.Delete) {
	}
	return buf.Result()
}

func (e *ChannelAdminLogEventsFilter) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Join = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.Leave = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.Invite = buf.PopBool()
	}
	if flags&1<<3 > 0 {
		e.Ban = buf.PopBool()
	}
	if flags&1<<4 > 0 {
		e.Unban = buf.PopBool()
	}
	if flags&1<<5 > 0 {
		e.Kick = buf.PopBool()
	}
	if flags&1<<6 > 0 {
		e.Unkick = buf.PopBool()
	}
	if flags&1<<7 > 0 {
		e.Promote = buf.PopBool()
	}
	if flags&1<<8 > 0 {
		e.Demote = buf.PopBool()
	}
	if flags&1<<9 > 0 {
		e.Info = buf.PopBool()
	}
	if flags&1<<10 > 0 {
		e.Settings = buf.PopBool()
	}
	if flags&1<<11 > 0 {
		e.Pinned = buf.PopBool()
	}
	if flags&1<<12 > 0 {
		e.Edit = buf.PopBool()
	}
	if flags&1<<13 > 0 {
		e.Delete = buf.PopBool()
	}
}

type Contact struct {
	UserId int32 `validate:"required"`
	Mutual bool  `validate:"required"`
}

func (e *Contact) CRC() uint32 {
	return uint32(0xf911c994)
}
func (e *Contact) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.UserId)
	buf.PutBool(e.Mutual)
	return buf.Result()
}

func (e *Contact) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.UserId = buf.PopInt()
	e.Mutual = buf.PopBool()
}

type StatsDateRangeDays struct {
	MinDate int32 `validate:"required"`
	MaxDate int32 `validate:"required"`
}

func (e *StatsDateRangeDays) CRC() uint32 {
	return uint32(0xb637edaf)
}
func (e *StatsDateRangeDays) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.MinDate)
	buf.PutInt(e.MaxDate)
	return buf.Result()
}

func (e *StatsDateRangeDays) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.MinDate = buf.PopInt()
	e.MaxDate = buf.PopInt()
}

type InputPeerNotifySettings struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	ShowPreviews    bool     `flag:"0"`
	Silent          bool     `flag:"1"`
	MuteUntil       int32    `flag:"2"`
	Sound           string   `flag:"3"`
}

func (e *InputPeerNotifySettings) CRC() uint32 {
	return uint32(0x9c3d198e)
}
func (e *InputPeerNotifySettings) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.ShowPreviews) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Silent) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.MuteUntil) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.Sound) {
		flag |= 1 << 3
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.ShowPreviews) {
		buf.PutBool(e.ShowPreviews)
	}
	if !zero.IsZeroVal(e.Silent) {
		buf.PutBool(e.Silent)
	}
	if !zero.IsZeroVal(e.MuteUntil) {
		buf.PutInt(e.MuteUntil)
	}
	if !zero.IsZeroVal(e.Sound) {
		buf.PutString(e.Sound)
	}
	return buf.Result()
}

func (e *InputPeerNotifySettings) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.ShowPreviews = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.Silent = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.MuteUntil = buf.PopInt()
	}
	if flags&1<<3 > 0 {
		e.Sound = buf.PopString()
	}
}

type SecureValueHash struct {
	Type SecureValueType `validate:"required"`
	Hash []byte          `validate:"required"`
}

func (e *SecureValueHash) CRC() uint32 {
	return uint32(0xed1ecdb0)
}
func (e *SecureValueHash) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Type.Encode())
	buf.PutMessage(e.Hash)
	return buf.Result()
}

func (e *SecureValueHash) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Type = *(buf.PopObj().(*SecureValueType))
	e.Hash = buf.PopMessage()
}

type SecureValue struct {
	__flagsPosition struct{}        // flags param position `validate:"required"`
	Type            SecureValueType `validate:"required"`
	Data            *SecureData     `flag:"0"`
	FrontSide       SecureFile      `flag:"1"`
	ReverseSide     SecureFile      `flag:"2"`
	Selfie          SecureFile      `flag:"3"`
	Translation     []SecureFile    `flag:"6"`
	Files           []SecureFile    `flag:"4"`
	PlainData       SecurePlainData `flag:"5"`
	Hash            []byte          `validate:"required"`
}

func (e *SecureValue) CRC() uint32 {
	return uint32(0x187fa0ca)
}
func (e *SecureValue) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Data) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.FrontSide) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.ReverseSide) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.Selfie) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.Files) {
		flag |= 1 << 4
	}
	if !zero.IsZeroVal(e.PlainData) {
		flag |= 1 << 5
	}
	if !zero.IsZeroVal(e.Translation) {
		flag |= 1 << 6
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Type.Encode())
	if !zero.IsZeroVal(e.Data) {
		buf.PutRawBytes(e.Data.Encode())
	}
	if !zero.IsZeroVal(e.FrontSide) {
		buf.PutRawBytes(e.FrontSide.Encode())
	}
	if !zero.IsZeroVal(e.ReverseSide) {
		buf.PutRawBytes(e.ReverseSide.Encode())
	}
	if !zero.IsZeroVal(e.Selfie) {
		buf.PutRawBytes(e.Selfie.Encode())
	}
	if !zero.IsZeroVal(e.Translation) {
		buf.PutVector(e.Translation)
	}
	if !zero.IsZeroVal(e.Files) {
		buf.PutVector(e.Files)
	}
	if !zero.IsZeroVal(e.PlainData) {
		buf.PutRawBytes(e.PlainData.Encode())
	}
	buf.PutMessage(e.Hash)
	return buf.Result()
}

func (e *SecureValue) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	e.Type = *(buf.PopObj().(*SecureValueType))
	if flags&1<<0 > 0 {
		e.Data = buf.PopObj().(*SecureData)
	}
	if flags&1<<1 > 0 {
		e.FrontSide = SecureFile(buf.PopObj())
	}
	if flags&1<<2 > 0 {
		e.ReverseSide = SecureFile(buf.PopObj())
	}
	if flags&1<<3 > 0 {
		e.Selfie = SecureFile(buf.PopObj())
	}
	if flags&1<<6 > 0 {
		e.Translation = buf.PopVector(reflect.TypeOf(SecureFile{})).([]SecureFile)
	}
	if flags&1<<4 > 0 {
		e.Files = buf.PopVector(reflect.TypeOf(SecureFile{})).([]SecureFile)
	}
	if flags&1<<5 > 0 {
		e.PlainData = SecurePlainData(buf.PopObj())
	}
	e.Hash = buf.PopMessage()
}

type AuthExportedAuthorization struct {
	Id    int32  `validate:"required"`
	Bytes []byte `validate:"required"`
}

func (e *AuthExportedAuthorization) CRC() uint32 {
	return uint32(0xdf969c2d)
}
func (e *AuthExportedAuthorization) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.Id)
	buf.PutMessage(e.Bytes)
	return buf.Result()
}

func (e *AuthExportedAuthorization) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Id = buf.PopInt()
	e.Bytes = buf.PopMessage()
}

type InputClientProxy struct {
	Address string `validate:"required"`
	Port    int32  `validate:"required"`
}

func (e *InputClientProxy) CRC() uint32 {
	return uint32(0x75588b3f)
}
func (e *InputClientProxy) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Address)
	buf.PutInt(e.Port)
	return buf.Result()
}

func (e *InputClientProxy) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Address = buf.PopString()
	e.Port = buf.PopInt()
}

type InlineBotSwitchPM struct {
	Text       string `validate:"required"`
	StartParam string `validate:"required"`
}

func (e *InlineBotSwitchPM) CRC() uint32 {
	return uint32(0x3c20629f)
}
func (e *InlineBotSwitchPM) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Text)
	buf.PutString(e.StartParam)
	return buf.Result()
}

func (e *InlineBotSwitchPM) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Text = buf.PopString()
	e.StartParam = buf.PopString()
}

type PaymentsValidatedRequestedInfo struct {
	__flagsPosition struct{}          // flags param position `validate:"required"`
	Id              string            `flag:"0"`
	ShippingOptions []*ShippingOption `flag:"1"`
}

func (e *PaymentsValidatedRequestedInfo) CRC() uint32 {
	return uint32(0xd1451883)
}
func (e *PaymentsValidatedRequestedInfo) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Id) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.ShippingOptions) {
		flag |= 1 << 1
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Id) {
		buf.PutString(e.Id)
	}
	if !zero.IsZeroVal(e.ShippingOptions) {
		buf.PutVector(e.ShippingOptions)
	}
	return buf.Result()
}

func (e *PaymentsValidatedRequestedInfo) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Id = buf.PopString()
	}
	if flags&1<<1 > 0 {
		e.ShippingOptions = buf.PopVector(reflect.TypeOf(*ShippingOption{})).([]*ShippingOption)
	}
}

type ContactsImportedContacts struct {
	Imported       []*ImportedContact `validate:"required"`
	PopularInvites []*PopularContact  `validate:"required"`
	RetryContacts  []int64            `validate:"required"`
	Users          []User             `validate:"required"`
}

func (e *ContactsImportedContacts) CRC() uint32 {
	return uint32(0x77d01c3b)
}
func (e *ContactsImportedContacts) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutVector(e.Imported)
	buf.PutVector(e.PopularInvites)
	buf.PutVector(e.RetryContacts)
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *ContactsImportedContacts) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Imported = buf.PopVector(reflect.TypeOf(*ImportedContact{})).([]*ImportedContact)
	e.PopularInvites = buf.PopVector(reflect.TypeOf(*PopularContact{})).([]*PopularContact)
	e.RetryContacts = buf.PopVector(reflect.TypeOf(int64{})).([]int64)
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type AccountSentEmailCode struct {
	EmailPattern string `validate:"required"`
	Length       int32  `validate:"required"`
}

func (e *AccountSentEmailCode) CRC() uint32 {
	return uint32(0x811f854f)
}
func (e *AccountSentEmailCode) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.EmailPattern)
	buf.PutInt(e.Length)
	return buf.Result()
}

func (e *AccountSentEmailCode) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.EmailPattern = buf.PopString()
	e.Length = buf.PopInt()
}

type ContactsResolvedPeer struct {
	Peer  Peer   `validate:"required"`
	Chats []Chat `validate:"required"`
	Users []User `validate:"required"`
}

func (e *ContactsResolvedPeer) CRC() uint32 {
	return uint32(0x7f077ad9)
}
func (e *ContactsResolvedPeer) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Peer.Encode())
	buf.PutVector(e.Chats)
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *ContactsResolvedPeer) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Peer = Peer(buf.PopObj())
	e.Chats = buf.PopVector(reflect.TypeOf(Chat{})).([]Chat)
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type Page struct {
	__flagsPosition struct{}    // flags param position `validate:"required"`
	Part            bool        `flag:"0,encoded_in_bitflags"`
	Rtl             bool        `flag:"1,encoded_in_bitflags"`
	V2              bool        `flag:"2,encoded_in_bitflags"`
	Url             string      `validate:"required"`
	Blocks          []PageBlock `validate:"required"`
	Photos          []Photo     `validate:"required"`
	Documents       []Document  `validate:"required"`
	Views           int32       `flag:"3"`
}

func (e *Page) CRC() uint32 {
	return uint32(0x98657f0d)
}
func (e *Page) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Part) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Rtl) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.V2) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.Views) {
		flag |= 1 << 3
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Part) {
	}
	if !zero.IsZeroVal(e.Rtl) {
	}
	if !zero.IsZeroVal(e.V2) {
	}
	buf.PutString(e.Url)
	buf.PutVector(e.Blocks)
	buf.PutVector(e.Photos)
	buf.PutVector(e.Documents)
	if !zero.IsZeroVal(e.Views) {
		buf.PutInt(e.Views)
	}
	return buf.Result()
}

func (e *Page) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Part = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.Rtl = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.V2 = buf.PopBool()
	}
	e.Url = buf.PopString()
	e.Blocks = buf.PopVector(reflect.TypeOf(PageBlock{})).([]PageBlock)
	e.Photos = buf.PopVector(reflect.TypeOf(Photo{})).([]Photo)
	e.Documents = buf.PopVector(reflect.TypeOf(Document{})).([]Document)
	if flags&1<<3 > 0 {
		e.Views = buf.PopInt()
	}
}

type NearestDc struct {
	Country   string `validate:"required"`
	ThisDc    int32  `validate:"required"`
	NearestDc int32  `validate:"required"`
}

func (e *NearestDc) CRC() uint32 {
	return uint32(0x8e1a1775)
}
func (e *NearestDc) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Country)
	buf.PutInt(e.ThisDc)
	buf.PutInt(e.NearestDc)
	return buf.Result()
}

func (e *NearestDc) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Country = buf.PopString()
	e.ThisDc = buf.PopInt()
	e.NearestDc = buf.PopInt()
}

type PhoneCallProtocol struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	UdpP2P          bool     `flag:"0,encoded_in_bitflags"`
	UdpReflector    bool     `flag:"1,encoded_in_bitflags"`
	MinLayer        int32    `validate:"required"`
	MaxLayer        int32    `validate:"required"`
	LibraryVersions []string `validate:"required"`
}

func (e *PhoneCallProtocol) CRC() uint32 {
	return uint32(0xfc878fc8)
}
func (e *PhoneCallProtocol) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.UdpP2P) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.UdpReflector) {
		flag |= 1 << 1
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.UdpP2P) {
	}
	if !zero.IsZeroVal(e.UdpReflector) {
	}
	buf.PutInt(e.MinLayer)
	buf.PutInt(e.MaxLayer)
	buf.PutVector(e.LibraryVersions)
	return buf.Result()
}

func (e *PhoneCallProtocol) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.UdpP2P = buf.PopBool()
	}
	if flags&1<<1 > 0 {
		e.UdpReflector = buf.PopBool()
	}
	e.MinLayer = buf.PopInt()
	e.MaxLayer = buf.PopInt()
	e.LibraryVersions = buf.PopVector(reflect.TypeOf(string{})).([]string)
}

type PhonePhoneCall struct {
	PhoneCall PhoneCall `validate:"required"`
	Users     []User    `validate:"required"`
}

func (e *PhonePhoneCall) CRC() uint32 {
	return uint32(0xec82e140)
}
func (e *PhonePhoneCall) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.PhoneCall.Encode())
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *PhonePhoneCall) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.PhoneCall = PhoneCall(buf.PopObj())
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type HelpSupportName struct {
	Name string `validate:"required"`
}

func (e *HelpSupportName) CRC() uint32 {
	return uint32(0x8c05f1c9)
}
func (e *HelpSupportName) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Name)
	return buf.Result()
}

func (e *HelpSupportName) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Name = buf.PopString()
}

type MessageInteractionCounters struct {
	MsgId    int32 `validate:"required"`
	Views    int32 `validate:"required"`
	Forwards int32 `validate:"required"`
}

func (e *MessageInteractionCounters) CRC() uint32 {
	return uint32(0xad4fc9bd)
}
func (e *MessageInteractionCounters) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.MsgId)
	buf.PutInt(e.Views)
	buf.PutInt(e.Forwards)
	return buf.Result()
}

func (e *MessageInteractionCounters) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.MsgId = buf.PopInt()
	e.Views = buf.PopInt()
	e.Forwards = buf.PopInt()
}

type HelpConfigSimple struct {
	Date    int32              `validate:"required"`
	Expires int32              `validate:"required"`
	Rules   []*AccessPointRule `validate:"required"`
}

func (e *HelpConfigSimple) CRC() uint32 {
	return uint32(0x5a592a6c)
}
func (e *HelpConfigSimple) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.Date)
	buf.PutInt(e.Expires)
	buf.PutVector(e.Rules)
	return buf.Result()
}

func (e *HelpConfigSimple) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Date = buf.PopInt()
	e.Expires = buf.PopInt()
	e.Rules = buf.PopVector(reflect.TypeOf(*AccessPointRule{})).([]*AccessPointRule)
}

type ChannelAdminLogEvent struct {
	Id     int64                      `validate:"required"`
	Date   int32                      `validate:"required"`
	UserId int32                      `validate:"required"`
	Action ChannelAdminLogEventAction `validate:"required"`
}

func (e *ChannelAdminLogEvent) CRC() uint32 {
	return uint32(0x3b5a3e40)
}
func (e *ChannelAdminLogEvent) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutLong(e.Id)
	buf.PutInt(e.Date)
	buf.PutInt(e.UserId)
	buf.PutRawBytes(e.Action.Encode())
	return buf.Result()
}

func (e *ChannelAdminLogEvent) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Id = buf.PopLong()
	e.Date = buf.PopInt()
	e.UserId = buf.PopInt()
	e.Action = ChannelAdminLogEventAction(buf.PopObj())
}

type AccountAutoDownloadSettings struct {
	Low    *AutoDownloadSettings `validate:"required"`
	Medium *AutoDownloadSettings `validate:"required"`
	High   *AutoDownloadSettings `validate:"required"`
}

func (e *AccountAutoDownloadSettings) CRC() uint32 {
	return uint32(0x63cacf26)
}
func (e *AccountAutoDownloadSettings) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Low.Encode())
	buf.PutRawBytes(e.Medium.Encode())
	buf.PutRawBytes(e.High.Encode())
	return buf.Result()
}

func (e *AccountAutoDownloadSettings) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Low = buf.PopObj().(*AutoDownloadSettings)
	e.Medium = buf.PopObj().(*AutoDownloadSettings)
	e.High = buf.PopObj().(*AutoDownloadSettings)
}

type HelpSupport struct {
	PhoneNumber string `validate:"required"`
	User        User   `validate:"required"`
}

func (e *HelpSupport) CRC() uint32 {
	return uint32(0x17c6b5f6)
}
func (e *HelpSupport) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.PhoneNumber)
	buf.PutRawBytes(e.User.Encode())
	return buf.Result()
}

func (e *HelpSupport) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.PhoneNumber = buf.PopString()
	e.User = User(buf.PopObj())
}

type PaymentRequestedInfo struct {
	__flagsPosition struct{}     // flags param position `validate:"required"`
	Name            string       `flag:"0"`
	Phone           string       `flag:"1"`
	Email           string       `flag:"2"`
	ShippingAddress *PostAddress `flag:"3"`
}

func (e *PaymentRequestedInfo) CRC() uint32 {
	return uint32(0x909c3f94)
}
func (e *PaymentRequestedInfo) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Name) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Phone) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.Email) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.ShippingAddress) {
		flag |= 1 << 3
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Name) {
		buf.PutString(e.Name)
	}
	if !zero.IsZeroVal(e.Phone) {
		buf.PutString(e.Phone)
	}
	if !zero.IsZeroVal(e.Email) {
		buf.PutString(e.Email)
	}
	if !zero.IsZeroVal(e.ShippingAddress) {
		buf.PutRawBytes(e.ShippingAddress.Encode())
	}
	return buf.Result()
}

func (e *PaymentRequestedInfo) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Name = buf.PopString()
	}
	if flags&1<<1 > 0 {
		e.Phone = buf.PopString()
	}
	if flags&1<<2 > 0 {
		e.Email = buf.PopString()
	}
	if flags&1<<3 > 0 {
		e.ShippingAddress = buf.PopObj().(*PostAddress)
	}
}

type UserFull struct {
	__flagsPosition     struct{}            // flags param position `validate:"required"`
	Blocked             bool                `flag:"0,encoded_in_bitflags"`
	PhoneCallsAvailable bool                `flag:"4,encoded_in_bitflags"`
	PhoneCallsPrivate   bool                `flag:"5,encoded_in_bitflags"`
	CanPinMessage       bool                `flag:"7,encoded_in_bitflags"`
	HasScheduled        bool                `flag:"12,encoded_in_bitflags"`
	VideoCallsAvailable bool                `flag:"13,encoded_in_bitflags"`
	User                User                `validate:"required"`
	About               string              `flag:"1"`
	Settings            *PeerSettings       `validate:"required"`
	ProfilePhoto        Photo               `flag:"2"`
	NotifySettings      *PeerNotifySettings `validate:"required"`
	BotInfo             *BotInfo            `flag:"3"`
	PinnedMsgId         int32               `flag:"6"`
	CommonChatsCount    int32               `validate:"required"`
	FolderId            int32               `flag:"11"`
}

func (e *UserFull) CRC() uint32 {
	return uint32(0xedf17c12)
}
func (e *UserFull) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Blocked) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.About) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.ProfilePhoto) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.BotInfo) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.PhoneCallsAvailable) {
		flag |= 1 << 4
	}
	if !zero.IsZeroVal(e.PhoneCallsPrivate) {
		flag |= 1 << 5
	}
	if !zero.IsZeroVal(e.PinnedMsgId) {
		flag |= 1 << 6
	}
	if !zero.IsZeroVal(e.CanPinMessage) {
		flag |= 1 << 7
	}
	if !zero.IsZeroVal(e.FolderId) {
		flag |= 1 << 11
	}
	if !zero.IsZeroVal(e.HasScheduled) {
		flag |= 1 << 12
	}
	if !zero.IsZeroVal(e.VideoCallsAvailable) {
		flag |= 1 << 13
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Blocked) {
	}
	if !zero.IsZeroVal(e.PhoneCallsAvailable) {
	}
	if !zero.IsZeroVal(e.PhoneCallsPrivate) {
	}
	if !zero.IsZeroVal(e.CanPinMessage) {
	}
	if !zero.IsZeroVal(e.HasScheduled) {
	}
	if !zero.IsZeroVal(e.VideoCallsAvailable) {
	}
	buf.PutRawBytes(e.User.Encode())
	if !zero.IsZeroVal(e.About) {
		buf.PutString(e.About)
	}
	buf.PutRawBytes(e.Settings.Encode())
	if !zero.IsZeroVal(e.ProfilePhoto) {
		buf.PutRawBytes(e.ProfilePhoto.Encode())
	}
	buf.PutRawBytes(e.NotifySettings.Encode())
	if !zero.IsZeroVal(e.BotInfo) {
		buf.PutRawBytes(e.BotInfo.Encode())
	}
	if !zero.IsZeroVal(e.PinnedMsgId) {
		buf.PutInt(e.PinnedMsgId)
	}
	buf.PutInt(e.CommonChatsCount)
	if !zero.IsZeroVal(e.FolderId) {
		buf.PutInt(e.FolderId)
	}
	return buf.Result()
}

func (e *UserFull) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Blocked = buf.PopBool()
	}
	if flags&1<<4 > 0 {
		e.PhoneCallsAvailable = buf.PopBool()
	}
	if flags&1<<5 > 0 {
		e.PhoneCallsPrivate = buf.PopBool()
	}
	if flags&1<<7 > 0 {
		e.CanPinMessage = buf.PopBool()
	}
	if flags&1<<12 > 0 {
		e.HasScheduled = buf.PopBool()
	}
	if flags&1<<13 > 0 {
		e.VideoCallsAvailable = buf.PopBool()
	}
	e.User = User(buf.PopObj())
	if flags&1<<1 > 0 {
		e.About = buf.PopString()
	}
	e.Settings = buf.PopObj().(*PeerSettings)
	if flags&1<<2 > 0 {
		e.ProfilePhoto = Photo(buf.PopObj())
	}
	e.NotifySettings = buf.PopObj().(*PeerNotifySettings)
	if flags&1<<3 > 0 {
		e.BotInfo = buf.PopObj().(*BotInfo)
	}
	if flags&1<<6 > 0 {
		e.PinnedMsgId = buf.PopInt()
	}
	e.CommonChatsCount = buf.PopInt()
	if flags&1<<11 > 0 {
		e.FolderId = buf.PopInt()
	}
}

type CdnPublicKey struct {
	DcId      int32  `validate:"required"`
	PublicKey string `validate:"required"`
}

func (e *CdnPublicKey) CRC() uint32 {
	return uint32(0xc982eaba)
}
func (e *CdnPublicKey) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.DcId)
	buf.PutString(e.PublicKey)
	return buf.Result()
}

func (e *CdnPublicKey) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.DcId = buf.PopInt()
	e.PublicKey = buf.PopString()
}

type DialogFilterSuggested struct {
	Filter      *DialogFilter `validate:"required"`
	Description string        `validate:"required"`
}

func (e *DialogFilterSuggested) CRC() uint32 {
	return uint32(0x77744d4a)
}
func (e *DialogFilterSuggested) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Filter.Encode())
	buf.PutString(e.Description)
	return buf.Result()
}

func (e *DialogFilterSuggested) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Filter = buf.PopObj().(*DialogFilter)
	e.Description = buf.PopString()
}

type LangPackLanguage struct {
	__flagsPosition struct{} // flags param position `validate:"required"`
	Official        bool     `flag:"0,encoded_in_bitflags"`
	Rtl             bool     `flag:"2,encoded_in_bitflags"`
	Beta            bool     `flag:"3,encoded_in_bitflags"`
	Name            string   `validate:"required"`
	NativeName      string   `validate:"required"`
	LangCode        string   `validate:"required"`
	BaseLangCode    string   `flag:"1"`
	PluralCode      string   `validate:"required"`
	StringsCount    int32    `validate:"required"`
	TranslatedCount int32    `validate:"required"`
	TranslationsUrl string   `validate:"required"`
}

func (e *LangPackLanguage) CRC() uint32 {
	return uint32(0xeeca5ce3)
}
func (e *LangPackLanguage) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Official) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.BaseLangCode) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.Rtl) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.Beta) {
		flag |= 1 << 3
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Official) {
	}
	if !zero.IsZeroVal(e.Rtl) {
	}
	if !zero.IsZeroVal(e.Beta) {
	}
	buf.PutString(e.Name)
	buf.PutString(e.NativeName)
	buf.PutString(e.LangCode)
	if !zero.IsZeroVal(e.BaseLangCode) {
		buf.PutString(e.BaseLangCode)
	}
	buf.PutString(e.PluralCode)
	buf.PutInt(e.StringsCount)
	buf.PutInt(e.TranslatedCount)
	buf.PutString(e.TranslationsUrl)
	return buf.Result()
}

func (e *LangPackLanguage) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<0 > 0 {
		e.Official = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.Rtl = buf.PopBool()
	}
	if flags&1<<3 > 0 {
		e.Beta = buf.PopBool()
	}
	e.Name = buf.PopString()
	e.NativeName = buf.PopString()
	e.LangCode = buf.PopString()
	if flags&1<<1 > 0 {
		e.BaseLangCode = buf.PopString()
	}
	e.PluralCode = buf.PopString()
	e.StringsCount = buf.PopInt()
	e.TranslatedCount = buf.PopInt()
	e.TranslationsUrl = buf.PopString()
}

type StatsBroadcastStats struct {
	Period                    *StatsDateRangeDays           `validate:"required"`
	Followers                 *StatsAbsValueAndPrev         `validate:"required"`
	ViewsPerPost              *StatsAbsValueAndPrev         `validate:"required"`
	SharesPerPost             *StatsAbsValueAndPrev         `validate:"required"`
	EnabledNotifications      *StatsPercentValue            `validate:"required"`
	GrowthGraph               StatsGraph                    `validate:"required"`
	FollowersGraph            StatsGraph                    `validate:"required"`
	MuteGraph                 StatsGraph                    `validate:"required"`
	TopHoursGraph             StatsGraph                    `validate:"required"`
	InteractionsGraph         StatsGraph                    `validate:"required"`
	IvInteractionsGraph       StatsGraph                    `validate:"required"`
	ViewsBySourceGraph        StatsGraph                    `validate:"required"`
	NewFollowersBySourceGraph StatsGraph                    `validate:"required"`
	LanguagesGraph            StatsGraph                    `validate:"required"`
	RecentMessageInteractions []*MessageInteractionCounters `validate:"required"`
}

func (e *StatsBroadcastStats) CRC() uint32 {
	return uint32(0xbdf78394)
}
func (e *StatsBroadcastStats) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Period.Encode())
	buf.PutRawBytes(e.Followers.Encode())
	buf.PutRawBytes(e.ViewsPerPost.Encode())
	buf.PutRawBytes(e.SharesPerPost.Encode())
	buf.PutRawBytes(e.EnabledNotifications.Encode())
	buf.PutRawBytes(e.GrowthGraph.Encode())
	buf.PutRawBytes(e.FollowersGraph.Encode())
	buf.PutRawBytes(e.MuteGraph.Encode())
	buf.PutRawBytes(e.TopHoursGraph.Encode())
	buf.PutRawBytes(e.InteractionsGraph.Encode())
	buf.PutRawBytes(e.IvInteractionsGraph.Encode())
	buf.PutRawBytes(e.ViewsBySourceGraph.Encode())
	buf.PutRawBytes(e.NewFollowersBySourceGraph.Encode())
	buf.PutRawBytes(e.LanguagesGraph.Encode())
	buf.PutVector(e.RecentMessageInteractions)
	return buf.Result()
}

func (e *StatsBroadcastStats) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Period = buf.PopObj().(*StatsDateRangeDays)
	e.Followers = buf.PopObj().(*StatsAbsValueAndPrev)
	e.ViewsPerPost = buf.PopObj().(*StatsAbsValueAndPrev)
	e.SharesPerPost = buf.PopObj().(*StatsAbsValueAndPrev)
	e.EnabledNotifications = buf.PopObj().(*StatsPercentValue)
	e.GrowthGraph = StatsGraph(buf.PopObj())
	e.FollowersGraph = StatsGraph(buf.PopObj())
	e.MuteGraph = StatsGraph(buf.PopObj())
	e.TopHoursGraph = StatsGraph(buf.PopObj())
	e.InteractionsGraph = StatsGraph(buf.PopObj())
	e.IvInteractionsGraph = StatsGraph(buf.PopObj())
	e.ViewsBySourceGraph = StatsGraph(buf.PopObj())
	e.NewFollowersBySourceGraph = StatsGraph(buf.PopObj())
	e.LanguagesGraph = StatsGraph(buf.PopObj())
	e.RecentMessageInteractions = buf.PopVector(reflect.TypeOf(*MessageInteractionCounters{})).([]*MessageInteractionCounters)
}

type TopPeer struct {
	Peer   Peer    `validate:"required"`
	Rating float64 `validate:"required"`
}

func (e *TopPeer) CRC() uint32 {
	return uint32(0xedcdc05b)
}
func (e *TopPeer) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Peer.Encode())
	buf.PutDouble(e.Rating)
	return buf.Result()
}

func (e *TopPeer) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Peer = Peer(buf.PopObj())
	e.Rating = buf.PopDouble()
}

type InputSingleMedia struct {
	__flagsPosition struct{}        // flags param position `validate:"required"`
	Media           InputMedia      `validate:"required"`
	RandomId        int64           `validate:"required"`
	Message         string          `validate:"required"`
	Entities        []MessageEntity `flag:"0"`
}

func (e *InputSingleMedia) CRC() uint32 {
	return uint32(0x1cc6e91f)
}
func (e *InputSingleMedia) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Entities) {
		flag |= 1 << 0
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutRawBytes(e.Media.Encode())
	buf.PutLong(e.RandomId)
	buf.PutString(e.Message)
	if !zero.IsZeroVal(e.Entities) {
		buf.PutVector(e.Entities)
	}
	return buf.Result()
}

func (e *InputSingleMedia) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	e.Media = InputMedia(buf.PopObj())
	e.RandomId = buf.PopLong()
	e.Message = buf.PopString()
	if flags&1<<0 > 0 {
		e.Entities = buf.PopVector(reflect.TypeOf(MessageEntity{})).([]MessageEntity)
	}
}

type AccountWebAuthorizations struct {
	Authorizations []*WebAuthorization `validate:"required"`
	Users          []User              `validate:"required"`
}

func (e *AccountWebAuthorizations) CRC() uint32 {
	return uint32(0xed56c9fc)
}
func (e *AccountWebAuthorizations) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutVector(e.Authorizations)
	buf.PutVector(e.Users)
	return buf.Result()
}

func (e *AccountWebAuthorizations) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Authorizations = buf.PopVector(reflect.TypeOf(*WebAuthorization{})).([]*WebAuthorization)
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
}

type JSONObjectValue struct {
	Key   string    `validate:"required"`
	Value JSONValue `validate:"required"`
}

func (e *JSONObjectValue) CRC() uint32 {
	return uint32(0xc0de1bd9)
}
func (e *JSONObjectValue) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Key)
	buf.PutRawBytes(e.Value.Encode())
	return buf.Result()
}

func (e *JSONObjectValue) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Key = buf.PopString()
	e.Value = JSONValue(buf.PopObj())
}

type MessagesSearchCounter struct {
	__flagsPosition struct{}       // flags param position `validate:"required"`
	Inexact         bool           `flag:"1,encoded_in_bitflags"`
	Filter          MessagesFilter `validate:"required"`
	Count           int32          `validate:"required"`
}

func (e *MessagesSearchCounter) CRC() uint32 {
	return uint32(0xe844ebff)
}
func (e *MessagesSearchCounter) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.Inexact) {
		flag |= 1 << 1
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Inexact) {
	}
	buf.PutRawBytes(e.Filter.Encode())
	buf.PutInt(e.Count)
	return buf.Result()
}

func (e *MessagesSearchCounter) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<1 > 0 {
		e.Inexact = buf.PopBool()
	}
	e.Filter = MessagesFilter(buf.PopObj())
	e.Count = buf.PopInt()
}

type RestrictionReason struct {
	Platform string `validate:"required"`
	Reason   string `validate:"required"`
	Text     string `validate:"required"`
}

func (e *RestrictionReason) CRC() uint32 {
	return uint32(0xd072acb4)
}
func (e *RestrictionReason) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutString(e.Platform)
	buf.PutString(e.Reason)
	buf.PutString(e.Text)
	return buf.Result()
}

func (e *RestrictionReason) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Platform = buf.PopString()
	e.Reason = buf.PopString()
	e.Text = buf.PopString()
}

type ReceivedNotifyMessage struct {
	Id    int32 `validate:"required"`
	Flags int32 `validate:"required"`
}

func (e *ReceivedNotifyMessage) CRC() uint32 {
	return uint32(0xa384b779)
}
func (e *ReceivedNotifyMessage) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.Id)
	buf.PutInt(e.Flags)
	return buf.Result()
}

func (e *ReceivedNotifyMessage) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Id = buf.PopInt()
	e.Flags = buf.PopInt()
}

type StickerSet struct {
	__flagsPosition struct{}  // flags param position `validate:"required"`
	Archived        bool      `flag:"1,encoded_in_bitflags"`
	Official        bool      `flag:"2,encoded_in_bitflags"`
	Masks           bool      `flag:"3,encoded_in_bitflags"`
	Animated        bool      `flag:"5,encoded_in_bitflags"`
	InstalledDate   int32     `flag:"0"`
	Id              int64     `validate:"required"`
	AccessHash      int64     `validate:"required"`
	Title           string    `validate:"required"`
	ShortName       string    `validate:"required"`
	Thumb           PhotoSize `flag:"4"`
	ThumbDcId       int32     `flag:"4"`
	Count           int32     `validate:"required"`
	Hash            int32     `validate:"required"`
}

func (e *StickerSet) CRC() uint32 {
	return uint32(0xeeb46f27)
}
func (e *StickerSet) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.InstalledDate) {
		flag |= 1 << 0
	}
	if !zero.IsZeroVal(e.Archived) {
		flag |= 1 << 1
	}
	if !zero.IsZeroVal(e.Official) {
		flag |= 1 << 2
	}
	if !zero.IsZeroVal(e.Masks) {
		flag |= 1 << 3
	}
	if !zero.IsZeroVal(e.Thumb) || !zero.IsZeroVal(e.ThumbDcId) {
		flag |= 1 << 4
	}
	if !zero.IsZeroVal(e.Animated) {
		flag |= 1 << 5
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	if !zero.IsZeroVal(e.Archived) {
	}
	if !zero.IsZeroVal(e.Official) {
	}
	if !zero.IsZeroVal(e.Masks) {
	}
	if !zero.IsZeroVal(e.Animated) {
	}
	if !zero.IsZeroVal(e.InstalledDate) {
		buf.PutInt(e.InstalledDate)
	}
	buf.PutLong(e.Id)
	buf.PutLong(e.AccessHash)
	buf.PutString(e.Title)
	buf.PutString(e.ShortName)
	if !zero.IsZeroVal(e.Thumb) {
		buf.PutRawBytes(e.Thumb.Encode())
	}
	if !zero.IsZeroVal(e.ThumbDcId) {
		buf.PutInt(e.ThumbDcId)
	}
	buf.PutInt(e.Count)
	buf.PutInt(e.Hash)
	return buf.Result()
}

func (e *StickerSet) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	if flags&1<<1 > 0 {
		e.Archived = buf.PopBool()
	}
	if flags&1<<2 > 0 {
		e.Official = buf.PopBool()
	}
	if flags&1<<3 > 0 {
		e.Masks = buf.PopBool()
	}
	if flags&1<<5 > 0 {
		e.Animated = buf.PopBool()
	}
	if flags&1<<0 > 0 {
		e.InstalledDate = buf.PopInt()
	}
	e.Id = buf.PopLong()
	e.AccessHash = buf.PopLong()
	e.Title = buf.PopString()
	e.ShortName = buf.PopString()
	if flags&1<<4 > 0 {
		e.Thumb = PhotoSize(buf.PopObj())
	}
	if flags&1<<4 > 0 {
		e.ThumbDcId = buf.PopInt()
	}
	e.Count = buf.PopInt()
	e.Hash = buf.PopInt()
}

type AccountAuthorizationForm struct {
	__flagsPosition  struct{}             // flags param position `validate:"required"`
	RequiredTypes    []SecureRequiredType `validate:"required"`
	Values           []*SecureValue       `validate:"required"`
	Errors           []SecureValueError   `validate:"required"`
	Users            []User               `validate:"required"`
	PrivacyPolicyUrl string               `flag:"0"`
}

func (e *AccountAuthorizationForm) CRC() uint32 {
	return uint32(0xad2e1cd8)
}
func (e *AccountAuthorizationForm) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	var flag uint32
	if !zero.IsZeroVal(e.PrivacyPolicyUrl) {
		flag |= 1 << 0
	}
	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutVector(e.RequiredTypes)
	buf.PutVector(e.Values)
	buf.PutVector(e.Errors)
	buf.PutVector(e.Users)
	if !zero.IsZeroVal(e.PrivacyPolicyUrl) {
		buf.PutString(e.PrivacyPolicyUrl)
	}
	return buf.Result()
}

func (e *AccountAuthorizationForm) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	flags := buf.PopUint()
	e.RequiredTypes = buf.PopVector(reflect.TypeOf(SecureRequiredType{})).([]SecureRequiredType)
	e.Values = buf.PopVector(reflect.TypeOf(*SecureValue{})).([]*SecureValue)
	e.Errors = buf.PopVector(reflect.TypeOf(SecureValueError{})).([]SecureValueError)
	e.Users = buf.PopVector(reflect.TypeOf(User{})).([]User)
	if flags&1<<0 > 0 {
		e.PrivacyPolicyUrl = buf.PopString()
	}
}

type UpdatesState struct {
	Pts         int32 `validate:"required"`
	Qts         int32 `validate:"required"`
	Date        int32 `validate:"required"`
	Seq         int32 `validate:"required"`
	UnreadCount int32 `validate:"required"`
}

func (e *UpdatesState) CRC() uint32 {
	return uint32(0xa56c2a3e)
}
func (e *UpdatesState) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.Pts)
	buf.PutInt(e.Qts)
	buf.PutInt(e.Date)
	buf.PutInt(e.Seq)
	buf.PutInt(e.UnreadCount)
	return buf.Result()
}

func (e *UpdatesState) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.Pts = buf.PopInt()
	e.Qts = buf.PopInt()
	e.Date = buf.PopInt()
	e.Seq = buf.PopInt()
	e.UnreadCount = buf.PopInt()
}

type BotInfo struct {
	UserId      int32         `validate:"required"`
	Description string        `validate:"required"`
	Commands    []*BotCommand `validate:"required"`
}

func (e *BotInfo) CRC() uint32 {
	return uint32(0x98e81d3a)
}
func (e *BotInfo) Encode() []byte {
	err := validator.New().Struct(e)
	dry.PanicIfErr(err)

	buf := mtproto.NewEncodeBuf(512)
	buf.PutUint(e.CRC())
	buf.PutInt(e.UserId)
	buf.PutString(e.Description)
	buf.PutVector(e.Commands)
	return buf.Result()
}

func (e *BotInfo) DecodeFrom(buf *mtproto.Decoder) {
	crc := buf.PopUint()
	if crc != e.CRC() {
		panic("wrong type: " + fmt.Sprintf("%#v", crc))
	}
	e.UserId = buf.PopInt()
	e.Description = buf.PopString()
	e.Commands = buf.PopVector(reflect.TypeOf(*BotCommand{})).([]*BotCommand)
}
