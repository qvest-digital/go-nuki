package command

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

// InvalidPinError will be returned if the given pin is invalid
var InvalidPinError = fmt.Errorf("the given pin is invalid")

type Pin uint16

func NewPin(pin string) (Pin, error) {
	if len(pin) != 4 {
		return 0, InvalidPinError
	}
	rawPin, err := strconv.ParseUint(pin, 10, 16)
	if err != nil {
		return 0, InvalidPinError
	}

	return Pin(rawPin), nil
}

func (p Pin) AsByte() []byte {
	pinAsByte := make([]byte, 2)
	binary.LittleEndian.PutUint16(pinAsByte, uint16(p))

	return pinAsByte
}
