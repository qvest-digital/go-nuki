package command

import (
	"encoding/binary"
	"fmt"
	"time"
)

type NukiState uint8
type LockState uint8
type Trigger uint8
type DoorSensorState uint8

const (
	NukiStateUninitialized   = NukiState(0x00)
	NukiStatePairingMode     = NukiState(0x01)
	NukiStateDoorMode        = NukiState(0x02)
	NukiStateContinuousMode  = NukiState(0x03)
	NukiStateMaintenanceMode = NukiState(0x04)

	LockStateUncalibrated            = LockState(0x00)
	LockStateLocked                  = LockState(0x01)
	LockStateUnlocking               = LockState(0x02)
	LockStateUnlocked                = LockState(0x03)
	LockStateRtoActive               = LockState(0x03)
	LockStateLocking                 = LockState(0x04)
	LockStateUnlatched               = LockState(0x05)
	LockStateOpen                    = LockState(0x05)
	LockStateUnlockedLockAndGoActive = LockState(0x06)
	LockStateUnlatching              = LockState(0x07)
	LockStateOpening                 = LockState(0x07)
	LockStateCalibration             = LockState(0xFC)
	LockStateBootRun                 = LockState(0xFD)
	LockStateMotorBlocked            = LockState(0xFE)
	LockStateUndefined               = LockState(0xFF)

	TriggerSystem    = Trigger(0x00)
	TriggerManual    = Trigger(0x01)
	TriggerButton    = Trigger(0x02)
	TriggerAutomatic = Trigger(0x03)
	TriggerAutoLock  = Trigger(0x06)

	DoorSensorStateUnavailable      = DoorSensorState(0x00)
	DoorSensorStateDeactivated      = DoorSensorState(0x01)
	DoorSensorStateDoorClosed       = DoorSensorState(0x02)
	DoorSensorStateDoorOpened       = DoorSensorState(0x03)
	DoorSensorStateDoorStateUnknown = DoorSensorState(0x04)
	DoorSensorStateCalibrating      = DoorSensorState(0x05)
)

type KeyturnerStatesCommand Command

func (c Command) AsKeyturnerStatesCommand() KeyturnerStatesCommand {
	if !c.Is(IdKeyturnerStates) {
		return nil
	}

	return KeyturnerStatesCommand(c)
}

func (k KeyturnerStatesCommand) NukiState() NukiState {
	return NukiState(Command(k).Payload()[0])
}

func (k KeyturnerStatesCommand) LockState() LockState {
	return LockState(Command(k).Payload()[1])
}

func (k KeyturnerStatesCommand) Trigger() Trigger {
	return Trigger(Command(k).Payload()[2])
}

func (k KeyturnerStatesCommand) CurrentTime() time.Time {
	return time.Date(
		int(binary.LittleEndian.Uint16(Command(k).Payload()[3:5])),
		time.Month(Command(k).Payload()[5]),
		int(Command(k).Payload()[6]),
		int(Command(k).Payload()[7]),
		int(Command(k).Payload()[8]),
		int(Command(k).Payload()[9]),
		0,
		time.FixedZone("nuki", int(int16(binary.LittleEndian.Uint16(Command(k).Payload()[10:12])))),
	)
}

func (k KeyturnerStatesCommand) CriticalBatteryState() (critical bool, charging bool, battery uint8) {
	critical = (Command(k).Payload()[12] & 0b0000_0001) == 0b0000_0001
	charging = (Command(k).Payload()[12] & 0b0000_0010) == 0b0000_0010
	battery = (Command(k).Payload()[12] >> 2) * 2

	return
}

func (k KeyturnerStatesCommand) ConfigUpdateCount() uint8 {
	return Command(k).Payload()[13]
}

func (k KeyturnerStatesCommand) LockAndGoTimer() uint8 {
	return Command(k).Payload()[14]
}

func (k KeyturnerStatesCommand) LastLockAction() LockAction {
	return LockAction(Command(k).Payload()[15])
}

func (k KeyturnerStatesCommand) LastLockActionTrigger() Trigger {
	return Trigger(Command(k).Payload()[16])
}

func (k KeyturnerStatesCommand) LastLockActionCompletionStatus() CompletionStatus {
	return CompletionStatus(Command(k).Payload()[17])
}

func (k KeyturnerStatesCommand) DoorSensorState() DoorSensorState {
	return DoorSensorState(Command(k).Payload()[18])
}

func (k KeyturnerStatesCommand) NightModeActive() bool {
	return Command(k).Payload()[19] != 0
}

func (k KeyturnerStatesCommand) AccessoryBatteryState() (supported bool, kpBatteryCritical bool) {
	supported = (Command(k).Payload()[20] & 0b0000_0001) == 0b0000_0001
	kpBatteryCritical = (Command(k).Payload()[20] & 0b0000_0010) == 0b0000_0010

	return
}

func (k KeyturnerStatesCommand) String() string {
	batCritical, batCharge, batPercentage := k.CriticalBatteryState()
	kpSupported, kpBatCritical := k.AccessoryBatteryState()

	return fmt.Sprintf(
		"Nuki-State: %02x\nLock-State: %02x\nTrigger: %02x\nCurrent Time: %s\nBattery:\n\tCritical: %v\n\tCharging: %v\n\tPercentage: %d"+
			"\nConfig update count: %d\nLock 'n' Go timer: %d\nLast Lockaction: %02x\nLast Lockaction trigger: %02x\nLast lockaction completion status: %02x"+
			"\nDoor sensor state: %02x\nNightmode: %v\nAccessory Battery:\n\tSupported: %v\n\tBattery critical: %v",
		k.NukiState(),
		k.LockState(),
		k.Trigger(),
		k.CurrentTime().Format(time.RFC3339),
		batCritical, batCharge, batPercentage,
		k.ConfigUpdateCount(),
		k.LockAndGoTimer(),
		k.LastLockAction(),
		k.LastLockActionTrigger(),
		k.LastLockActionCompletionStatus(),
		k.DoorSensorState(),
		k.NightModeActive(),
		kpSupported, kpBatCritical,
	)
}
