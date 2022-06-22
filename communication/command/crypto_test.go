package command

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestEncryptCommand(t *testing.T) {
	nukiPubKey, _ := hex.DecodeString("2FE57DA347CD62431528DAAC5FBB290730FFF684AFC4CFC2ED90995F58CB3B74")
	privKey, _ := hex.DecodeString("8CAA54672307BFFDF5EA183FC607158D2011D008ECA6A1088614FF0853A5AA07")

	origNewNonce := newNonce192
	defer func() {
		newNonce192 = origNewNonce
	}()
	newNonce192 = func() []byte {
		v, _ := hex.DecodeString("37917F1AF31EC5940705F34D1E5550607D5B2F9FE7D496B6")
		return v
	}

	result := EncryptCommand(2, privKey, nukiPubKey, NewRequest(IdKeyturnerStates))

	assert.Equal(t, "37917F1AF31EC5940705F34D1E5550607D5B2F9FE7D496B6020000001A00670D124926004366532E8D927A33FE84E782A9594D39157D065E", strings.ToUpper(hex.EncodeToString(result)))
}

func TestDecryptCommand(t *testing.T) {
	nukiPubKey, _ := hex.DecodeString("2FE57DA347CD62431528DAAC5FBB290730FFF684AFC4CFC2ED90995F58CB3B74")
	privKey, _ := hex.DecodeString("8CAA54672307BFFDF5EA183FC607158D2011D008ECA6A1088614FF0853A5AA07")
	encryptedCmd, _ := hex.DecodeString("90B0757CFED0243017EAF5E089F8583B9839D61B050924D2020000002700B13938B67121B6D528E7DE206B0D7C5A94587A471B33EBFB012CED8F1261135566ED756E3910B5")

	authId, result := DecryptCommand(encryptedCmd, privKey, nukiPubKey)

	assert.Equal(t, uint32(2), authId)
	assert.Equal(t, "020100E0070307080F1E3C0000200A", strings.ToUpper(hex.EncodeToString(result.Payload())))
	assert.True(t, result.CheckCRC())
}
