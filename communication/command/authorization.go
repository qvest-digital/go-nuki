package command

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/box"
	"github.com/kevinburke/nacl/randombytes"
	"strings"
)

type ClientId uint32
type ClientIdType uint8
type AuthorizationId uint32

const (
	ClientIdTypeApp    = ClientIdType(0)
	ClientIdTypeBridge = ClientIdType(1)
	ClientIdTypeFob    = ClientIdType(2)
	ClientIdTypeKeypad = ClientIdType(3)
)

func NewAuthorizationAuthenticator(
	nukiNonce []byte, nukiPubKey []byte,
	privateKey []byte, publicKey []byte) Command {

	sharedKey := box.Precompute(nacl.Key(nukiPubKey), nacl.Key(privateKey))

	valueR := make([]byte, 0, len(publicKey)+len(nukiPubKey)+len(nukiNonce))
	valueR = append(valueR, publicKey...)
	valueR = append(valueR, nukiPubKey...)
	valueR = append(valueR, nukiNonce...)

	hash := hmac.New(sha256.New, (*sharedKey)[:])
	hash.Write(valueR)

	return NewCommand(IdAuthorizationAuthenticator, hash.Sum(nil))
}

func NewAuthorizationData(
	nukiNonce []byte, nukiPubKey []byte,
	privateKey []byte,
	id ClientId, idType ClientIdType, name string) Command {

	sharedKey := box.Precompute(nacl.Key(nukiPubKey), nacl.Key(privateKey))

	valueR := make([]byte, 0, 1+4+32+32+32)
	valueR = append(valueR, uint8(idType))

	idAsByte := make([]byte, 4)
	binary.LittleEndian.PutUint32(idAsByte, uint32(id))
	valueR = append(valueR, idAsByte...)

	n := name
	if len(n) > 32 {
		n = name[:32]
	} else if len(n) < 32 {
		n += strings.Repeat("\x00", 32-len(n))
	}
	valueR = append(valueR, n...)

	nonce := newNonce256()

	valueR = append(valueR, nonce[:]...)
	valueR = append(valueR, nukiNonce...)

	hash := hmac.New(sha256.New, (*sharedKey)[:])
	hash.Write(valueR)

	authenticator := hash.Sum(nil)

	payload := make([]byte, 0, len(authenticator)+1+4+32+32)
	payload = append(payload, authenticator...)
	payload = append(payload, uint8(idType))
	payload = append(payload, idAsByte...)
	payload = append(payload, n...)
	payload = append(payload, nonce[:]...)

	return NewCommand(IdAuthorizationData, payload)
}

func NewAuthorizationIdConfirmation(
	nukiNonce []byte, nukiPubKey []byte,
	privateKey []byte,
	authId AuthorizationId) Command {

	sharedKey := box.Precompute(nacl.Key(nukiPubKey), nacl.Key(privateKey))

	valueR := make([]byte, 0, len(nukiNonce)+4)

	idAsByte := make([]byte, 4)
	binary.LittleEndian.PutUint32(idAsByte, uint32(authId))
	valueR = append(valueR, idAsByte...)
	valueR = append(valueR, nukiNonce...)

	hash := hmac.New(sha256.New, (*sharedKey)[:])
	hash.Write(valueR)

	authenticator := hash.Sum(nil)

	payload := make([]byte, 0, len(authenticator)+4)
	payload = append(payload, authenticator...)
	payload = append(payload, idAsByte...)

	return NewCommand(IdAuthorizationIDConfirmation, payload)
}

// only for monkey patching purposes
var newNonce256 = func() []byte {
	nonce := new([32]byte)
	randombytes.MustRead(nonce[:])
	return nonce[:]
}
