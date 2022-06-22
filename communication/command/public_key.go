package command

func NewPublicKey(publicKey []byte) Command {
	return NewCommand(IdPublicKey, publicKey)
}

type PublicKeyCommand Command

func (c Command) AsPublicKeyCommand() PublicKeyCommand {
	if !c.Is(IdPublicKey) {
		return nil
	}

	return PublicKeyCommand(c)
}

func (p PublicKeyCommand) PublicKey() []byte {
	return Command(p).Payload()
}
