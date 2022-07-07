package nuki

import (
	"context"
	"fmt"
	"github.com/tarent/go-nuki/communication/command"
)

// Reboot will trigger a reboot of the connected device.
// After the reboot you have to re-establish the connection to the device via EstablishConnection!
func (c *Client) Reboot(ctx context.Context, pin string) error {
	parsedPin, err := c.checkPreconditionAndParsePin(pin)
	if err != nil {
		return err
	}

	if c.client == nil {
		return ConnectionNotEstablishedError
	}
	if c.udioCom == nil {
		return UnauthenticatedError
	}

	err = c.udioCom.Send(command.NewRequest(command.IdChallenge))
	if err != nil {
		return fmt.Errorf("unable to send request for challenge: %w", err)
	}

	challenge, err := c.udioCom.WaitForSpecificResponse(ctx, command.IdChallenge, c.responseTimeout)
	if err != nil {
		return fmt.Errorf("error while waiting for challenge: %w", err)
	}

	err = c.udioCom.Send(command.NewRequestReboot(parsedPin, challenge.AsChallengeCommand().Nonce()))
	if err != nil {
		return fmt.Errorf("unable to send action: %w", err)
	}

	c.Close()
	return nil
}
