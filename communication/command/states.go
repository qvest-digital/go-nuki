package command

import (
	"encoding/binary"
	"fmt"
	"time"
)

type StatesCommand Command
type NukiState uint8
type LockState uint8
type Trigger uint8
type DoorSensorState uint8

const (
	NukiStateUninitialized   = NukiState(0x00)
	NukiStatePairingMode     = NukiState(0x01)
	NukiStateDoorMode        = NukiState(0x02)
	NukiStateMaintenanceMode = NukiState(0x04)

	LockStateUncalibrated = LockState(0x00)
	LockStateLocked       = LockState(0x01)
	LockStateUndefined    = LockState(0xFF)

	TriggerSystem    = Trigger(0x00)
	TriggerManual    = Trigger(0x01)
	TriggerButton    = Trigger(0x02)
	TriggerAutomatic = Trigger(0x03)

	DoorSensorStateUnavailable      = DoorSensorState(0x00)
	DoorSensorStateDeactivated      = DoorSensorState(0x01)
	DoorSensorStateDoorClosed       = DoorSensorState(0x02)
	DoorSensorStateDoorOpened       = DoorSensorState(0x03)
	DoorSensorStateDoorStateUnknown = DoorSensorState(0x04)
	DoorSensorStateCalibrating      = DoorSensorState(0x05)
)

type StatesType uint8

const (
	StatesTypeUnknown   = StatesType(0)
	StatesTypeSmartLock = StatesType(1)
	StatesTypeOpener    = StatesType(2)
)

type StatesSmartLockCommand Command

const (
	LockStateSmartLockUnlocking = LockState(0x02)
	LockStateSmartLockUnlocked  = LockState(0x03)
	LockStateSmartLockLocking   = LockState(0x04)
	LockStateSmartLockUnlatched = LockState(0x05)

	LockStateSmartLockUnlockedLockAndGoActive = LockState(0x06)
	LockStateSmartLockUnlatching              = LockState(0x07)
	LockStateSmartLockCalibration             = LockState(0xFC)
	LockStateSmartLockBootRun                 = LockState(0xFD)
	LockStateSmartLockMotorBlocked            = LockState(0xFE)

	TriggerSmartLockAutoLock = Trigger(0x06)
)

type StatesOpenerCommand Command

const (
	NukiStateOpenerContinuousMode = NukiState(0x03)

	LockStateOpenerRTOActive = LockState(0x03)
	LockStateOpenerOpen      = LockState(0x05)
	LockStateOpenerOpening   = LockState(0x07)
)

func (c Command) AsStatesCommand() StatesCommand {
	if !c.Is(IdKeyturnerStates) && !c.Is(IdOpenerStates) {
		return nil
	}

	return StatesCommand(c)
}

func (s StatesCommand) AsSmartLockStates() StatesSmartLockCommand {
	if s.Type() != StatesTypeSmartLock {
		return nil
	}
	return StatesSmartLockCommand(s)
}

func (s StatesCommand) AsOpenerStates() StatesOpenerCommand {
	if s.Type() != StatesTypeOpener {
		return nil
	}
	return StatesOpenerCommand(s)
}

func (s StatesCommand) Type() StatesType {
	if len(Command(s).Payload()) >= 22 {
		return StatesTypeOpener
	} else if len(Command(s).Payload()) >= 19 {
		return StatesTypeSmartLock
	}
	return StatesTypeUnknown
}

func (s StatesCommand) NukiState() NukiState {
	return NukiState(Command(s).Payload()[0])
}

func (s StatesCommand) LockState() LockState {
	return LockState(Command(s).Payload()[1])
}

func (s StatesCommand) Trigger() Trigger {
	return Trigger(Command(s).Payload()[2])
}

func (s StatesCommand) CurrentTime() time.Time {
	return time.Date(
		int(binary.LittleEndian.Uint16(Command(s).Payload()[3:5])),
		time.Month(Command(s).Payload()[5]),
		int(Command(s).Payload()[6]),
		int(Command(s).Payload()[7]),
		int(Command(s).Payload()[8]),
		int(Command(s).Payload()[9]),
		0,
		time.FixedZone("nuki", int(int16(binary.LittleEndian.Uint16(Command(s).Payload()[10:12])))),
	)
}

