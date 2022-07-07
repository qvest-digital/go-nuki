package nuki

import (
	"context"
	"fmt"
	"github.com/tarent/go-nuki/communication/command"
)

// ReadStates will request the current states for the connected and paired nuki device and return the result.
func (c *Client) ReadStates(ctx context.Context) (command.StatesCommand, error) {
	if c.client == nil {
		return nil, ConnectionNotEstablishedError
	}
	if c.udioCom == nil {
		return nil, UnauthenticatedError
	}

	err := c.udioCom.Send(command.NewRequest(command.IdStates))
	if err != nil {
		return nil, fmt.Errorf("unable to send request for device states: %w", err)
	}

	statesCommand, err := c.udioCom.WaitForSpecificResponse(ctx, command.IdStates, c.responseTimeout)
	if err != nil {
		return nil, fmt.Errorf("error while waiting for device states: %w", err)
	}

	return statesCommand.AsStatesCommand(), nil
}
