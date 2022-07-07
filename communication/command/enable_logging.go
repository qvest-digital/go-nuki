package command

func NewEnableLogging(enable bool, pin Pin, nonce []byte) Command {
	payload := make([]byte, 0, 1+len(nonce)+2)

	payload = append(payload, 0x00)
	if enable {
		payload[0] = 0x01
	}
	payload = append(payload, nonce...)
	payload = append(payload, pin.AsByte()...)

	return NewCommand(IdEnableLogging, payload)
}