func (s StatesSmartLockCommand) CriticalBatteryState() (critical bool, charging bool, battery uint8) {
	critical = (Command(s).Payload()[12] & 0b0000_0001) == 0b0000_0001
	charging = (Command(s).Payload()[12] & 0b0000_0010) == 0b0000_0010
	battery = (Command(s).Payload()[12] >> 2) * 2

	return
}

func (s StatesOpenerCommand) CriticalBatteryState() bool {
	return (Command(s).Payload()[12] & 0b0000_0001) == 0b0000_0001
}

func (s StatesCommand) ConfigUpdateCount() uint8 {
	return Command(s).Payload()[13]
}

func (s StatesSmartLockCommand) LockAndGoTimer() uint8 {
	return Command(s).Payload()[14]
}

func (s StatesOpenerCommand) RingToOpenTimer() uint8 {
	return Command(s).Payload()[14]
}

func (s StatesCommand) LastLockAction() LockAction {
	return LockAction(Command(s).Payload()[15])
}

func (s StatesCommand) LastLockActionTrigger() Trigger {
	return Trigger(Command(s).Payload()[16])
}

func (s StatesCommand) LastLockActionCompletionStatus() CompletionStatus {
	return CompletionStatus(Command(s).Payload()[17])
}

func (s StatesCommand) DoorSensorState() DoorSensorState {
	return DoorSensorState(Command(s).Payload()[18])
}

func (s StatesSmartLockCommand) NightModeActive() bool {
	return Command(s).Payload()[19] != 0
}

func (s StatesSmartLockCommand) AccessoryBatteryState() (supported bool, kpBatteryCritical bool) {
	supported = (Command(s).Payload()[20] & 0b0000_0001) == 0b0000_0001
	kpBatteryCritical = (Command(s).Payload()[20] & 0b0000_0010) == 0b0000_0010

	return
}

func (s StatesCommand) String() string {
	subPart := ""
	if s.Type() == StatesTypeSmartLock {
		subPart = s.AsSmartLockStates().String()
	} else if s.Type() == StatesTypeOpener {
		subPart = s.AsOpenerStates().String()
	}

	return fmt.Sprintf(
		"Nuki-State: 0x%02x\nLock-State: 0x%02x\nTrigger: 0x%02x\nCurrent Time: %s"+
			"\nConfig update count: %d\nLast Lockaction: 0x%02x\nLast Lockaction trigger: 0x%02x\nLast Lockaction completion status: 0x%02x"+
			"\nDoor sensor state: 0x%02x\n%s",
		s.NukiState(),
		s.LockState(),
		s.Trigger(),
		s.CurrentTime().Format(time.RFC3339),
		s.ConfigUpdateCount(),
		s.LastLockAction(),
		s.LastLockActionTrigger(),
		s.LastLockActionCompletionStatus(),
		s.DoorSensorState(),
		subPart,
	)
}

func (s StatesSmartLockCommand) String() string {
	batCritical, batCharge, batPercentage := s.CriticalBatteryState()
	kpSupported, kpBatCritical := s.AccessoryBatteryState()

	return fmt.Sprintf(
		"Battery:\n\tCritical: %v\n\tCharging: %v\n\tPercentage: %d"+
			"\nLock 'n' Go timer: %d\nNightmode: %v\nAccessory Battery:\n\tSupported: %v\n\tBattery critical: %v",
		batCritical, batCharge, batPercentage,
		s.LockAndGoTimer(),
		s.NightModeActive(),
		kpSupported, kpBatCritical,
	)
}

func (s StatesOpenerCommand) String() string {
	return fmt.Sprintf(
		"Battery critical: %v\nRing to Open timer: %d",
		s.CriticalBatteryState(),
		s.RingToOpenTimer(),
	)
}
