package command

import (
	"encoding/binary"
	"fmt"
	"github.com/howeyc/crc16"
)

const (
	IdRequestData                 = Id(0x0001)
	IdPublicKey                   = Id(0x0003)
	IdChallenge                   = Id(0x0004)
	IdAuthorizationAuthenticator  = Id(0x0005)
	IdAuthorizationData           = Id(0x0006)
	IdAuthorizationID             = Id(0x0007)
	IdRemoveUserAuthorization     = Id(0x0008)
	IdRequestAuthorizationEntries = Id(0x0009)
	IdAuthorizationEntry          = Id(0x000A)
	IdAuthorizationDataInvite     = Id(0x000B)
	IdKeyturnerStates             = Id(0x000C)
	IdOpenerStates                = Id(0x000C)
	IdLockAction                  = Id(0x000D)
	IdStatus                      = Id(0x000E)
	IdMostRecentCommand           = Id(0x000F)
	IdOpeningsClosingsSummary     = Id(0x0010)
	IdBatteryReport               = Id(0x0011)
	IdErrorReport                 = Id(0x0012)
	IdSetConfig                   = Id(0x0013)
	IdRequestConfig               = Id(0x0014)
	IdConfig                      = Id(0x0015)
	IdSetSecurityPIN              = Id(0x0019)
	IdRequestCalibration          = Id(0x001A)
	IdSetCalibrated               = Id(0x001A)
	IdRequestReboot               = Id(0x001D)
	IdAuthorizationIDConfirmation = Id(0x001E)
	IdAuthorizationIDInvite       = Id(0x001F)
	IdVerifySecurityPIN           = Id(0x0020)
	IdUpdateTime                  = Id(0x0021)
	IdUpdateUserAuthorization     = Id(0x0025)
	IdAuthorizationEntryCount     = Id(0x0027)
	IdStartBusSignalRecording     = Id(0x002F)
	IdRequestLogEntries           = Id(0x0031)
	IdLogEntry                    = Id(0x0032)
	IdLogEntryCount               = Id(0x0033)
	IdEnableLogging               = Id(0x0034)
	IdSetAdvancedConfig           = Id(0x0035)
	IdRequestAdvancedConfig       = Id(0x0036)
	IdAdvancedConfig              = Id(0x0037)
	IdAddTimeControlEntry         = Id(0x0039)
	IdTimeControlEntryID          = Id(0x003A)
	IdRemoveTimeControlEntry      = Id(0x003B)
	IdRequestTimeControlEntries   = Id(0x003C)
	IdTimeControlEntryCount       = Id(0x003D)
	IdTimeControlEntry            = Id(0x003E)
	IdUpdateTimeControlEntry      = Id(0x003F)
	IdAddKeypadCode               = Id(0x0041)
	IdKeypadCodeID                = Id(0x0042)
	IdRequestKeypadCodes          = Id(0x0043)
	IdKeypadCodeCount             = Id(0x0044)
	IdKeypadCode                  = Id(0x0045)
	IdUpdateKeypadCode            = Id(0x0046)
	IdRemoveKeypadCode            = Id(0x0047)
	IdKeypadAction                = Id(0x0048)
	IdContinuousModeAction        = Id(0x0057)
	IdSimpleLockAction            = Id(0x0100)
)

type Id uint16

type Command []byte

func (c Command) Id() Id {
	if len(c) >= 2 {
		return Id(binary.LittleEndian.Uint16(c[0:2]))
	}

	return 0
}

func (c Command) Payload() []byte {
	if len(c) < 4 {
		return nil
	}

	//first 2 byte are id and last 2 byte are crc
	return c[2 : len(c)-2]
}

func (c Command) CrcSum() uint16 {
	if len(c) >= 4 {
		return binary.LittleEndian.Uint16(c[len(c)-2:])
	}

	return 0
}

func (c Command) CheckCRC() bool {
	if len(c) == 0 {
		return true
	}

	return crc16.ChecksumCCITTFalse(c[:len(c)-2]) == c.CrcSum()
}

func (c Command) Is(t Id) bool {
	return c.Id() == t
}

func (c Command) String() string {
	var id Id
	var payload []byte
	var crc []byte

	if len(c) >= 2 {
		id = Id(binary.LittleEndian.Uint16(c[0:2]))
	}
	if len(c) >= 4 {
		crc = c[len(c)-2:]
		payload = c[2 : len(c)-2]
	}

	return fmt.Sprintf("Raw: %x; ID: %04x; Payload: %x; CRC: %x", []byte(c), id, payload, crc)
}

func NewCommand(id Id, payload []byte) Command {
	/*
	  +--------------------+----------+---------+
	  | command identifier | payload  |  CRC    |
	  | 2 Byte             |  n Byte  |  2 Byte |
	  +--------------------+----------+---------+
	*/

	cmd := make([]byte, 0, 4+len(payload))

	idPart := make([]byte, 2)
	binary.LittleEndian.PutUint16(idPart, uint16(id))

	cmd = append(cmd, idPart...)
	cmd = append(cmd, payload...)

	crcPart := make([]byte, 2)
	crc := crc16.ChecksumCCITTFalse(cmd)
	binary.LittleEndian.PutUint16(crcPart, crc)

	cmd = append(cmd, crcPart...)

	return cmd
}
