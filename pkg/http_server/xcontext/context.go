package xcontext

import (
	"context"
	"fmt"
	"time"
)

type UserInfo struct {
	UserID    int64     `json:"user_id"`
	Role      string    `json:"role"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (p *UserInfo) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return fmt.Errorf("token has been expired")
	}

	return nil
}

func (p *UserInfo) AddExpired(expirationTime time.Duration) {
	p.ExpiredAt = time.Now().Add(expirationTime)
}

// ImportUserInfoToContext implements import the user info which retrieved from token
// and inject it into the given context.
func ImportUserInfoToContext(ctx context.Context, info *UserInfo) context.Context {
	return context.WithValue(ctx, &userInfoKey{}, info)
}

// ExtractUserInfoFromContext returns an user info which was injected from [ImportUserInfoToContext].
func ExtractUserInfoFromContext(ctx context.Context) (*UserInfo, error) {
	info, ok := ctx.Value(&userInfoKey{}).(*UserInfo)

	if !ok || info == nil {
		return nil, fmt.Errorf("authorization is not valid")
	}

	return info, nil
}
