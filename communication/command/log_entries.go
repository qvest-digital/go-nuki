package command

import "encoding/binary"

type LogSortOrder uint8

const (
	LogSortOrderAscending  = LogSortOrder(0x00)
	LogSortOrderDescending = LogSortOrder(0x01)
)

func NewRequestLogEntriesCountCommand(pin Pin, nonce []byte) Command {
	return newRequestLogEntriesCommand(0, 0, 0, 0x01, pin, nonce)
}

func NewRequestLogEntriesCommand(startIndex uint32, count uint16, order LogSortOrder, pin Pin, nonce []byte) Command {
	return newRequestLogEntriesCommand(startIndex, count, order, 0x00, pin, nonce)
}

func newRequestLogEntriesCommand(startIndex uint32, count uint16, order LogSortOrder, totalCount uint8, pin Pin, nonce []byte) Command {
	payload := make([]byte, 6, 4+2+1+1+len(nonce)+2)
	binary.LittleEndian.PutUint32(payload[0:4], startIndex)
	binary.LittleEndian.PutUint16(payload[4:6], count)
	payload = append(payload, uint8(order))
	payload = append(payload, totalCount)
	payload = append(payload, nonce...)
	payload = append(payload, pin.AsByte()...)

	return NewCommand(IdRequestLogEntries, payload)
}
