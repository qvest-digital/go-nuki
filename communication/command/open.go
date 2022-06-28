package command

type OpenAction uint8

const (
	OpenActionActivateRTO             = OpenAction(0x01)
	OpenActionDeactivateRTO           = OpenAction(0x02)
	OpenActionElectricStrikeActuation = OpenAction(0x03)
	OpenActionActivateCm              = OpenAction(0x04)
	OpenActionDeactivateCm            = OpenAction(0x05)
	OpenActionFobAction1              = OpenAction(0x81)
	OpenActionFobAction2              = OpenAction(0x82)
	OpenActionFobAction3              = OpenAction(0x83)

	OpenActionFlagGeofence = 0b0000_0001
	OpenActionFlagForce    = 0b0000_0010
)

func NewOpenAction(action OpenAction, appId uint32, flags uint8, nameSuffix *string, nonce []byte) Command {
	return NewLockAction(LockAction(action), appId, flags, nameSuffix, nonce)
}
