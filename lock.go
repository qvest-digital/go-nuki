package nuki

import (
	"context"
	"fmt"
	"github.com/tarent/go-nuki/communication/command"
)

// ReadLockState will request the current lock state for the connected and paired nuki device and return the result.
func (c *Client) ReadLockState(ctx context.Context) (command.KeyturnerStatesCommand, error) {
	if c.client == nil {
		return nil, ConnectionNotEstablishedError
	}
	if c.udioCom == nil {
		return nil, UnauthenticatedError
	}

	err := c.udioCom.Send(command.NewRequest(command.IdKeyturnerStates))
	if err != nil {
		return nil, fmt.Errorf("unable to send request for keyturner states: %w", err)
	}

	statesCommand, err := c.udioCom.WaitForSpecificResponse(ctx, command.IdKeyturnerStates, c.responseTimeout)
	if err != nil {
		return nil, fmt.Errorf("error while waiting for keyturner states: %w", err)
	}

	return statesCommand.AsKeyturnerStatesCommand(), nil
}

// PerformLock will request the connected and paired nuki device to lock.
func (c *Client) PerformLock(ctx context.Context, appId command.ClientId) error {
	return c.performLockAction(ctx, appId, command.LockActionLock, command.LockStateLocking, command.LockStateLocked)
}

// PerformUnlock will request the connected and paired nuki device to unlock.
func (c *Client) PerformUnlock(ctx context.Context, appId command.ClientId) error {
	return c.performLockAction(ctx, appId, command.LockActionUnlock, command.LockStateUnlocking, command.LockStateUnlocked)
}

func (c *Client) performLockAction(ctx context.Context, appId command.ClientId, action command.LockAction, firstState, secondState command.LockState) error {
	if c.client == nil {
		return ConnectionNotEstablishedError
	}
	if c.udioCom == nil {
		return UnauthenticatedError
	}

	err := c.udioCom.Send(command.NewRequest(command.IdChallenge))
	if err != nil {
		return fmt.Errorf("unable to send request for challenge: %w", err)
	}

	challenge, err := c.udioCom.WaitForSpecificResponse(ctx, command.IdChallenge, c.responseTimeout)
	if err != nil {
		return fmt.Errorf("error while waiting for challenge: %w", err)
	}

	err = c.udioCom.Send(command.NewLockAction(
		action,
		uint32(appId),
		0,
		nil,
		challenge.AsChallengeCommand().Nonce(),
	))
	if err != nil {
		return fmt.Errorf("unable to send lock action: %w", err)
	}

	status1, err := c.udioCom.WaitForSpecificResponse(ctx, command.IdStatus, c.responseTimeout)
	if err != nil {
		return fmt.Errorf("error while waiting for first status: %w", err)
	}

	if !status1.AsStatusCommand().IsAccepted() {
		return fmt.Errorf("unexpected status: expect 0x%02x got 0x%02x", command.CompletionStatusAccepted, status1.AsStatusCommand().Status())
	}

	ktStates1, err := c.udioCom.WaitForSpecificResponse(ctx, command.IdKeyturnerStates, c.responseTimeout)
	if err != nil {
		return fmt.Errorf("error while waiting for first keyturner states: %w", err)
	}

	if ktStates1.AsKeyturnerStatesCommand().LockState() != firstState {
		return fmt.Errorf("unexpected lock state: expect 0x%02x got 0x%02x", firstState, ktStates1.AsKeyturnerStatesCommand().LockState())
	}

	ktStates2, err := c.udioCom.WaitForSpecificResponse(ctx, command.IdKeyturnerStates, c.responseTimeout)
	if err != nil {
		return fmt.Errorf("error while waiting for second keyturner states: %w", err)
	}

	if ktStates2.AsKeyturnerStatesCommand().LockState() != secondState {
		return fmt.Errorf("unexpected lock state: expect 0x%02x got 0x%02x", secondState, ktStates2.AsKeyturnerStatesCommand().LockState())
	}

	status2, err := c.udioCom.WaitForSpecificResponse(ctx, command.IdStatus, c.responseTimeout)
	if err != nil {
		return fmt.Errorf("error while waiting for second status: %w", err)
	}

	if !status2.AsStatusCommand().IsComplete() {
		return fmt.Errorf("unexpected status: expect 0x%02x got 0x%02x", command.CompletionStatusComplete, status2.AsStatusCommand().Status())
	}

	return nil
}
