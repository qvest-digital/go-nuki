[![Go](https://github.com/tarent/go-nuki/actions/workflows/build.yml/badge.svg)](https://github.com/tarent/go-nuki/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tarent/go-nuki)](https://goreportcard.com/report/github.com/tarent/go-nuki)

# go-nuki

A go library to control [nuki](https://nuki.io) devices - such as Smart Locks - via bluetooth. This library is orienting
on the official [Nuki Smart Lock API V2.2.1](https://developer.nuki.io/page/nuki-smart-lock-api-2/2/)
and [Nuki Opener API v.1.1.1](https://developer.nuki.io/page/nuki-opener-api-1/7/) documentations. Since this library is
based on the [go-ble](https://github.com/go-ble/ble) bluetooth lib, only **Linux** and **Mac OS** are supported
currently.

This lib was successfully tested with a **Nuki Smart Lock 2.0** and **Nuki Opener 2.0** but it should also work with
other Nuki Smart Locks / Opener.

# Features

Actually the following features are implemented:

* [x] Pairing
* [x] Receiving lock status
* [x] Locking
* [x] Unlocking
* [x] Open
* [x] Receive log entries
* [x] Enable/Disable event logging
* [x] Read applied device configuration
* [ ] advanced device configuration
    * [ ] set security pin
    * [x] update time
    * [x] update current time
    * [ ] add keypad codes
* [ ] trigger calibration
* [x] trigger reboot

# Example

```go
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/tarent/go-nuki"
	"github.com/tarent/go-nuki/communication"
	"github.com/tarent/go-nuki/communication/command"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/box"
)

func main() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := nuki.NewClient(device)
	defer nukiClient.Close()

	// the device's MAC address can be found by 'hcitool lescan' for example
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

	// start pairing: the device must be in pairing mode ( normally by pressing the button on the lock for 5 seconds )
	err = nukiClient.Pair(context.Background(), privateKey, publicKey, 13, command.ClientIdTypeApp, "Lib-Nuki-Example")
	if err != nil {
		panic(err)
	}

	//after pairing was successful, save the information into a file or similar
	toSave := map[string]interface{}{
		"authId":     nukiClient.AuthenticationId(),
		"privKey":    hex.EncodeToString((*privateKey)[:]),
		"pubKey":     hex.EncodeToString((*publicKey)[:]),
		"nukiPubKey": hex.EncodeToString(nukiClient.PublicKey()),
	}

	fmt.Printf("Save content:\n%#v", toSave)

	switch nukiClient.GetDeviceType() {
	case communication.DeviceTypeSmartLock:
		err = nukiClient.PerformUnlock(context.Background(), 13)
	case communication.DeviceTypeOpener:
		err = nukiClient.PerformOpen(context.Background(), 13)
	}
	if err != nil {
		panic(err)
	}
}
```

If you already authorized:

```go
package main

import (
	"context"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/tarent/go-nuki"
	"github.com/tarent/go-nuki/communication"
	"github.com/tarent/go-nuki/communication/command"
	"github.com/kevinburke/nacl"
)

func main() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

	nukiClient := nuki.NewClient(device)
	defer nukiClient.Close()

	// the device's MAC address can be found by 'hcitool lescan' for example
	nukiDeviceAddr := ble.NewAddr("54:D2:AA:BB:CC:DD")
	err = nukiClient.EstablishConnection(context.Background(), nukiDeviceAddr)
	if err != nil {
		panic(err)
	}

	//if you already paired before, you have to reuse this informations
	authId := command.AuthorizationId(111111) //load from file
	privateKey := nacl.Key(make([]byte, 32))  //load from file
	publicKey := nacl.Key(make([]byte, 32))   //load from file
	nukiPublicKey := []byte{}                 //load from file

	err = nukiClient.Authenticate(privateKey, publicKey, nukiPublicKey, authId)
	if err != nil {
		panic(err)
	}

	switch nukiClient.GetDeviceType() {
	case communication.DeviceTypeSmartLock:
		err = nukiClient.PerformUnlock(context.Background(), 13)
	case communication.DeviceTypeOpener:
		err = nukiClient.PerformOpen(context.Background(), 13)
	}
	if err != nil {
		panic(err)
	}
}
```

Disable internal logging:

```go
package main

import (
	"github.com/tarent/go-nuki/logger"
)

func main() {
	logger.Debug = nil
	logger.Info = nil

	//...
}
```

For more examples how to use the nuki lib see [examples](example_test.go).

# Raw Device Communication

Is the feature were you looking for not implemented (yet)? Maybe you have look to do it by your own. The two main
"communication-channels" between the client and the device are public and can be used. These channels are also used by
the library itself! In principle, it is quite simple: the client sends a command (byte-payload) and the device sends
commands back. Because the commands is nothing other than byte-slices, you can build the commands by you own (if these
are not available yet.) In the following there is an example how to use these "raw" communication:

```go
package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/kevinburke/nacl/box"
	"github.com/tarent/go-nuki"
	"github.com/tarent/go-nuki/communication/command"
	"os"
	"time"
)

func main() {
	device, err := linux.NewDevice()
	if err != nil {
		panic(err)
	}

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

	//unencrypted communication between application and device
	err = nukiClient.GeneralDataIOCommunicator().Send(command.NewCommand(command.Id(0xaffe), []byte{0x01, 0x02, 0x03, 0x04}))
	if err != nil {
		panic(err)
	}
	response, err := nukiClient.GeneralDataIOCommunicator().WaitForResponse(context.Background(), 10*time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response from device: %s\n", response.String())

	//decrypted communication between application and device
	err = nukiClient.UserSpecificDataIOCommunicator().Send(command.NewCommand(command.Id(0xfeef), []byte{0x01, 0x02, 0x03, 0x04}))
	if err != nil {
		panic(err)
	}
	response, err = nukiClient.UserSpecificDataIOCommunicator().WaitForResponse(context.Background(), 10*time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response from device: %s\n", response.String())
}
```

For more details how the communication of devices will work, look at the api documentations from nuki. Also feel free to
look at the already implemented features to understand how the different communicator will work.

# Test environment

For development/testing purposes there is a docker environment which includes all necessary libraries/tools for
bluetooth.

```shell
docker-compose build  #build container
docker-compose up -d  #start container
```

The deployment will build a binary and copy them inside the docker container and starts a remote-debugging port (**
2345**).

```shell
./deploy.sh "54:D2:72:AA:BB:CC" #build and deploy inside the container
```