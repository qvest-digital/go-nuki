package command

import "encoding/binary"

func NewRequest(dataType Id) Command {
	typeAsByte := make([]byte, 2)
	binary.LittleEndian.PutUint16(typeAsByte, uint16(dataType))

	return NewCommand(IdRequestData, typeAsByte)
}
