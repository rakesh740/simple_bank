package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrTokenExpired = errors.New("token has expired")
)

func NewPayload(username string, duration time.Duration) (Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Payload{}, err
	}
	payload := Payload{
		id,
		username,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	return payload, nil
}

func (payload Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt.Time) {
		return ErrTokenExpired
	}
	return nil
}
