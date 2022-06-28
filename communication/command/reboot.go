package command

func NewRequestReboot(pin Pin, nonce []byte) Command {
	payload := make([]byte, 0, len(nonce)+2)
	payload = append(payload, nonce...)
	payload = append(payload, pin.AsByte()...)

	return NewCommand(IdRequestReboot, payload)
}
