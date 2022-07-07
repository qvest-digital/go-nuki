package nuki

import (
	"context"
	"fmt"
	"github.com/tarent/go-nuki/communication"
	"github.com/tarent/go-nuki/communication/command"
)

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
