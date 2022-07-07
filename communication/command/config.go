package command

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

type ConfigCommand Command
type ConfigType uint8

const (
	ConfigTypeUnknown   = ConfigType(0)
	ConfigTypeSmartLock = ConfigType(1)
	ConfigTypeOpener    = ConfigType(2)
)

type ConfigSmartLockCommand Command
type ConfigOpenerCommand Command

type OpenerCapabilities uint8

const (
	OpenerCapabilitiesOnlyDoorOpening = OpenerCapabilities(0x00)
	OpenerCapabilitiesBoth            = OpenerCapabilities(0x01)
	OpenerCapabilitiesOnlyRto         = OpenerCapabilities(0x02)
)

type DaylightSavingTimeMode uint8

const (
	DaylightSavingTimeModeDisabled = DaylightSavingTimeMode(0x00)
	DaylightSavingTimeModeEuropean = DaylightSavingTimeMode(0x01)
	DaylightSavingTimeModeUnknown  = DaylightSavingTimeMode(0xFF)
)

type OpenerOperationMode uint8

const (
	OpenerOperationModeGenericDoorOpener          = OpenerOperationMode(0x00)
	OpenerOperationModeAnalogueIntercom           = OpenerOperationMode(0x01)
	OpenerOperationModeDigitalIntercom            = OpenerOperationMode(0x02)
	OpenerOperationModeDigitalIntercomSiedle      = OpenerOperationMode(0x03)
	OpenerOperationModeDigitalIntercomTCS         = OpenerOperationMode(0x04)
	OpenerOperationModeDigitalIntercomBticino     = OpenerOperationMode(0x05)
	OpenerOperationModeAnalogIntercomSiedleHTS    = OpenerOperationMode(0x06)
	OpenerOperationModeDigitalIntercomSTR         = OpenerOperationMode(0x07)
	OpenerOperationModeDigitalIntercomRitto       = OpenerOperationMode(0x08)
	OpenerOperationModeDigitalIntercomFermax      = OpenerOperationMode(0x09)
	OpenerOperationModeDigitalIntercomComelit     = OpenerOperationMode(0x0A)
	OpenerOperationModeDigitalIntercomUrmetBiBus  = OpenerOperationMode(0x0B)
	OpenerOperationModeDigitalIntercomUrmet2Voice = OpenerOperationMode(0x0C)
	OpenerOperationModeDigitalIntercomGolmar      = OpenerOperationMode(0x0D)
	OpenerOperationModeDigitalIntercomSKS         = OpenerOperationMode(0x0E)
	OpenerOperationModeDigitalIntercomSpare       = OpenerOperationMode(0x0F)
)

type AdvertisingMode uint8

const (
	AdvertisingModeAutomatic = AdvertisingMode(0x00)
	AdvertisingModeNormal    = AdvertisingMode(0x01)
	AdvertisingModeSlow      = AdvertisingMode(0x02)
	AdvertisingModeSlowest   = AdvertisingMode(0x03)
)

type HomeKitStatus uint8

const (
	HomeKitStatusNotAvailable     = HomeKitStatus(0x00)
	HomeKitStatusDisabled         = HomeKitStatus(0x01)
	HomeKitStatusEnabled          = HomeKitStatus(0x02)
	HomeKitStatusEnabledAndPaired = HomeKitStatus(0x03)
)

func NewRequestConfig(nonce []byte) Command {
	return NewCommand(IdRequestConfig, nonce)
}

func (c Command) AsConfigCommand() ConfigCommand {
	if !c.Is(IdConfig) {
		return nil
	}

	return ConfigCommand(c)
}

func (c ConfigCommand) AsSmartLockConfig() ConfigSmartLockCommand {
	if c.Type() != ConfigTypeSmartLock {
		return nil
	}
	return ConfigSmartLockCommand(c)
}

func (c ConfigCommand) AsOpenerConfig() ConfigOpenerCommand {
	if c.Type() != ConfigTypeOpener {
		return nil
	}
	return ConfigOpenerCommand(c)
}

func (c ConfigCommand) Type() ConfigType {
	if len(Command(c).Payload()) == 74 {
		return ConfigTypeSmartLock
	} else if len(Command(c).Payload()) == 72 {
		return ConfigTypeOpener
	}
	return ConfigTypeUnknown
}

