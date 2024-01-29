package communication

import (
	"fmt"
	"github.com/go-ble/ble"
)

type DeviceType uint8

const (
	DeviceTypeUnknown   = DeviceType(0x00)
	DeviceTypeSmartLock = DeviceType(0x01)
	DeviceTypeOpener    = DeviceType(0x02)
)

type deviceSpecification struct {
	GeneralDataInputOutputUUID string
	UserDataInputOutputUUID    string
}

var deviceSetups = map[DeviceType]deviceSpecification{
	DeviceTypeSmartLock: {
		GeneralDataInputOutputUUID: "a92ee101-5501-11e4-916c-0800200c9a66",
		UserDataInputOutputUUID:    "a92ee202-5501-11e4-916c-0800200c9a66",
	},
	DeviceTypeOpener: {
		GeneralDataInputOutputUUID: "a92ae101-5501-11e4-916c-0800200c9a66",
		UserDataInputOutputUUID:    "a92ae202-5501-11e4-916c-0800200c9a66",
	},
}

type uuidChooser func(deviceSpecification) string

func chooseGDIO(specification deviceSpecification) string {
	return specification.GeneralDataInputOutputUUID
}

func chooseUDIO(specification deviceSpecification) string {
	return specification.UserDataInputOutputUUID
}

func setupGeneralDataInputOutputCharacteristic(client ble.Client, receiver func(payload []byte)) (char *ble.Characteristic, dType DeviceType, err error) {
	return setupDataInputOutputCharacteristic(client, chooseGDIO, "general data input output", receiver)
}

func setupUserDataInputOutputCharacteristic(client ble.Client, receiver func(payload []byte)) (char *ble.Characteristic, dType DeviceType, err error) {
	return setupDataInputOutputCharacteristic(client, chooseUDIO, "user-specific data input output", receiver)
}

func setupDataInputOutputCharacteristic(client ble.Client, uuidChooser uuidChooser, name string, receiver func(payload []byte)) (char *ble.Characteristic, dType DeviceType, err error) {
	profile, err := client.DiscoverProfile(false)
	if err != nil {
		return nil, DeviceTypeUnknown, fmt.Errorf("unable to discover profile: %w", err)
	}

	for deviceType, setup := range deviceSetups {
		char = profile.FindCharacteristic(ble.NewCharacteristic(ble.MustParse(uuidChooser(setup))))
		if char != nil {
			dType = deviceType
			break
		}
	}
	if char == nil {
		return nil, DeviceTypeUnknown, fmt.Errorf("unable to find " + name + " characteristic")
	}

	_, err = client.DiscoverDescriptors(nil, char)
	if err != nil {
		return nil, dType, fmt.Errorf("unable to discover "+name+" characteristic descriptors: %w", err)
	}

	err = client.Subscribe(char, true, receiver)
	if err != nil {
		return nil, dType, fmt.Errorf("unable to subscribe %s: %w", name, err)
	}

	return char, dType, nil
}
