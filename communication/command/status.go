package command

type CompletionStatus uint8

const (
	CompletionStatusComplete = CompletionStatus(0x00)
	CompletionStatusAccepted = CompletionStatus(0x01)
)

type StatusCommand Command

func (c Command) AsStatusCommand() StatusCommand {
	if !c.Is(IdStatus) {
		return nil
	}

	return StatusCommand(c)
}

func (s StatusCommand) Status() CompletionStatus {
	return CompletionStatus(Command(s).Payload()[0])
}

func (s StatusCommand) IsComplete() bool {
	return s.Status() == CompletionStatusComplete
}

func (s StatusCommand) IsAccepted() bool {
	return s.Status() == CompletionStatusAccepted
}
