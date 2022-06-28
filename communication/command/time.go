package command

import (
	"encoding/binary"
	"time"
)

func NewUpdateTime(t time.Time, pin Pin, nonce []byte) Command {
	payload := make([]byte, 0, 7+len(nonce)+2)

	yearAsByte := make([]byte, 2)
	binary.LittleEndian.PutUint16(yearAsByte, uint16(t.Year()))
	payload = append(payload, yearAsByte...)
	payload = append(payload, uint8(t.Month()))
	payload = append(payload, uint8(t.Day()))
	payload = append(payload, uint8(t.Hour()))
	payload = append(payload, uint8(t.Minute()))
	payload = append(payload, uint8(t.Second()))
	payload = append(payload, nonce...)
	payload = append(payload, pin.AsByte()...)

	return NewCommand(IdUpdateTime, payload)
}
