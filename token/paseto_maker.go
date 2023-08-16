package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// PasteoMaker is a real implementation of the TokenMaker interface.
type PasteoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasteoMaker creates a new PasteoMaker.
func NewPasteoMaker(symmetricKey string) (TokenMaker, error) {
	if len(symmetricKey) < chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasteoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil

}

// CreateToken creates a new token for a specific username and duration.
func (maker *PasteoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

// VerifyToken checks if the token is valid or not.
func (maker *PasteoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// check if the token is expired or not
	if payload.ExpireAt.Before(time.Now()) {
		return nil, fmt.Errorf("token is expired")
	}

	return payload, nil
}
