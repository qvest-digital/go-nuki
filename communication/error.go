package communication

import (
	"errors"
	"github.com/tarent/go-nuki/communication/command"
)

var (
	ERROR_BAD_CRC    = errors.New("CRC of received command is invalid")
	ERROR_BAD_LENGTH = errors.New("length of retrieved command payload does not match expected length")
	ERROR_UNKNOWN    = errors.New("unknown error")

	P_ERROR_NOT_PAIRING       = errors.New("public key is being requested via request data command, but the Smart Lock is not in pairing mode")
	P_ERROR_BAD_AUTHENTICATOR = errors.New("received authenticator does not match the own calculated authenticator")
	P_ERROR_BAD_PARAMETER     = errors.New("provided parameter is outside of its valid range")
	P_ERROR_MAX_USER          = errors.New("the maximum number of users has been reached")

	K_ERROR_NOT_AUTHORIZED         = errors.New("either the provided authorization id is invalid or the payload could not be decrypted using the shared key for this authorization id")
	K_ERROR_BAD_PIN                = errors.New("the provided pin does not match the stored one")
	K_ERROR_BAD_NONCE              = errors.New("the provided nonce does not match the last stored one of this authorization id or has already been used before")
	K_ERROR_BAD_PARAMETER          = errors.New("a provided parameter is outside of its valid range")
	K_ERROR_INVALID_AUTH_ID        = errors.New("the desired authorization id could not be deleted because it does not exist")
	K_ERROR_DISABLED               = errors.New("the provided authorization id is currently disabled")
	K_ERROR_REMOTE_NOT_ALLOWED     = errors.New("the request has been forwarded by the Nuki Bridge and the provided authorization id has not been granted remote access")
	K_ERROR_TIME_NOT_ALLOWED       = errors.New("the provided authorization id has not been granted access at the current time")
	K_ERROR_TOO_MANY_PIN_ATTEMPTS  = errors.New("an invalid pin has been provided too often")
	K_ERROR_TOO_MANY_ENTRIES       = errors.New("no more entries can be stored")
	K_ERROR_CODE_ALREADY_EXISTS    = errors.New("a Keypad Code should be added but the given code already exists")
	K_ERROR_CODE_INVALID           = errors.New("a Keypad Code that has been entered is invalid")
	K_ERROR_CODE_INVALID_TIMEOUT_1 = errors.New("an invalid pin has been provided multiple times (1)")
	K_ERROR_CODE_INVALID_TIMEOUT_2 = errors.New("an invalid pin has been provided multiple times (2)")
	K_ERROR_CODE_INVALID_TIMEOUT_3 = errors.New("an invalid pin has been provided multiple times (3)")
	K_ERROR_AUTO_UNLOCK_TOO_RECENT = errors.New("an incoming auto unlock request and if a lock action has already been executed within short time")
	K_ERROR_POSITION_UNKNOWN       = errors.New("the request has been forwarded by the Nuki Bridge and the Smart Lock is unsure about its actual lock position")
	K_ERROR_MOTOR_BLOCKED          = errors.New("the motor blocks")
	K_ERROR_CLUTCH_FAILURE         = errors.New("there is a problem with the clutch during motor movement")
	K_ERROR_MOTOR_TIMEOUT          = errors.New("the motor moves for a given period of time but did not block")
	K_ERROR_BUSY                   = errors.New("there is already a lock action processing")
	K_ERROR_CANCELED               = errors.New("the user canceled the motor movement by pressing the button")
	K_ERROR_SL_NOT_CALIBRATED      = errors.New("the Smart Lock has not yet been calibrated")
	K_ERROR_OPENER_NOT_CALIBRATED  = errors.New("the Opener is not in operating mode 0x00 and has not yet been trained")
	K_ERROR_MOTOR_POSITION_LIMIT   = errors.New("the internal position database is not able to store any more values")
	K_ERROR_MOTOR_LOW_VOLTAGE      = errors.New("the motor blocks because of low voltage")
	K_ERROR_MOTOR_POWER_FAILURE    = errors.New("power drain during motor movement is zero")
	K_ERROR_CLUTCH_POWER_FAILURE   = errors.New("the power drain during clutch movement is zero")
	K_ERROR_RECORDING_TIMEOUT      = errors.New("the BUS signal recording duration > 30s without receiving a signal")
	K_ERROR_VOLTAGE_TOO_LOW        = errors.New("the battery voltage is too low and a calibration will therefore not be started")
	K_ERROR_LOW_VOLTAGE            = errors.New("operating mode is > 1 and no voltage is detected on BUS connection")
	K_ERROR_FIRMWARE_UPDATE_NEEDED = errors.New("a firmware update is mandatory")
	K_ERROR_OPERATING_MODE_UNKNOWN = errors.New("operating mode is not in the valid range of the firmware")
)

