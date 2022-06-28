package communication

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/go-ble/ble"
	"github.com/tarent/go-nuki/communication/command"
	"github.com/tarent/go-nuki/logger"
	"time"
)

const (
	charUuidUdioSmartLock = "a92ee202-5501-11e4-916c-0800200c9a66"
	charUuidUdioOpener    = "a92ae202-5501-11e4-916c-0800200c9a66"
)

var DecryptionError = fmt.Errorf("unable to decrypt message")
var UnexpectedAuthId = fmt.Errorf("unexpected authorization id")

type udioCommunicator struct {
	commandChan chan command.Command
	errorChan   chan error

	curEncryptedCommand command.Command
	authId              uint32
	privKey             []byte
	nukiPubKey          []byte

	client     ble.Client
	udioChar   *ble.Characteristic
	deviceType DeviceType
}

// NewUserSpecificDataIOCommunicator establish a new communicator to the "user-specific data io" characteristic to the connected nuki device.
func NewUserSpecificDataIOCommunicator(client ble.Client, authId uint32, userPrivateKey, nukiPublicKey []byte) (Communicator, error) {
	com := &udioCommunicator{
		commandChan: make(chan command.Command),
		errorChan:   make(chan error),
		deviceType:  DeviceTypeUnknown,
		authId:      authId,
		privKey:     userPrivateKey,
		nukiPubKey:  nukiPublicKey,
	}

	profile, err := client.DiscoverProfile(false)
	if err != nil {
		return nil, fmt.Errorf("unable to discover profile: %w", err)
	}

	udio := profile.FindCharacteristic(ble.NewCharacteristic(ble.MustParse(charUuidUdioSmartLock)))
	if udio == nil {
		udio = profile.FindCharacteristic(ble.NewCharacteristic(ble.MustParse(charUuidUdioOpener)))
		if udio == nil {
			return nil, fmt.Errorf("unable to find user-specific data input output characteristic")
		}
		com.deviceType = DeviceTypeOpener
	} else {
		com.deviceType = DeviceTypeSmartLock
	}

	_, err = client.DiscoverDescriptors(nil, udio)
	if err != nil {
		return nil, fmt.Errorf("unable to discover user-specific data input output characteristic descriptors: %w", err)
	}

	err = client.Subscribe(udio, true, com.receive)
	if err != nil {
		return nil, fmt.Errorf("unable to subscribe UDIO: %w", err)
	}

	com.client = client
	com.udioChar = udio

	return com, nil
}

func (u *udioCommunicator) GetDeviceType() DeviceType {
	return u.deviceType
}

func (u *udioCommunicator) Send(cmd command.Command) error {
	if logger.Info != nil {
		logger.Info.Printf("[UDIO][OUT][PLAIN] %s", cmd.String())
	}

	encryptedCmd := command.EncryptCommand(u.authId, u.privKey, u.nukiPubKey, cmd)

	if logger.Debug != nil {
		logger.Debug.Printf("[UDIO][OUT][ENCRYPTED] %s", hex.EncodeToString(encryptedCmd))
	}

	err := u.client.WriteCharacteristic(u.udioChar, encryptedCmd, false)
	if err != nil {
		return fmt.Errorf("error while send command: %w", err)
	}

	return nil
}

func (u *udioCommunicator) WaitForResponse(ctx context.Context, timeout time.Duration) (command.Command, error) {
	return waitForResponse(ctx, u.deviceType, timeout, u.commandChan, u.errorChan)
}

func (u *udioCommunicator) WaitForSpecificResponse(ctx context.Context, expectedType command.Id, timeout time.Duration) (command.Command, error) {
	return waitForSpecificResponse(ctx, u.deviceType, expectedType, timeout, u.commandChan, u.errorChan, "[UDIO][IN]")
}

func (u *udioCommunicator) receive(payload []byte) {
	if logger.Debug != nil {
		logger.Debug.Printf("[UDIO][IN][PART][ENCRYPTED] %x", payload)
	}

	u.curEncryptedCommand = append(u.curEncryptedCommand, payload...)

	if len(payload) == mtu {
		//we expect more data
		return
	}

	//command seems to be completed

	authId, decryptedCommand := command.DecryptCommand(u.curEncryptedCommand, u.privKey, u.nukiPubKey)

	if decryptedCommand == nil {
		u.errorChan <- DecryptionError
		return
	}
	if logger.Info != nil {
		logger.Info.Printf("[UDIO][IN][COMPLETE][PLAIN] %s", decryptedCommand.String())
	}

	if u.authId != authId {
		u.errorChan <- UnexpectedAuthId
		return
	}

	if !decryptedCommand.CheckCRC() {
		u.errorChan <- ERROR_BAD_CRC
		return
	}

	u.commandChan <- decryptedCommand
	u.curEncryptedCommand = []byte{} //clear command
}

func (u *udioCommunicator) Close() error {
	if err := u.client.Unsubscribe(u.udioChar, true); err != nil {
		return fmt.Errorf("unable to unsubscribe UDIO: %w", err)
	}

	return nil
}
