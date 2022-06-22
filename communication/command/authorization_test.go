package command

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewAuthorizationAuthenticator(t *testing.T) {
	nukiNonce, _ := hex.DecodeString("6CD4163D159050C798553EAA57E278A579AFFCBC56F09FC57FE879E51C42DF17")
	nukiPubKey, _ := hex.DecodeString("2FE57DA347CD62431528DAAC5FBB290730FFF684AFC4CFC2ED90995F58CB3B74")
	privKey, _ := hex.DecodeString("8CAA54672307BFFDF5EA183FC607158D2011D008ECA6A1088614FF0853A5AA07")
	pubKey, _ := hex.DecodeString("F88127CCF48023B5CBE9101D24BAA8A368DA94E8C2E3CDE2DED29CE96AB50C15")

	result := NewAuthorizationAuthenticator(nukiNonce, nukiPubKey, privKey, pubKey)

	assert.Equal(t, "0500B09A0D3979A029E5FD027B519EAA200BC14AD3E163D3BE4563843E021073BCB1C357", strings.ToUpper(hex.EncodeToString(result)))
}

func TestNewAuthorizationData(t *testing.T) {
	nukiNonce, _ := hex.DecodeString("E0742CFEA39CB46109385BF91286A3C02F40EE86B0B62FC34033094DE41E2C0D")
	nukiPubKey, _ := hex.DecodeString("2FE57DA347CD62431528DAAC5FBB290730FFF684AFC4CFC2ED90995F58CB3B74")
	privKey, _ := hex.DecodeString("8CAA54672307BFFDF5EA183FC607158D2011D008ECA6A1088614FF0853A5AA07")

	origNewNonce := newNonce256
	defer func() {
		newNonce256 = origNewNonce
	}()
	newNonce256 = func() []byte {
		nonce, _ := hex.DecodeString("52AFE0A664B4E9B56DC6BD4CB718A6C9FED6BE17A7411072AA0D315378140577")
		return nonce
	}
	name := "Marc (Test)"

	id := ClientId(0)
	idType := ClientIdTypeApp

	result := NewAuthorizationData(nukiNonce, nukiPubKey, privKey, id, idType, name)

	assert.Equal(t, "0600CF1B9E7801E3196E6594E76D57908EE500AAD5C33F4B6E0BBEA0DDEF82967BFC00000000004D6172632028546573742900000000000000000000000000000000000000000052AFE0A664B4E9B56DC6BD4CB718A6C9FED6BE17A7411072AA0D31537814057769F2", strings.ToUpper(hex.EncodeToString(result)))
}

func TestNewAuthorizationIdConfirmation(t *testing.T) {
	nukiNonce, _ := hex.DecodeString("EA479915982F13C61D997A56678AD77791BFA7E95229A3DD34F87132BF3E3C97")
	nukiPubKey, _ := hex.DecodeString("2FE57DA347CD62431528DAAC5FBB290730FFF684AFC4CFC2ED90995F58CB3B74")
	privKey, _ := hex.DecodeString("8CAA54672307BFFDF5EA183FC607158D2011D008ECA6A1088614FF0853A5AA07")

	result := NewAuthorizationIdConfirmation(nukiNonce, nukiPubKey, privKey, 2)

	assert.Equal(t, "1E003A41B91A66FBC4D22EFEFBB7272140829695A3917433D5BEB981B76166D13F8A02000000CDF5", strings.ToUpper(hex.EncodeToString(result)))
}
