package command

import "encoding/binary"

func NewEnableLogging(enable bool, pin uint16, nonce []byte) Command {
	payload := make([]byte, 0, 1+len(nonce)+2)

	payload = append(payload, 0x00)
	if enable {
		payload[0] = 0x01
	}
	payload = append(payload, nonce...)

	pinAsByte := make([]byte, 2)
	binary.LittleEndian.PutUint16(pinAsByte, pin)
	payload = append(payload, pinAsByte...)

	return NewCommand(IdEnableLogging, payload)
}
