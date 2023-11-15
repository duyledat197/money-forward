package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"user-management/internal/entities"
	"user-management/internal/repositories"
	"user-management/pkg/cache"
	"user-management/pkg/crypto_utils"
	"user-management/pkg/database"
	"user-management/pkg/http_server/xcontext"
	"user-management/pkg/id_utils"
	"user-management/pkg/postgres_client"
	"user-management/pkg/token_utils"
)

// AuthService is a auth service exporter to used for other layers.
type AuthService interface {
	Login(context.Context, *entities.User) (*entities.User, string, error)
}

// authService is a representation of service that implements business logic for auth domain.
type authService struct {
	pgClient    *postgres_client.PostgresClient
	idGenerator id_utils.IDGenerator

	tknGenerator token_utils.Authenticator[*xcontext.UserInfo]

	userByUserNameCache cache.Cache[string, *entities.User]

	userRepo interface {
		GetUserByUserName(ctx context.Context, db database.Executor, authName string) (*entities.User, error)
	}
}

func NewAuthService(
	pgClient *postgres_client.PostgresClient,
	idGenerator id_utils.IDGenerator,
	tknGenerator token_utils.Authenticator[*xcontext.UserInfo],
	userByUserNameCache cache.Cache[string, *entities.User],

) AuthService {
	return &authService{
		pgClient:     pgClient,
		idGenerator:  idGenerator,
		tknGenerator: tknGenerator,

		userByUserNameCache: userByUserNameCache,

		// for repositories
		userRepo: repositories.NewUserRepository(),
	}
}
func (s *authService) Login(ctx context.Context, req *entities.User) (*entities.User, string, error) {

	user, _ := s.userByUserNameCache.Get(ctx, req.UserName)
	if user == nil {
		var err error
		user, err = s.userRepo.GetUserByUserName(ctx, s.pgClient, req.UserName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, "", fmt.Errorf("username or password is not correctly")
			}
			return nil, "", err
		}
	}

	if err := crypto_utils.CheckPassword(req.Password, user.Password); err != nil {
		return nil, "", err
	}

	tkn, err := s.tknGenerator.Generate(&xcontext.UserInfo{
		UserID: user.ID,
		Role:   string(user.Role),
	}, 24*time.Hour)
	if err != nil {
		return nil, "", err
	}

	return user, tkn, nil
}
