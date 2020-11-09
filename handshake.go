package mtproto

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/xelaj/go-dry"

	ige "github.com/xelaj/mtproto/aes_ige"
	"github.com/xelaj/mtproto/encoding/tl"
	"github.com/xelaj/mtproto/keys"
	"github.com/xelaj/mtproto/serialize"
)

// https://tlgrm.ru/docs/mtproto/auth_key
// https://core.telegram.org/mtproto/auth_key
func (m *MTProto) makeAuthKey() error {
	m.serviceModeActivated = true
	nonceFirst := serialize.RandomInt128()

	pqParams := new(serialize.ResPQ)
	if err := m.MakeRequest(&ReqPQParams{nonceFirst}, pqParams); err != nil {
		return errors.Wrap(err, "requesting first pq")
	}
	if nonceFirst.Cmp(pqParams.Nonce.Int) != 0 {
		return errors.New("handshake: Wrong nonce")
	}

	found := false
	for _, b := range pqParams.Fingerprints {
		fgpr, err := keys.RSAFingerprint(m.publicKey)
		if err != nil {
			return err
		}

		if uint64(b) == binary.LittleEndian.Uint64(fgpr) {
			found = true
			break
		}
	}
	if !found {
		return errors.New("handshake: Can't find fingerprint")
	}

	// (encoding) p_q_inner_data
	pq := big.NewInt(0).SetBytes(pqParams.Pq)
	p, q := splitPQ(pq)
	nonceSecond := serialize.RandomInt256()
	nonceServer := pqParams.ServerNonce

	message, err := tl.Encode(&serialize.PQInnerData{
		Pq:          pqParams.Pq,
		P:           p.Bytes(),
		Q:           q.Bytes(),
		Nonce:       nonceFirst,
		ServerNonce: nonceServer,
		NewNonce:    nonceSecond,
	})
	if err != nil {
		return err
	}

	hashAndMsg := make([]byte, 255)
	copy(hashAndMsg, append(dry.Sha1(string(message)), message...))

	encryptedMessage := doRSAencrypt(hashAndMsg, m.publicKey)

	fgpr, err := keys.RSAFingerprint(m.publicKey)
	if err != nil {
		return err
	}

	keyFingerprint := int64(binary.LittleEndian.Uint64(fgpr))

	var dhResponse serialize.ServerDHParams
	if err := m.MakeRequest(&ReqDHParamsParams{
		Nonce:                nonceFirst,
		ServerNonce:          nonceServer,
		P:                    p.Bytes(),
		Q:                    q.Bytes(),
		PublicKeyFingerprint: keyFingerprint,
		EncryptedData:        encryptedMessage,
	}, &dhResponse); err != nil {
		return errors.Wrap(err, "sending ReqDHParams")
	}

	dhParams, ok := dhResponse.(*serialize.ServerDHParamsOk)
	if !ok {
		return fmt.Errorf("need *serialize.ServerDHParamsOk, got: %T", dhResponse)
	}

	if nonceFirst.Cmp(dhParams.Nonce.Int) != 0 {
		return errors.New("handshake: Wrong nonce")
	}

	if nonceServer.Cmp(dhParams.ServerNonce.Int) != 0 {
		return errors.New("handshake: Wrong server_nonce")
	}

	// проверку по хешу, удаление рандомных байт происходит в этой функции
	decodedMessage := ige.DecryptMessageWithTempKeys(dhParams.EncryptedAnswer, nonceSecond.Int, nonceServer.Int)
	dhi := new(serialize.ServerDHInnerData)
	if err := tl.Decode(decodedMessage, dhi); err != nil {
		return err
	}

	if nonceFirst.Cmp(dhi.Nonce.Int) != 0 {
		return errors.New("Handshake: Wrong nonce")
	}
	if nonceServer.Cmp(dhi.ServerNonce.Int) != 0 {
		return errors.New("Handshake: Wrong server_nonce")
	}

	// вот это видимо как раз и есть часть диффи хеллмана, поэтому просто оставим как есть надеюсь сработает
	_, g_b, g_ab := makeGAB(dhi.G, big.NewInt(0).SetBytes(dhi.GA), big.NewInt(0).SetBytes(dhi.DhPrime))

	authKey := g_ab.Bytes()
	if authKey[0] == 0 {
		authKey = authKey[1:]
	}

	m.SetAuthKey(authKey)

	// что это я пока не знаю, видимо какой то очень специфичный способ сгенерить ключи
	t4 := make([]byte, 32+1+8)
	copy(t4[0:], nonceSecond.Bytes())
	t4[32] = 1
	copy(t4[33:], dry.Sha1Byte(m.GetAuthKey())[0:8])
	nonceHash1 := dry.Sha1Byte(t4)[4:20]
	salt := make([]byte, tl.LongLen)
	copy(salt, nonceSecond.Bytes()[:8])
	xor(salt, nonceServer.Bytes()[:8])
	m.serverSalt = int64(binary.LittleEndian.Uint64(salt))

	// (encoding) client_DH_inner_data
	clientDHDataMsg, err := tl.Encode(&serialize.ClientDHInnerData{
		Nonce:       nonceFirst,
		ServerNonce: nonceServer,
		Retry:       0,
		GB:          g_b.Bytes(),
	})
	if err != nil {
		return err
	}

	encryptedMessage = ige.EncryptMessageWithTempKeys(clientDHDataMsg, nonceSecond.Int, nonceServer.Int)

	var dhGenStatus serialize.SetClientDHParamsAnswer
	if err := m.MakeRequest(&SetClientDHParamsParams{
		Nonce:         nonceFirst,
		ServerNonce:   nonceServer,
		EncryptedData: encryptedMessage,
	}, &dhGenStatus); err != nil {
		return errors.Wrap(err, "sending clientDHParams")
	}

	dhg, ok := dhGenStatus.(*serialize.DHGenOk)
	if !ok {
		return errors.New("Handshake: Need DHGenOk")
	}
	if nonceFirst.Cmp(dhg.Nonce.Int) != 0 {
		return fmt.Errorf("Handshake: Wrong nonce: %v, %v", nonceFirst, dhg.Nonce)
	}
	if nonceServer.Cmp(dhg.ServerNonce.Int) != 0 {
		return fmt.Errorf("Handshake: Wrong server_nonce: %v, %v", nonceServer, dhg.ServerNonce)
	}
	if !bytes.Equal(nonceHash1, dhg.NewNonceHash1.Bytes()) {
		return fmt.Errorf(
			"handshake: Wrong new_nonce_hash1: %v, %v",
			hex.EncodeToString(nonceHash1),
			hex.EncodeToString(dhg.NewNonceHash1.Bytes()),
		)
	}
	m.serviceModeActivated = false

	// (all ok)
	err = m.SaveSession()
	return errors.Wrap(err, "saving session")
}