func (c ConfigCommand) NukiId() uint32 {
	return binary.LittleEndian.Uint32(Command(c).Payload()[0:4])
}

func (c ConfigCommand) Name() string {
	return string(Command(c).Payload()[4:36])
}

func (c ConfigCommand) Latitude() float32 {
	return math.Float32frombits(
		binary.LittleEndian.Uint32(Command(c).Payload()[36:40]),
	)
}

func (c ConfigCommand) Longitude() float32 {
	return math.Float32frombits(
		binary.LittleEndian.Uint32(Command(c).Payload()[40:44]),
	)
}

func (c ConfigOpenerCommand) Capabilities() OpenerCapabilities {
	return OpenerCapabilities(Command(c).Payload()[44])
}

func (c ConfigSmartLockCommand) AutoUnlatch() bool {
	return Command(c).Payload()[44] != 0
}

func (c ConfigCommand) PairingEnabled() bool {
	return Command(c).Payload()[45] != 0
}

func (c ConfigCommand) ButtonEnabled() bool {
	return Command(c).Payload()[46] != 0
}

func (c ConfigCommand) LEDEnabled() bool {
	return Command(c).Payload()[47] != 0
}

func (c ConfigSmartLockCommand) LEDBrightness() uint8 {
	return Command(c).Payload()[48]
}

func (c ConfigCommand) CurrentTimeAsUTC() time.Time {
	var result time.Time
	if c.Type() == ConfigTypeSmartLock {
		result = time.Date(
			int(binary.LittleEndian.Uint16(Command(c).Payload()[49:51])),
			time.Month(Command(c).Payload()[51]),
			int(Command(c).Payload()[52]),
			int(Command(c).Payload()[53]),
			int(Command(c).Payload()[54]),
			int(Command(c).Payload()[55]),
			0,
			time.UTC,
		)
	} else if c.Type() == ConfigTypeOpener {
		result = time.Date(
			int(binary.LittleEndian.Uint16(Command(c).Payload()[48:50])),
			time.Month(Command(c).Payload()[50]),
			int(Command(c).Payload()[51]),
			int(Command(c).Payload()[52]),
			int(Command(c).Payload()[53]),
			int(Command(c).Payload()[54]),
			0,
			time.UTC,
		)
	}

	return result
}

func (c ConfigCommand) TimezoneOffset() time.Duration {
	var offsetInMin uint16
	if c.Type() == ConfigTypeSmartLock {
		offsetInMin = binary.LittleEndian.Uint16(Command(c).Payload()[56:58])
	} else if c.Type() == ConfigTypeOpener {
		offsetInMin = binary.LittleEndian.Uint16(Command(c).Payload()[55:57])
	}

	return time.Duration(int16(offsetInMin)) * time.Minute
}

func (c ConfigCommand) DSTMode() DaylightSavingTimeMode {
	if c.Type() == ConfigTypeSmartLock {
		return DaylightSavingTimeMode(Command(c).Payload()[58])
	} else if c.Type() == ConfigTypeOpener {
		return DaylightSavingTimeMode(Command(c).Payload()[57])
	}

	return DaylightSavingTimeModeUnknown
}

func (c ConfigCommand) HasFob() bool {
	if c.Type() == ConfigTypeSmartLock {
		return Command(c).Payload()[59] != 0x00
	} else if c.Type() == ConfigTypeOpener {
		return Command(c).Payload()[58] != 0x00
	}

	return false
}

func (c ConfigCommand) FobAction1() uint8 {
	if c.Type() == ConfigTypeSmartLock {
		return Command(c).Payload()[60]
	} else if c.Type() == ConfigTypeOpener {
		return Command(c).Payload()[59]
	}

	return 0x00
}

func (c ConfigCommand) FobAction2() uint8 {
	if c.Type() == ConfigTypeSmartLock {
		return Command(c).Payload()[61]
	} else if c.Type() == ConfigTypeOpener {
		return Command(c).Payload()[60]
	}

	return 0x00
}

func (c ConfigCommand) FobAction3() uint8 {
	if c.Type() == ConfigTypeSmartLock {
		return Command(c).Payload()[62]
	} else if c.Type() == ConfigTypeOpener {
		return Command(c).Payload()[61]
	}

	return 0x00
}

