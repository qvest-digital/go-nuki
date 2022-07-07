package command

import (
	"encoding/binary"
	"strings"
)

type LockAction uint8

const (
	LockActionUnlock               = LockAction(0x01)
	LockActionLock                 = LockAction(0x02)
	LockActionUnlatch              = LockAction(0x03)
	LockActionLockAndGo            = LockAction(0x04)
	LockActionLockAndGoWithUnlatch = LockAction(0x05)
	LockActionFullLock             = LockAction(0x06)
	LockActionFobAction1           = LockAction(0x81)
	LockActionFobAction2           = LockAction(0x82)
	LockActionFobAction3           = LockAction(0x83)

	LockActionFlagAutoUnlock = 0b0000_0001
	LockActionFlagForce      = 0b0000_0010
)

func NewLockAction(action LockAction, appId uint32, flags uint8, nameSuffix *string, nonce []byte) Command {
	payloadLen := 1 + 4 + 1 + 32
	if nameSuffix != nil {
		payloadLen += 20
	}

	payload := make([]byte, 0, payloadLen)
	payload = append(payload, uint8(action))

	idAsByte := make([]byte, 4)
	binary.LittleEndian.PutUint32(idAsByte, appId)

	payload = append(payload, idAsByte...)
	payload = append(payload, flags)

	if nameSuffix != nil {
		n := *nameSuffix
		if len(n) > 20 {
			n = n[:20]
		} else if len(n) < 20 {
			n += strings.Repeat("\x00", 20-len(n))
		}

		payload = append(payload, n...)
	}
	payload = append(payload, nonce...)

	return NewCommand(IdLockAction, payload)
}
