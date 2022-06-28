package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/kevinburke/nacl/box"
	"github.com/tarent/go-nuki"
	"github.com/tarent/go-nuki/communication"
	"github.com/tarent/go-nuki/communication/command"
	"github.com/tarent/go-nuki/logger"
	"os"
)

func main() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	// disable debug logging
	logger.Debug = nil

	nukiClient := nuki.NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr(os.Args[1])
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	//generate key-pair
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	err = nukiClient.Pair(context.Background(), privateKey, publicKey, 13, command.ClientIdTypeApp, "Go-Nuki-Example")
	if err != nil {
		panic(err)
	}

	//after pairing was successful, save the information into a file or similar
	toSave, _ := json.Marshal(map[string]interface{}{
		"authId":     nukiClient.AuthenticationId(),
		"privKey":    hex.EncodeToString((*privateKey)[:]),
		"pubKey":     hex.EncodeToString((*publicKey)[:]),
		"nukiPubKey": hex.EncodeToString(nukiClient.PublicKey()),
	})

	fmt.Printf("Save content: %s\n", toSave)

	var states interface{ String() string }
	switch nukiClient.GetDeviceType() {
	case communication.DeviceTypeSmartLock:
		states, err = nukiClient.ReadLockerState(context.Background())
	case communication.DeviceTypeOpener:
		states, err = nukiClient.ReadOpenerState(context.Background())
	}

	if err != nil {
		panic(err)
	}
	fmt.Printf("Device-State: %s\n",
		states.String(),
	)

	err = nukiClient.GetLogEntryStream(context.Background(), 0, 0xffff, command.LogSortOrderDescending, "0000", func(log command.LogEntryCommand) {
		fmt.Printf("%s\n", log.String())
	})
	if err != nil {
		panic(err)
	}
}
