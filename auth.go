package nuki

import (
	"context"
	"fmt"
	"github.com/kevinburke/nacl"
	"github.com/tarent/go-nuki/communication"
	"github.com/tarent/go-nuki/communication/command"
)

// Pair will perform a pairing process with the connected device.
// To pair a device, the device must be in pairing mode and a bluetooth connection must be established before.
// See EstablishConnection to establish a connection. The pairing must be done only once. After the successful
// pairing the authId, private- and public-key should be saved.
func (c *Client) Pair(ctx context.Context, privateKey, publicKey nacl.Key, id command.ClientId, idType command.ClientIdType, name string) error {
	if c.client == nil {
		return ConnectionNotEstablishedError
	}

	err := c.gdioCom.Send(command.NewRequest(command.IdPublicKey))
	if err != nil {
		return fmt.Errorf("unable to send request for public key: %w", err)
	}

	pubKeyResp, err := c.gdioCom.WaitForSpecificResponse(ctx, command.IdPublicKey, c.responseTimeout)
	if err != nil {
		return fmt.Errorf("error while waiting for public key response: %w", err)
	}
	nukiPublicKey := pubKeyResp.AsPublicKeyCommand().PublicKey()

	err = c.gdioCom.Send(command.NewPublicKey((*publicKey)[:]))
	if err != nil {
		return fmt.Errorf("error while sending public key to device: %w", err)
	}

	challenge1, err := c.gdioCom.WaitForSpecificResponse(ctx, command.IdChallenge, c.responseTimeout)
	if err != nil {
		return fmt.Errorf("error while waiting for first challenge: %w", err)
	}

	err = c.gdioCom.Send(command.NewAuthorizationAuthenticator(
		challenge1.AsChallengeCommand().Nonce(),
		nukiPublicKey,
		(*privateKey)[:],
		(*publicKey)[:],
	))
	if err != nil {
		return fmt.Errorf("error while sending authorization authenticator: %w", err)
	}

	challenge2, err := c.gdioCom.WaitForSpecificResponse(ctx, command.IdChallenge, c.responseTimeout)
	if err != nil {
		return fmt.Errorf("error while waiting for second challenge: %w", err)
	}

	err = c.gdioCom.Send(command.NewAuthorizationData(
		challenge2.AsChallengeCommand().Nonce(),
		nukiPublicKey,
		(*privateKey)[:],
		id,
		idType,
		name,
	))
	if err != nil {
		return fmt.Errorf("error while seinding authorization data: %w", err)
	}

	authIdResp, err := c.gdioCom.WaitForSpecificResponse(ctx, command.IdAuthorizationID, c.responseTimeout)
	if err != nil {
		return fmt.Errorf("error while waiting for authorization id: %w", err)
	}

	//TODO: verify authenticator!

	err = c.gdioCom.Send(command.NewAuthorizationIdConfirmation(
		authIdResp.AsAuthorizationIdCommand().Nonce(),
		nukiPublicKey,
		(*privateKey)[:],
		authIdResp.AsAuthorizationIdCommand().AuthorizationId(),
	))
	if err != nil {
		return fmt.Errorf("error while sending authorization id confirmation: %w", err)
	}

	status, err := c.gdioCom.WaitForSpecificResponse(ctx, command.IdStatus, c.responseTimeout)
	if err != nil {
		return fmt.Errorf("error while waiting authorization id confirmation response: %w", err)
	}

	if !status.AsStatusCommand().IsComplete() {
		return fmt.Errorf("pairing failed unexpectedly: status is not completed")
	}

	//done
	err = c.Authenticate(privateKey, publicKey, nukiPublicKey, authIdResp.AsAuthorizationIdCommand().AuthorizationId())
	if err != nil {
		return fmt.Errorf("error while authenticate: %w", err)
	}

	return nil
}

// Authenticate will use the given authentication data and use them for further communication to nuki device.
// The data should be the same which is used for pairing before.
func (c *Client) Authenticate(privateKey, publicKey nacl.Key, nukiPublicKey []byte, authId command.AuthorizationId) error {
	c.privateKey = privateKey
	c.publicKey = publicKey
	c.nukiPublicKey = nukiPublicKey
	c.authId = authId

	if c.client == nil {
		return nil
	}

	var err error
	c.udioCom, err = communication.NewUserSpecificDataIOCommunicator(
		c.client,
		uint32(authId),
		(*privateKey)[:],
		nukiPublicKey,
	)
	if err != nil {
		return err
	}

	return nil
}

// AuthenticationId will return the authId which was generated after the pairing process. See Pair.
func (c *Client) AuthenticationId() command.AuthorizationId {
	return c.authId
}

// PublicKey will return the public key of the connected nuki device.
func (c *Client) PublicKey() []byte {
	return c.nukiPublicKey
}
