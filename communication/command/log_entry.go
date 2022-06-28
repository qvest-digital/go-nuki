package command

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

type LoggingType uint8

const (
	LoggingTypeLoggingEnabledDisabled           = LoggingType(0x01)
	LoggingTypeLockAction                       = LoggingType(0x02)
	LoggingTypeCalibration                      = LoggingType(0x03)
	LoggingTypeInitializationRun                = LoggingType(0x04)
	LoggingTypeKeypadAction                     = LoggingType(0x05)
	LoggingTypeDoorSensor                       = LoggingType(0x06)
	LoggingTypeDoorSensorLoggingEnabledDisabled = LoggingType(0x07)
)

type LogEntryCommand Command
type LogEntryLogging Command
type LogEntryLockAction Command
type LogEntryKeypadAction Command
type LogEntryDoorSensor Command
type LogEntryDoorSensorLogging Command

func (c Command) AsLogEntryCommand() LogEntryCommand {
	if !c.Is(IdLogEntry) {
		return nil
	}

	return LogEntryCommand(c)
}

func (l LogEntryCommand) Index() uint32 {
	return binary.LittleEndian.Uint32(Command(l).Payload()[0:4])
}

func (l LogEntryCommand) Timestamp() time.Time {
	return time.Date(
		int(binary.LittleEndian.Uint16(Command(l).Payload()[4:6])),
		time.Month(Command(l).Payload()[6]),
		int(Command(l).Payload()[7]),
		int(Command(l).Payload()[8]),
		int(Command(l).Payload()[9]),
		int(Command(l).Payload()[10]),
		0,
		time.UTC,
	)
}

func (l LogEntryCommand) AuthId() uint32 {
	return binary.LittleEndian.Uint32(Command(l).Payload()[11:15])
}

func (l LogEntryCommand) Name() string {
	return string(Command(l).Payload()[15:47])
}

func (l LogEntryCommand) Type() LoggingType {
	return LoggingType(Command(l).Payload()[47])
}

func (l LogEntryCommand) String() string {
	var part string
	var logType string
	switch l.Type() {
	case LoggingTypeLoggingEnabledDisabled:
		part = l.AsLogging().String()
		logType = "Logging Enabled/Disabled"
	case LoggingTypeLockAction:
		part = l.AsLockAction().String()
		logType = "Lock action"
	case LoggingTypeCalibration:
		part = l.AsLockAction().String()
		logType = "Calibration"
	case LoggingTypeInitializationRun:
		part = l.AsLockAction().String()
		logType = "Initialization run"
	case LoggingTypeKeypadAction:
		part = l.AsKeypadAction().String()
		logType = "Keypad action"
	case LoggingTypeDoorSensor:
		part = l.AsDoorSensor().String()
		logType = "Door sensor"
	case LoggingTypeDoorSensorLoggingEnabledDisabled:
		part = l.AsDoorSensorLogging().String()
		logType = "Door sensor logging"
	default:
		part = hex.EncodeToString(Command(l).Payload()[48:])
		logType = "Unknown"
	}

	return fmt.Sprintf("[%d][%s][%d][%s]: %s > %s",
		l.Index(),
		l.Timestamp().Format(time.RFC3339),
		l.AuthId(),
		l.Name(),
		logType,
		part,
	)
}

func (l LogEntryCommand) AsLogging() LogEntryLogging {
	if l.Type() != LoggingTypeLoggingEnabledDisabled {
		return nil
	}

	return LogEntryLogging(l)
}

func (l LogEntryLogging) IsLoggingEnabled() bool {
	return Command(l).Payload()[48] == 0x01
}

func (l LogEntryLogging) String() string {
	return fmt.Sprintf("logging enabled: %v",
		l.IsLoggingEnabled(),
	)
}

func (l LogEntryCommand) AsLockAction() LogEntryLockAction {
	if l.Type() != LoggingTypeLockAction &&
		l.Type() != LoggingTypeCalibration &&
		l.Type() != LoggingTypeInitializationRun {
		return nil
	}

	return LogEntryLockAction(l)
}

func (l LogEntryLockAction) LockAction() LockAction {
	return LockAction(Command(l).Payload()[48])
}

func (l LogEntryLockAction) Trigger() Trigger {
	return Trigger(Command(l).Payload()[49])
}

func (l LogEntryLockAction) Flags() uint8 {
	return Command(l).Payload()[50]
}

func (l LogEntryLockAction) CompletionStatus() uint8 {
	return Command(l).Payload()[51]
}

func (l LogEntryLockAction) String() string {
	return fmt.Sprintf("LockAction: 0x%02x; Trigger: 0x%02x; Flags: 0x%02x; Completion status: 0x%02x",
		l.LockAction(),
		l.Trigger(),
		l.Flags(),
		l.CompletionStatus(),
	)
}

func (l LogEntryCommand) AsKeypadAction() LogEntryKeypadAction {
	if l.Type() != LoggingTypeKeypadAction {
		return nil
	}

	return LogEntryKeypadAction(l)
}

func (l LogEntryKeypadAction) LockAction() LockAction {
	return LockAction(Command(l).Payload()[48])
}

func (l LogEntryKeypadAction) Source() uint8 {
	return Command(l).Payload()[49]
}

func (l LogEntryKeypadAction) CompletionStatus() uint8 {
	return Command(l).Payload()[50]
}

func (l LogEntryKeypadAction) CodeId() uint16 {
	return binary.LittleEndian.Uint16(Command(l).Payload()[51:53])
}

func (l LogEntryKeypadAction) String() string {
	return fmt.Sprintf("LockAction: 0x%02x; Source: 0x%02x; Completion status: 0x%02x; CodeId: %d",
		l.LockAction(),
		l.Source(),
		l.CompletionStatus(),
		l.CodeId(),
	)
}

func (l LogEntryCommand) AsDoorSensor() LogEntryDoorSensor {
	if l.Type() != LoggingTypeDoorSensor {
		return nil
	}

	return LogEntryDoorSensor(l)
}

func (l LogEntryDoorSensor) IsDoorOpened() bool {
	return Command(l).Payload()[48] == 0x00
}

func (l LogEntryDoorSensor) IsDoorClosed() bool {
	return Command(l).Payload()[48] == 0x01
}

func (l LogEntryDoorSensor) IsSensorJammed() bool {
	return Command(l).Payload()[48] == 0x02
}

func (l LogEntryDoorSensor) String() string {
	return fmt.Sprintf("Door opened: %v; Door closed: %v; Sensor jammed: %v",
		l.IsDoorOpened(),
		l.IsDoorClosed(),
		l.IsSensorJammed(),
	)
}

func (l LogEntryCommand) AsDoorSensorLogging() LogEntryDoorSensorLogging {
	if l.Type() != LoggingTypeDoorSensorLoggingEnabledDisabled {
		return nil
	}

	return LogEntryDoorSensorLogging(l)
}

func (l LogEntryDoorSensorLogging) IsLoggingEnabled() bool {
	return Command(l).Payload()[48] == 0x00
}

func (l LogEntryDoorSensorLogging) String() string {
	return fmt.Sprintf("logging enabled: %v",
		l.IsLoggingEnabled(),
	)
}
