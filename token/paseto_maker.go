package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey string
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size must be at least %d characters", chacha20poly1305.KeySize)
	}

	return &PasetoMaker{
		paseto.NewV2(),
		symmetricKey,
	}, nil
}

func (p *PasetoMaker) CreateToken(username string, duration time.Duration) (token string, err error) {

	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return p.paseto.Encrypt([]byte(p.symmetricKey), payload, nil)
}

func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {

	payload := &Payload{}

	if err := p.paseto.Decrypt(token, []byte(p.symmetricKey), payload, nil); err != nil {
		return &Payload{}, err
	}
	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}
