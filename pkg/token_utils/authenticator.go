package token_utils

type Authenticator interface {
	Generate(payload *Payload) (*Token, error)
	Verify(token string) (*Payload, error)
}
