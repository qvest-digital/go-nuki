package communication

import (
	"context"
	"fmt"
	"github.com/go-ble/ble"
	"github.com/tarent/go-nuki/communication/command"
	"github.com/tarent/go-nuki/logger"
	"time"
)

const (
	charUuidGdioSmartLock = "a92ee101-5501-11e4-916c-0800200c9a66"
	charUuidGdioOpener    = "a92ae101-5501-11e4-916c-0800200c9a66"
)

type gdioCommunicator struct {
	commandChan chan command.Command
	errorChan   chan error

	curCommand command.Command

	client     ble.Client
	gdioChar   *ble.Characteristic
	deviceType DeviceType
}

// NewGeneralDataIOCommunicator establish a new communicator to the "general data io" characteristic to the connected nuki device.
func NewGeneralDataIOCommunicator(client ble.Client) (Communicator, error) {
	com := &gdioCommunicator{
		commandChan: make(chan command.Command),
		errorChan:   make(chan error),
		deviceType:  DeviceTypeUnknown,
	}

	profile, err := client.DiscoverProfile(false)
	if err != nil {
		return nil, fmt.Errorf("unable to discover profile: %w", err)
	}
	gdio := profile.FindCharacteristic(ble.NewCharacteristic(ble.MustParse(charUuidGdioSmartLock)))
	if gdio == nil {
		gdio = profile.FindCharacteristic(ble.NewCharacteristic(ble.MustParse(charUuidGdioOpener)))
		if gdio == nil {
			return nil, fmt.Errorf("unable to find general data input output characteristic")
		}
		com.deviceType = DeviceTypeOpener
	} else {
		com.deviceType = DeviceTypeSmartLock
	}

	_, err = client.DiscoverDescriptors(nil, gdio)
	if err != nil {
		return nil, fmt.Errorf("unable to discover general data input output characteristic descriptors: %w", err)
	}

	err = client.Subscribe(gdio, true, com.receive)
	if err != nil {
		return nil, fmt.Errorf("unable to subscribe GDIO: %w", err)
	}

	com.client = client
	com.gdioChar = gdio

	return com, nil
}

func (g *gdioCommunicator) GetDeviceType() DeviceType {
	return g.deviceType
}

func (g *gdioCommunicator) Send(cmd command.Command) error {
	if logger.Info != nil {
		logger.Info.Printf("[GDIO][OUT] %s", cmd.String())
	}
	err := g.client.WriteCharacteristic(g.gdioChar, cmd, false)
	if err != nil {
		return fmt.Errorf("error while send command: %w", err)
	}

	return nil
}

func (g *gdioCommunicator) WaitForResponse(ctx context.Context, timeout time.Duration) (command.Command, error) {
	return waitForResponse(ctx, g.deviceType, timeout, g.commandChan, g.errorChan)
}

func (g *gdioCommunicator) WaitForSpecificResponse(ctx context.Context, expectedType command.Id, timeout time.Duration) (command.Command, error) {
	return waitForSpecificResponse(ctx, g.deviceType, expectedType, timeout, g.commandChan, g.errorChan, "[GDIO][IN]")
}

func (g *gdioCommunicator) receive(payload []byte) {
	if logger.Debug != nil {
		logger.Debug.Printf("[GDIO][IN][PART] %x", payload)
	}

	g.curCommand = append(g.curCommand, payload...)

	if len(payload) == mtu {
		//we expect more data
		return
	}

	//command seems to be completed
	if logger.Info != nil {
		logger.Info.Printf("[GDIO][IN][COMPLETE] %s", g.curCommand.String())
	}

	if !g.curCommand.CheckCRC() {
		g.errorChan <- ERROR_BAD_CRC
		return
	}

	g.commandChan <- g.curCommand
	g.curCommand = []byte{} //clear command
}

func (g *gdioCommunicator) Close() error {
	if err := g.client.Unsubscribe(g.gdioChar, true); err != nil {
		return fmt.Errorf("unable to unsubscribe GDIO: %w", err)
	}

	return nil
}
