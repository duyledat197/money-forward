package token_utils

import "time"

type Authenticator[T Claims] interface {
	Generate(payload T, expirationTime time.Duration) (string, error)
	Verify(token string) (T, error)
}

type Claims interface {
	Valid() error
	AddExpired(time.Duration)
}
