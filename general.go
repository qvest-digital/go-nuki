package nuki

import (
	"context"
	"fmt"
	"github.com/go-ble/ble"
	"github.com/kevinburke/nacl"
	"github.com/tarent/go-nuki/communication"
	"github.com/tarent/go-nuki/communication/command"
	"time"
)

// ConnectionNotEstablishedError will be returned if the connection is not established before
var ConnectionNotEstablishedError = fmt.Errorf("the connection is not established")

// UnauthenticatedError will be returned if the client is not authenticated before
var UnauthenticatedError = fmt.Errorf("the client is not authenticated")

// InvalidPinError will be returned if the given pin is invalid
var InvalidPinError = fmt.Errorf("the given pin is invalid")

type Client struct {
	client          ble.Client
	responseTimeout time.Duration

	privateKey    nacl.Key
	publicKey     nacl.Key
	nukiPublicKey []byte
	authId        command.AuthorizationId

	gdioCom communication.Communicator
	udioCom communication.Communicator
}

func NewClient(bleDevice ble.Device) *Client {
	ble.SetDefaultDevice(bleDevice)

	return &Client{
		responseTimeout: 10 * time.Second,
	}
}

// WithTimeout sets the timeout which is used for each response waiting.
func (c *Client) WithTimeout(duration time.Duration) *Client {
	c.responseTimeout = duration
	return c
}

// EstablishConnection establish a connection to the given nuki device.
// Returns an error if there was a problem with connecting to the device.
func (c *Client) EstablishConnection(ctx context.Context, deviceAddress ble.Addr) error {
	bleClient, err := ble.Dial(ctx, deviceAddress)
	if err != nil {
		return fmt.Errorf("error while establish connection: %w", err)
	}
	c.client = bleClient

	c.gdioCom, err = communication.NewGeneralDataIOCommunicator(bleClient)
	if err != nil {
		return fmt.Errorf("error while establish communication: %w", err)
	}

	return nil
}

// GeneralDataIOCommunicator will return the communicator which is responsible for general data io.
// This is only available after the connection is established (EstablishConnection)
func (c *Client) GeneralDataIOCommunicator() communication.Communicator {
	return c.gdioCom
}

// UserSpecificDataIOCommunicator will return the communicator which is responsible for user specific data io.
// This is only available after the connection is established (EstablishConnection) and the authentication is done (Pair or Authenticate).
func (c *Client) UserSpecificDataIOCommunicator() communication.Communicator {
	return c.udioCom
}

// Close will close all underlying resources. This function should be called after
// the client will not be used anymore.
func (c *Client) Close() error {
	errors := make([]error, 0, 3)

	if c.gdioCom != nil {
		if err := c.gdioCom.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if c.udioCom != nil {
		if err := c.udioCom.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if c.client != nil {
		if err := c.client.Conn().Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("error while closing resources: [%v]", errors)
	}

	return nil
}