func (c ConfigOpenerCommand) OperationMode() OpenerOperationMode {
	return OpenerOperationMode(Command(c).Payload()[62])
}

func (c ConfigSmartLockCommand) SingleLock() bool {
	return Command(c).Payload()[63] != 0x00
}

func (c ConfigCommand) AdvertisingMode() AdvertisingMode {
	if c.Type() == ConfigTypeSmartLock {
		return AdvertisingMode(Command(c).Payload()[64])
	} else if c.Type() == ConfigTypeOpener {
		return AdvertisingMode(Command(c).Payload()[63])
	}

	return 0x00
}

func (c ConfigCommand) HasKeypad() bool {
	if c.Type() == ConfigTypeSmartLock {
		return Command(c).Payload()[65] != 0x00
	} else if c.Type() == ConfigTypeOpener {
		return Command(c).Payload()[64] != 0x00
	}

	return false
}

func (c ConfigCommand) FirmwareVersion() Version {
	if c.Type() == ConfigTypeSmartLock {
		return Command(c).Payload()[66:69]
	} else if c.Type() == ConfigTypeOpener {
		return Command(c).Payload()[65:68]
	}

	return nil
}

func (c ConfigCommand) HardwareRevision() Version {
	if c.Type() == ConfigTypeSmartLock {
		return Command(c).Payload()[69:71]
	} else if c.Type() == ConfigTypeOpener {
		return Command(c).Payload()[68:70]
	}

	return nil
}

func (c ConfigSmartLockCommand) HomeKitStatus() HomeKitStatus {
	return HomeKitStatus(Command(c).Payload()[71])
}

func (c ConfigCommand) TimeZoneId() TimeZoneId {
	tz := TimeZoneId(0xFFFF)
	if c.Type() == ConfigTypeSmartLock {
		tz = TimeZoneId(binary.LittleEndian.Uint16(Command(c).Payload()[72:74]))
	} else if c.Type() == ConfigTypeOpener {
		tz = TimeZoneId(binary.LittleEndian.Uint16(Command(c).Payload()[70:72]))
	}

	return tz
}

func (c ConfigSmartLockCommand) String() string {
	return fmt.Sprintf("Auto unlatch: %v\nLED brightness: %d\nSingle lock: %v\nHomeKit status: 0x%02x",
		c.AutoUnlatch(),
		c.LEDBrightness(),
		c.SingleLock(),
		c.HomeKitStatus(),
	)
}

func (c ConfigOpenerCommand) String() string {
	return fmt.Sprintf("Capabilities: 0x%02x\nOperation mode: 0x%02x",
		c.Capabilities(),
		c.OperationMode(),
	)
}

func (c ConfigCommand) String() string {
	subPart := ""
	if c.Type() == ConfigTypeSmartLock {
		subPart = c.AsSmartLockConfig().String()
	} else if c.Type() == ConfigTypeOpener {
		subPart = c.AsOpenerConfig().String()
	}

	return fmt.Sprintf("Nuki-ID: %08x\nName: %s\nLatitude: %f\nLongitude: %f\n"+
		"Pairing enabled: %v\nButton enabled: %v\nLED enabled: %v\nCurrent Time (UTC): %s, TZ-Offset: %s\nDST-Mode: 0x%02x\n"+
		"Has Fob: %v\nFob Action#1: 0x%02x\nFob Action#2: 0x%02x\nFob Action#3: 0x%02x\n"+
		"Advertising mode: 0x%02x\nHas Keypad: %v\nFirmware: %s\nHardware: %s\nTimezone: %s\n"+
		"%s",
		c.NukiId(),
		c.Name(),
		c.Latitude(),
		c.Longitude(),
		c.PairingEnabled(),
		c.ButtonEnabled(),
		c.LEDEnabled(),
		c.CurrentTimeAsUTC().Format(time.RFC3339),
		c.TimezoneOffset().String(),
		c.DSTMode(),
		c.HasFob(),
		c.FobAction1(),
		c.FobAction2(),
		c.FobAction3(),
		c.AdvertisingMode(),
		c.HasKeypad(),
		c.FirmwareVersion().String(),
		c.HardwareRevision().String(),
		c.TimeZoneId().String(),
		subPart,
	)
}
