package command

import (
	"encoding/binary"
	"github.com/howeyc/crc16"
	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/box"
	"github.com/kevinburke/nacl/randombytes"
)

/*
  +---------------------------------------------+-------------------------------------------------+
  | ADATA                                       | PDATA                                           |
  +---------------------------------------------+-------------------------------------------------+
  | nonce   | authorization id | message length | authorization id | command id | payload | CRC   |
  | 24 Byte |  4 Byte          | 2 Byte         | 4Byte            | 2 Byte     | n Byte  | 2Byte |
  +---------------------------------------------+-------------------------------------------------+
  | PLAIN (UNENCRYPTED)                         | ENCRYPTED                                       |
  +---------------------------------------------+-------------------------------------------------+
*/

func EncryptCommand(authId uint32, privateKey []byte, nukiPubKey []byte, plainCmd Command) Command {
	pdata := make([]byte, 0, len(plainCmd)+4)

	idAsByte := make([]byte, 4)
	binary.LittleEndian.PutUint32(idAsByte, authId)

	pdata = append(pdata, idAsByte...)
	pdata = append(pdata, plainCmd[:len(plainCmd)-2]...) //remove old crc

	crcPart := make([]byte, 2)
	crc := crc16.ChecksumCCITTFalse(pdata)
	binary.LittleEndian.PutUint16(crcPart, crc)

	pdata = append(pdata, crcPart...)

	sharedKey := box.Precompute(nacl.Key(nukiPubKey), nacl.Key(privateKey))
	nonce := newNonce192()

	encrypted := box.SealAfterPrecomputation(nonce[:], pdata, nacl.Nonce(nonce), sharedKey)
	encrypted = encrypted[24:] //remove nonce from encrypted

	message := make([]byte, 0, 24+4+2+len(encrypted))
	message = append(message, nonce...)
	message = append(message, idAsByte...)

	lengthAsByte := make([]byte, 2)
	binary.LittleEndian.PutUint16(lengthAsByte, uint16(len(encrypted)))
	message = append(message, lengthAsByte...)
	message = append(message, encrypted...)

	return message
}

func DecryptCommand(encryptedCmd Command, privateKey []byte, nukiPubKey []byte) (authId uint32, decrypted Command) {
	nonce := encryptedCmd[:24]
	//authId := encryptedCmd[24:28]
	length := binary.LittleEndian.Uint16(encryptedCmd[28:30])

	encrypted := encryptedCmd[30 : 30+length]

	sharedKey := box.Precompute(nacl.Key(nukiPubKey), nacl.Key(privateKey))
	decrypted, ok := box.OpenAfterPrecomputation(nil, encrypted, nacl.Nonce(nonce), sharedKey)
	if !ok {
		return 0, nil
	}

	if decrypted.CheckCRC() {
		//fix crc
		binary.LittleEndian.PutUint16(decrypted[len(decrypted)-2:], crc16.ChecksumCCITTFalse(decrypted[4:len(decrypted)-2]))
	}

	return binary.LittleEndian.Uint32(decrypted[:4]), decrypted[4:]
}

// only for monkey patching purposes
var newNonce192 = func() []byte {
	nonce := new([24]byte)
	randombytes.MustRead(nonce[:])
	return nonce[:]
}
