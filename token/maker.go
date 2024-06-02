package token

import "time"

type Maker interface {
	CreateToken(username string, duration time.Duration) (token string, err error)
	VerifyToken(token string) (*Payload, error)
}
