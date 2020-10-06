package telegram

// import (
// 	"reflect"
//
// 	"github.com/pkg/errors"
// 	"github.com/xelaj/mtproto/serialize"
// )
//
// type HelpGetConfigParams struct{}
//
// func (_ *HelpGetConfigParams) CRC() uint32 {
// 	return 0xc4f9186b
// }
//
// func (t *HelpGetConfigParams) Encode() []byte {
// 	buf := serialize.NewEncoder()
// 	buf.PutUint(t.CRC())
// 	return buf.Result()
// }


// type AuthSignInParams struct {
// 	PhoneNumber   string
// 	PhoneCodeHash string
// 	PhoneCode     string
// }
//
// func (_ *AuthSignInParams) CRC() uint32 {
// 	return 0xbcd51581
// }
//
// func (t *AuthSignInParams) Encode() []byte {
// 	buf := serialize.NewEncoder()
// 	buf.PutUint(t.CRC())
// 	buf.PutString(t.PhoneNumber)
// 	buf.PutString(t.PhoneCodeHash)
// 	buf.PutString(t.PhoneCode)
// 	return buf.Result()
// }
//
// func (c *Client) AuthSignIn(PhoneNumber, PhoneCodeHash, PhoneCode string) (AuthAuthorization, error) {
// 	data, err := c.MakeRequest(&AuthSignInParams{
// 		PhoneNumber:   PhoneNumber,
// 		PhoneCodeHash: PhoneCodeHash,
// 		PhoneCode:     PhoneCode,
// 	})
// 	if err != nil {
// 		return nil, errors.Wrap(err, "sending AuthSignIn")
// 	}
//
// 	resp, ok := data.(AuthAuthorization)
// 	if !ok {
// 		panic(errors.New("got invalid response type: " + reflect.TypeOf(data).String()))
// 	}
//
// 	return resp, nil
// }