func Error(c command.Command, deviceType DeviceType) error {
	if c.Id() != command.IdErrorReport {
		return nil
	}
	if len(c) < 3 {
		return nil
	}

	switch c[2] {
	case 0xFD:
		return ERROR_BAD_CRC
	case 0xFE:
		return ERROR_BAD_LENGTH
	case 0xFF:
		return ERROR_UNKNOWN
	case 0x10:
		return P_ERROR_NOT_PAIRING
	case 0x11:
		return P_ERROR_BAD_AUTHENTICATOR
	case 0x12:
		return P_ERROR_BAD_PARAMETER
	case 0x13:
		return P_ERROR_MAX_USER
	case 0x20:
		return K_ERROR_NOT_AUTHORIZED
	case 0x21:
		return K_ERROR_BAD_PIN
	case 0x22:
		return K_ERROR_BAD_NONCE
	case 0x23:
		return K_ERROR_BAD_PARAMETER
	case 0x24:
		return K_ERROR_INVALID_AUTH_ID
	case 0x25:
		return K_ERROR_DISABLED
	case 0x26:
		return K_ERROR_REMOTE_NOT_ALLOWED
	case 0x27:
		return K_ERROR_TIME_NOT_ALLOWED
	case 0x28:
		return K_ERROR_TOO_MANY_PIN_ATTEMPTS
	case 0x29:
		return K_ERROR_TOO_MANY_ENTRIES
	case 0x2A:
		return K_ERROR_CODE_ALREADY_EXISTS
	case 0x2B:
		return K_ERROR_CODE_INVALID
	case 0x2C:
		return K_ERROR_CODE_INVALID_TIMEOUT_1
	case 0x2D:
		return K_ERROR_CODE_INVALID_TIMEOUT_2
	case 0x2E:
		return K_ERROR_CODE_INVALID_TIMEOUT_3
	case 0x40:
		return K_ERROR_AUTO_UNLOCK_TOO_RECENT
	case 0x41:
		return K_ERROR_POSITION_UNKNOWN
	case 0x42:
		return K_ERROR_MOTOR_BLOCKED
	case 0x43:
		return K_ERROR_CLUTCH_FAILURE
	case 0x44:
		return K_ERROR_MOTOR_TIMEOUT
	case 0x45:
		return K_ERROR_BUSY
	case 0x46:
		return K_ERROR_CANCELED
	case 0x47:
		if deviceType == DeviceTypeSmartLock {
			return K_ERROR_SL_NOT_CALIBRATED
		}
		return K_ERROR_OPENER_NOT_CALIBRATED
	case 0x48:
		return K_ERROR_MOTOR_POSITION_LIMIT
	case 0x49:
		if deviceType == DeviceTypeSmartLock {
			return K_ERROR_MOTOR_LOW_VOLTAGE
		}
		return K_ERROR_LOW_VOLTAGE
	case 0x4A:
		return K_ERROR_MOTOR_POWER_FAILURE
	case 0x4B:
		if deviceType == DeviceTypeSmartLock {
			return K_ERROR_CLUTCH_POWER_FAILURE
		}
		return K_ERROR_RECORDING_TIMEOUT
	case 0x4C:
		return K_ERROR_VOLTAGE_TOO_LOW
	case 0x4D:
		return K_ERROR_FIRMWARE_UPDATE_NEEDED
	case 0x50:
		return K_ERROR_OPERATING_MODE_UNKNOWN
	}

	return ERROR_UNKNOWN
}
