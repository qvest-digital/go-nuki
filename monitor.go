package nuki

import (
	"context"
	"encoding/hex"

	"github.com/go-ble/ble"
)

type StateChangeHandler func(context.Context, *Client)

// MonitorStateChanges listens to advertisements from the given devices and triggers the callback
// once their state changes. The client - which must be paired with ClientIdTypeBridge - can then
// call ReadStates to fetch the new state and reset the flag.
func MonitorStateChanges(ctx context.Context, clb StateChangeHandler, clients ...*Client) error {
	h := func(a ble.Advertisement) {
		for _, c := range clients {
			if a.Addr().String() == c.addr.String() {
				_, stateChanged := ParseAdvertisement(a)
				if stateChanged && !c.stateChanged {
					go clb(ctx, c)
				}
				c.stateChanged = stateChanged
			}
		}
	}
	return ble.Scan(ctx, true, h, nil)
}

func ParseAdvertisement(a ble.Advertisement) (id string, stateChanged bool) {
	md := a.ManufacturerData()
	id = hex.EncodeToString(md[20:24])
	txPower := int8(md[len(md)-1])
	stateChanged = txPower & 1 == 1
	return
}
