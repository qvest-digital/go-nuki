package nuki

import (
	"context"
	"fmt"
	"github.com/tarent/go-nuki/communication"
	"github.com/tarent/go-nuki/communication/command"
)

// ReadOpenerState will request the current opener state for the connected and paired nuki device and return the result.
func (c *Client) ReadOpenerState(ctx context.Context) (command.OpenerStatesCommand, error) {
	if c.client == nil {
		return nil, ConnectionNotEstablishedError
	}
	if c.udioCom == nil {
		return nil, UnauthenticatedError
	}
	if c.udioCom.GetDeviceType() != communication.DeviceTypeOpener {
		return nil, fmt.Errorf("unexpected device type: this operation is only available for opener")
	}

	err := c.udioCom.Send(command.NewRequest(command.IdOpenerStates))
	if err != nil {
		return nil, fmt.Errorf("unable to send request for opener states: %w", err)
	}

	statesCommand, err := c.udioCom.WaitForSpecificResponse(ctx, command.IdOpenerStates, c.responseTimeout)
	if err != nil {
		return nil, fmt.Errorf("error while waiting for opener states: %w", err)
	}

	return statesCommand.AsOpenerStatesCommand(), nil
}

// PerformOpen will trigger the electric strike actuation to open the door and return the result.
func (c *Client) PerformOpen(ctx context.Context, appId command.ClientId) error {
	return c.PerformOpenAction(ctx, appId, command.OpenActionElectricStrikeActuation)
}

// PerformOpenAction will request the connected and paired nuki opener to perform the given open action.
func (c *Client) PerformOpenAction(ctx context.Context, appId command.ClientId, action command.OpenAction) error {
	if c.udioCom.GetDeviceType() != communication.DeviceTypeOpener {
		return fmt.Errorf("unexpected device type: this operation is only available for opener")
	}

	return c.PerformAction(ctx, func(nonce []byte) command.Command {
		return command.NewOpenAction(action, uint32(appId), 0, nil, nonce)
	})
}
