package nuki

import (
	"context"
	"fmt"
	"github.com/tarent/go-nuki/communication"
	"github.com/tarent/go-nuki/communication/command"
)

// ReadLockerState will request the current locker state for the connected and paired nuki device and return the result.
func (c *Client) ReadLockerState(ctx context.Context) (command.KeyturnerStatesCommand, error) {
	if c.client == nil {
		return nil, ConnectionNotEstablishedError
	}
	if c.udioCom == nil {
		return nil, UnauthenticatedError
	}
	if c.udioCom.GetDeviceType() != communication.DeviceTypeSmartLock {
		return nil, fmt.Errorf("unexpected device type: this operation is only available for smart lock")
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

// PerformLock will request the connected and paired nuki smart lock to lock.
func (c *Client) PerformLock(ctx context.Context, appId command.ClientId) error {
	return c.PerformLockAction(ctx, appId, command.LockActionLock)
}

// PerformUnlock will request the connected and paired nuki smart lock to unlock.
func (c *Client) PerformUnlock(ctx context.Context, appId command.ClientId) error {
	return c.PerformLockAction(ctx, appId, command.LockActionUnlock)
}

// PerformLockAction will request the connected and paired nuki smart lock to perform the given lock action.
func (c *Client) PerformLockAction(ctx context.Context, appId command.ClientId, action command.LockAction) error {
	if c.udioCom.GetDeviceType() != communication.DeviceTypeSmartLock {
		return fmt.Errorf("unexpected device type: this operation is only available for smart lock")
	}

	return c.PerformAction(ctx, func(nonce []byte) command.Command {
		return command.NewLockAction(action, uint32(appId), 0, nil, nonce)
	})
}
