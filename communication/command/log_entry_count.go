package command

import (
	"encoding/binary"
	"fmt"
)

type LogEntryCountCommand Command

func (c Command) AsLogEntriesCountCommand() LogEntryCountCommand {
	if !c.Is(IdLogEntryCount) {
		return nil
	}

	return LogEntryCountCommand(c)
}

func (l LogEntryCountCommand) IsLoggingEnabled() bool {
	return Command(l).Payload()[0] == 0x01
}

func (l LogEntryCountCommand) Count() uint16 {
	return binary.LittleEndian.Uint16(Command(l).Payload()[1:3])
}

func (l LogEntryCountCommand) IsDoorSensorEnabled() bool {
	return Command(l).Payload()[3] == 0x01
}

func (l LogEntryCountCommand) IsDoorSensorLoggingEnabled() bool {
	return Command(l).Payload()[4] == 0x01
}

func (l LogEntryCountCommand) String() string {
	return fmt.Sprintf("Logging enabled: %v\nCount: %d\nDoor sensor enabled: %v\nDoor sensor logging enabled: %v",
		l.IsLoggingEnabled(),
		l.Count(),
		l.IsDoorSensorEnabled(),
		l.IsDoorSensorLoggingEnabled(),
	)
}
