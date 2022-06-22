package command

import "encoding/binary"

type AuthorizationIdCommand Command

func (c Command) AsAuthorizationIdCommand() AuthorizationIdCommand {
	if !c.Is(IdAuthorizationID) {
		return nil
	}

	return AuthorizationIdCommand(c)
}

func (a AuthorizationIdCommand) Authenticator() []byte {
	return Command(a).Payload()[:32]
}

func (a AuthorizationIdCommand) AuthorizationId() AuthorizationId {
	return AuthorizationId(binary.LittleEndian.Uint32(Command(a).Payload()[32:36]))
}

func (a AuthorizationIdCommand) UUID() []byte {
	return Command(a).Payload()[36:52]
}

func (a AuthorizationIdCommand) Nonce() []byte {
	return Command(a).Payload()[52:]
}
