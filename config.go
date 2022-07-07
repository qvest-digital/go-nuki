package nuki

import (
	"context"
	"fmt"
	"github.com/tarent/go-nuki/communication/command"
	"time"
)

// UpdateTime set the given time on the connected device.
func (c *Client) UpdateTime(ctx context.Context, pin string, t time.Time) error {
	parsedPin, err := c.checkPreconditionAndParsePin(pin)
	if err != nil {
		return err
	}

	return c.PerformAction(ctx, func(nonce []byte) command.Command {
		return command.NewUpdateTime(t, parsedPin, nonce)
	})
}

// ReadConfig will request and return the applied config of the connected device.
func (c *Client) ReadConfig(ctx context.Context) (command.ConfigCommand, error) {
	if c.client == nil {
		return nil, ConnectionNotEstablishedError
	}
	if c.udioCom == nil {
		return nil, UnauthenticatedError
	}

	err := c.udioCom.Send(command.NewRequest(command.IdChallenge))
	if err != nil {
		return nil, fmt.Errorf("unable to send request for challenge: %w", err)
	}

	challenge, err := c.udioCom.WaitForSpecificResponse(ctx, command.IdChallenge, c.responseTimeout)
	if err != nil {
		return nil, fmt.Errorf("error while waiting for challenge: %w", err)
	}
	err = c.udioCom.Send(command.NewRequestConfig(challenge.AsChallengeCommand().Nonce()))
	if err != nil {
		return nil, fmt.Errorf("unable to send request for config: %w", err)
	}

	config, err := c.udioCom.WaitForSpecificResponse(ctx, command.IdConfig, c.responseTimeout)
	if err != nil {
		return nil, fmt.Errorf("error while waiting for config: %w", err)
	}

	return config.AsConfigCommand(), nil
}
