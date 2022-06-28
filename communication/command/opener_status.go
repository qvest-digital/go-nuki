package command

import (
	"encoding/binary"
	"fmt"
	"time"
)

type OpenerStatesCommand Command

func (c Command) AsOpenerStatesCommand() OpenerStatesCommand {
	if !c.Is(IdOpenerStates) {
		return nil
	}

	return OpenerStatesCommand(c)
}

func (o OpenerStatesCommand) NukiState() NukiState {
	return NukiState(Command(o).Payload()[0])
}

func (o OpenerStatesCommand) LockState() LockState {
	return LockState(Command(o).Payload()[1])
}

func (o OpenerStatesCommand) Trigger() Trigger {
	return Trigger(Command(o).Payload()[2])
}

func (o OpenerStatesCommand) CurrentTime() time.Time {
	return time.Date(
		int(binary.LittleEndian.Uint16(Command(o).Payload()[3:5])),
		time.Month(Command(o).Payload()[5]),
		int(Command(o).Payload()[6]),
		int(Command(o).Payload()[7]),
		int(Command(o).Payload()[8]),
		int(Command(o).Payload()[9]),
		0,
		time.FixedZone("nuki", int(int16(binary.LittleEndian.Uint16(Command(o).Payload()[10:12])))),
	)
}

func (o OpenerStatesCommand) CriticalBatteryState() (critical bool, charging bool, battery uint8) {
	critical = (Command(o).Payload()[12] & 0b0000_0001) == 0b0000_0001
	charging = (Command(o).Payload()[12] & 0b0000_0010) == 0b0000_0010
	battery = (Command(o).Payload()[12] >> 2) * 2

	return
}

func (o OpenerStatesCommand) ConfigUpdateCount() uint8 {
	return Command(o).Payload()[13]
}

func (o OpenerStatesCommand) RingToOpenTimer() uint8 {
	return Command(o).Payload()[14]
}

func (o OpenerStatesCommand) LastLockAction() LockAction {
	return LockAction(Command(o).Payload()[15])
}

func (o OpenerStatesCommand) LastLockActionTrigger() Trigger {
	return Trigger(Command(o).Payload()[16])
}

func (o OpenerStatesCommand) LastLockActionCompletionStatus() CompletionStatus {
	return CompletionStatus(Command(o).Payload()[17])
}

func (o OpenerStatesCommand) DoorSensorState() DoorSensorState {
	return DoorSensorState(Command(o).Payload()[18])
}

func (o OpenerStatesCommand) String() string {
	batCritical, batCharge, batPercentage := o.CriticalBatteryState()

	return fmt.Sprintf(
		"Nuki-State: %02x\nLock-State: %02x\nTrigger: %02x\nCurrent Time: %s\nBattery:\n\tCritical: %v\n\tCharging: %v\n\tPercentage: %d"+
			"\nConfig update count: %d\nRing to Open timer: %d\nLast Lockaction: %02x\nLast Lockaction trigger: %02x\nLast lockaction completion status: %02x"+
			"\nDoor sensor state: %02x",
		o.NukiState(),
		o.LockState(),
		o.Trigger(),
		o.CurrentTime().Format(time.RFC3339),
		batCritical, batCharge, batPercentage,
		o.ConfigUpdateCount(),
		o.RingToOpenTimer(),
		o.LastLockAction(),
		o.LastLockActionTrigger(),
		o.LastLockActionCompletionStatus(),
		o.DoorSensorState(),
	)
}
