package command

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// InvalidPinError will be returned if the given pin is invalid
var InvalidPinError = fmt.Errorf("the given pin is invalid")

type Pin uint16

func NewPin(pin string) (Pin, error) {
	if len(pin) != 4 || hex.DecodedLen(len(pin)) != 2 {
		return 0, InvalidPinError
	}
	rawPin, err := hex.DecodeString(pin)
	if err != nil {
		return 0, InvalidPinError
	}

	return Pin(binary.BigEndian.Uint16(rawPin)), nil
}

func (p Pin) AsByte() []byte {
	pinAsByte := make([]byte, 2)
	binary.LittleEndian.PutUint16(pinAsByte, uint16(p))

	return pinAsByte
}
