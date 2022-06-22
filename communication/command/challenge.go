package command

type ChallengeCommand Command

func (c Command) AsChallengeCommand() ChallengeCommand {
	if !c.Is(IdChallenge) {
		return nil
	}

	return ChallengeCommand(c)
}

func (c ChallengeCommand) Nonce() []byte {
	return Command(c).Payload()
}
