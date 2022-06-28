package nuki

import (
	"context"
	"fmt"
	"github.com/tarent/go-nuki/communication/command"
)

// GetLogEntriesCount will return the count of persisting logs.
func (c *Client) GetLogEntriesCount(ctx context.Context, pin string) (command.LogEntryCountCommand, error) {
	pinAsInt, err := c.checkPreconditionAndParsePin(pin)
	if err != nil {
		return nil, err
	}

	err = c.udioCom.Send(command.NewRequest(command.IdChallenge))
	if err != nil {
		return nil, fmt.Errorf("unable to send request for challenge: %w", err)
	}

	challenge, err := c.udioCom.WaitForSpecificResponse(ctx, command.IdChallenge, c.responseTimeout)
	if err != nil {
		return nil, fmt.Errorf("error while waiting for challenge: %w", err)
	}
	err = c.udioCom.Send(command.NewRequestLogEntriesCountCommand(
		pinAsInt,
		challenge.AsChallengeCommand().Nonce(),
	))
	if err != nil {
		return nil, fmt.Errorf("unable to send request for log count: %w", err)
	}

	logEntryCount, err := c.udioCom.WaitForSpecificResponse(ctx, command.IdLogEntryCount, c.responseTimeout)
	if err != nil {
		return nil, fmt.Errorf("error while waiting for log entry count: %w", err)
	}

	status, err := c.udioCom.WaitForSpecificResponse(ctx, command.IdStatus, c.responseTimeout)
	if err != nil {
		return nil, fmt.Errorf("error while waiting for status: %w", err)
	}

	if !status.AsStatusCommand().IsComplete() {
		return nil, fmt.Errorf("unexpected status: expect 0x%02x got 0x%02x", command.CompletionStatusComplete, status.AsStatusCommand().Status())
	}

	return logEntryCount.AsLogEntriesCountCommand(), nil
}

// GetLogEntryStream will start consume the persisted logs from the device. While the callback function will be called
// foreach received log entry. This function is blocking which mean it will return after the log receiving is done.
func (c *Client) GetLogEntryStream(ctx context.Context, start uint32, count uint16, order command.LogSortOrder, pin string, clb func(command.LogEntryCommand)) error {
	pinAsInt, err := c.checkPreconditionAndParsePin(pin)
	if err != nil {
		return err
	}

	err = c.udioCom.Send(command.NewRequest(command.IdChallenge))
	if err != nil {
		return fmt.Errorf("unable to send request for challenge: %w", err)
	}

	challenge, err := c.udioCom.WaitForSpecificResponse(ctx, command.IdChallenge, c.responseTimeout)
	if err != nil {
		return fmt.Errorf("error while waiting for challenge: %w", err)
	}
	err = c.udioCom.Send(command.NewRequestLogEntriesCommand(
		start,
		count,
		order,
		pinAsInt,
		challenge.AsChallengeCommand().Nonce(),
	))
	if err != nil {
		return fmt.Errorf("unable to send request for log entries: %w", err)
	}

	for {
		resp, err := c.udioCom.WaitForResponse(ctx, c.responseTimeout)
		if err != nil {
			return fmt.Errorf("error while waiting for log entry: %w", err)
		}
		if resp.Is(command.IdLogEntry) {
			clb(resp.AsLogEntryCommand())
		} else if resp.Is(command.IdStatus) {
			break //we are done
		} else {
			return fmt.Errorf("unexpected response type")
		}
	}

	return nil
}

// GetLogEntries will return the persisted log entries from the device. All logentries will be saved in memory! For a huge
// load of log entries consider the usage of GetLogEntryStream instead.
func (c *Client) GetLogEntries(ctx context.Context, start uint32, count uint16, order command.LogSortOrder, pin string) ([]command.LogEntryCommand, error) {
	result := make([]command.LogEntryCommand, 0, count)
	err := c.GetLogEntryStream(ctx, start, count, order, pin, func(logEntry command.LogEntryCommand) {
		result = append(result, logEntry)
	})

	return result, err
}

// EnableLogging will enable the logging on the connected nuki device.
func (c *Client) EnableLogging(ctx context.Context, pin string) error {
	return c.SetLogging(ctx, pin, true)
}

// DisableLogging will disable the logging on the connected nuki device.
func (c *Client) DisableLogging(ctx context.Context, pin string) error {
	return c.SetLogging(ctx, pin, false)
}

// SetLogging will set the logging on the connected nuki device.
func (c *Client) SetLogging(ctx context.Context, pin string, enable bool) error {
	pinAsInt, err := c.checkPreconditionAndParsePin(pin)
	if err != nil {
		return err
	}

	return c.PerformAction(ctx, func(nonce []byte) command.Command {
		return command.NewEnableLogging(enable, pinAsInt, nonce)
	})
}
