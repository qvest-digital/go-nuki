package nuki

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/box"
	"github.com/tarent/go-nuki/communication/command"
	"time"
)

func ExampleClient_EstablishConnection() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_Pair() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	//generate key-pair
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	err = nukiClient.Pair(context.Background(), privateKey, publicKey, 13, command.ClientIdTypeApp, "Lib-Nuki-Example")
	if err != nil {
		panic(err)
	}

	//after pairing was successful, save the information into a file or similar
	toSave := map[string]interface{}{
		"authId":     nukiClient.AuthenticationId(),
		"privKey":    privateKey,
		"pubKey":     publicKey,
		"nukiPubKey": nukiClient.PublicKey(),
	}

	fmt.Printf("Save content:\n%#v", toSave)
}

func ExampleClient_Authenticate() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	authId := command.AuthorizationId(111111) //load from file
	privateKey := nacl.Key(make([]byte, 32))  //load from file
	publicKey := nacl.Key(make([]byte, 32))   //load from file
	nukiPublicKey := []byte{}                 //load from file

	err = nukiClient.Authenticate(privateKey, publicKey, nukiPublicKey, authId)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_ReadLockerState() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	authId := command.AuthorizationId(111111) //load from file
	privateKey := nacl.Key(make([]byte, 32))  //load from file
	publicKey := nacl.Key(make([]byte, 32))   //load from file
	nukiPublicKey := []byte{}                 //load from file

	err = nukiClient.Authenticate(privateKey, publicKey, nukiPublicKey, authId)
	if err != nil {
		panic(err)
	}

	state, err := nukiClient.ReadLockerState(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Status:\n%s", state)
}

func ExampleClient_ReadOpenerState() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	authId := command.AuthorizationId(111111) //load from file
	privateKey := nacl.Key(make([]byte, 32))  //load from file
	publicKey := nacl.Key(make([]byte, 32))   //load from file
	nukiPublicKey := []byte{}                 //load from file

	err = nukiClient.Authenticate(privateKey, publicKey, nukiPublicKey, authId)
	if err != nil {
		panic(err)
	}

	state, err := nukiClient.ReadOpenerState(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Status:\n%s", state)
}

func ExampleClient_PerformAction() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	authId := command.AuthorizationId(111111) //load from file
	privateKey := nacl.Key(make([]byte, 32))  //load from file
	publicKey := nacl.Key(make([]byte, 32))   //load from file
	nukiPublicKey := []byte{}                 //load from file

	err = nukiClient.Authenticate(privateKey, publicKey, nukiPublicKey, authId)
	if err != nil {
		panic(err)
	}

	err = nukiClient.PerformAction(context.Background(), func(nonce []byte) command.Command {
		suffix := "logSuffix"

		return command.NewLockAction(
			command.LockActionLockAndGo,
			13,
			command.LockActionFlagForce|command.LockActionFlagAutoUnlock,
			&suffix,
			nonce,
		)
	})
	if err != nil {
		panic(err)
	}
}

func ExampleClient_PerformLock() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	authId := command.AuthorizationId(111111) //load from file
	privateKey := nacl.Key(make([]byte, 32))  //load from file
	publicKey := nacl.Key(make([]byte, 32))   //load from file
	nukiPublicKey := []byte{}                 //load from file

	err = nukiClient.Authenticate(privateKey, publicKey, nukiPublicKey, authId)
	if err != nil {
		panic(err)
	}

	err = nukiClient.PerformLock(context.Background(), 13)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_PerformUnlock() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	//generate key-pair
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	err = nukiClient.Pair(context.Background(), privateKey, publicKey, 13, command.ClientIdTypeApp, "Lib-Nuki-Example")
	if err != nil {
		panic(err)
	}

	err = nukiClient.PerformUnlock(context.Background(), 13)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_PerformLockAction() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	//generate key-pair
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	err = nukiClient.Pair(context.Background(), privateKey, publicKey, 13, command.ClientIdTypeApp, "Lib-Nuki-Example")
	if err != nil {
		panic(err)
	}

	err = nukiClient.PerformLockAction(context.Background(), 13, command.LockActionLockAndGo)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_GetLogEntriesCount() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	//generate key-pair
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	err = nukiClient.Pair(context.Background(), privateKey, publicKey, 13, command.ClientIdTypeApp, "Lib-Nuki-Example")
	if err != nil {
		panic(err)
	}

	logCount, err := nukiClient.GetLogEntriesCount(context.Background(), "0000")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Count: %s\n", logCount.String())
}

func ExampleClient_GetLogEntries() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	//generate key-pair
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	err = nukiClient.Pair(context.Background(), privateKey, publicKey, 13, command.ClientIdTypeApp, "Lib-Nuki-Example")
	if err != nil {
		panic(err)
	}

	logs, err := nukiClient.GetLogEntries(context.Background(), 0, 10, command.LogSortOrderDescending, "0000")
	if err != nil {
		panic(err)
	}

	for _, log := range logs {
		fmt.Printf("%s\n", log.String())
	}
}

func ExampleClient_GetLogEntryStream() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	//generate key-pair
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	err = nukiClient.Pair(context.Background(), privateKey, publicKey, 13, command.ClientIdTypeApp, "Lib-Nuki-Example")
	if err != nil {
		panic(err)
	}

	err = nukiClient.GetLogEntryStream(context.Background(), 0, 0xffff, command.LogSortOrderDescending, "0000", func(log command.LogEntryCommand) {
		fmt.Printf("%s\n", log.String())
	})
	if err != nil {
		panic(err)
	}
}

func ExampleClient_SetLogging() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	//generate key-pair
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	err = nukiClient.Pair(context.Background(), privateKey, publicKey, 13, command.ClientIdTypeApp, "Lib-Nuki-Example")
	if err != nil {
		panic(err)
	}

	err = nukiClient.SetLogging(context.Background(), "0000", true)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_PerformOpen() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	authId := command.AuthorizationId(111111) //load from file
	privateKey := nacl.Key(make([]byte, 32))  //load from file
	publicKey := nacl.Key(make([]byte, 32))   //load from file
	nukiPublicKey := []byte{}                 //load from file

	err = nukiClient.Authenticate(privateKey, publicKey, nukiPublicKey, authId)
	if err != nil {
		panic(err)
	}

	err = nukiClient.PerformOpen(context.Background(), 13)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_PerformOpenAction() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	authId := command.AuthorizationId(111111) //load from file
	privateKey := nacl.Key(make([]byte, 32))  //load from file
	publicKey := nacl.Key(make([]byte, 32))   //load from file
	nukiPublicKey := []byte{}                 //load from file

	err = nukiClient.Authenticate(privateKey, publicKey, nukiPublicKey, authId)
	if err != nil {
		panic(err)
	}

	err = nukiClient.PerformOpenAction(context.Background(), 13, command.OpenActionActivateRTO)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_UpdateTime() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := NewClient(device)
	defer nukiClient.Close()

	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	authId := command.AuthorizationId(111111) //load from file
	privateKey := nacl.Key(make([]byte, 32))  //load from file
	publicKey := nacl.Key(make([]byte, 32))   //load from file
	nukiPublicKey := []byte{}                 //load from file

	err = nukiClient.Authenticate(privateKey, publicKey, nukiPublicKey, authId)
	if err != nil {
		panic(err)
	}

	err = nukiClient.UpdateTime(context.Background(), "0000", time.Now())
	if err != nil {
		panic(err)
	}
}
