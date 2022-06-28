package communication

import (
	"context"
	"fmt"
	"github.com/tarent/go-nuki/communication/command"
	"github.com/tarent/go-nuki/logger"
	"time"
)

// mtu is the maximum transfer unit (how many bytes will be transfer at once)
const mtu = 20

type DeviceType uint8

const (
	DeviceTypeUnknown   = DeviceType(0x00)
	DeviceTypeSmartLock = DeviceType(0x01)
	DeviceTypeOpener    = DeviceType(0x02)
)

// TimeoutErr will be occurred if the used timeout exceeded
var TimeoutErr = fmt.Errorf("timeout exceeded")

type Communicator interface {
	// Send will send the given command to the connected nuki device
	Send(cmd command.Command) error

	// WaitForResponse will wait for response. It will return an error if the timeout exceeded, the context is closed
	// or the response was erroneous.
	WaitForResponse(ctx context.Context, timeout time.Duration) (command.Command, error)

	// WaitForSpecificResponse will wait for response. It will return an error if the timeout exceeded, the context is closed,
	// the response was erroneous or the response command is not the expected type.
	WaitForSpecificResponse(ctx context.Context, expectedType command.Id, timeout time.Duration) (command.Command, error)

	// GetDeviceType will return the discovered device type.
	GetDeviceType() DeviceType

	// Close will close the underlying connection and unsubscribe all subscriptions.
	// This method should be called if the communicator is not needed anymore.
	Close() error
}

func waitForResponse(ctx context.Context, deviceType DeviceType, timeout time.Duration, cmdChan chan command.Command, errChan chan error) (command.Command, error) {
	select {
	case cmd := <-cmdChan:
		if cmd.Is(command.IdErrorReport) {
			return nil, Error(cmd, deviceType)
		}

		return cmd, nil
	case err := <-errChan:
		return nil, err
	case <-time.After(timeout):
		return nil, TimeoutErr
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func waitForSpecificResponse(ctx context.Context, deviceType DeviceType, expectedType command.Id, timeout time.Duration, cmdChan chan command.Command, errChan chan error, logPrefix string) (command.Command, error) {
	for {
		select {
		case cmd := <-cmdChan:
			if expectedType != command.IdErrorReport && cmd.Is(command.IdErrorReport) {
				return nil, Error(cmd, deviceType)
			}

			if !cmd.Is(expectedType) {
				if logger.Debug != nil {
					logger.Debug.Printf("%s Unexpected response type: 0x%04x. Skip this command because of waiting for type: 0x%04x",
						logPrefix, cmd.Id(), expectedType,
					)
				}

				continue
			}

			return cmd, nil
		case err := <-errChan:
			return nil, err
		case <-time.After(timeout):
			return nil, TimeoutErr
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
