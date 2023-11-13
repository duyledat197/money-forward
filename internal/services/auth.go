package services

import (
	"context"

	"user-management/internal/entities"
	"user-management/internal/repositories"
	"user-management/pkg/cache"
	"user-management/pkg/crypto_utils"
	"user-management/pkg/database"
	"user-management/pkg/id_utils"
	"user-management/pkg/postgres_client"
	"user-management/pkg/token_utils"
)

// AuthService is a auth service exporter to used for other layers.
type AuthService interface {
	Login(context.Context, *entities.User) (*entities.User, *token_utils.Token, error)
}

// authService is a representation of service that implements business logic for auth domain.
type authService struct {
	pgClient    *postgres_client.PostgresClient
	idGenerator id_utils.IDGenerator

	tknGenerator token_utils.Authenticator

	userRepo interface {
		GetUserByUserName(ctx context.Context, db database.Executor, authName string) (*entities.User, error)
	}
}

func NewAuthService(
	pgClient *postgres_client.PostgresClient,
	idGenerator id_utils.IDGenerator,
	tknGenerator token_utils.Authenticator,
	userByUserNameCache cache.Cache[int64, *entities.User],

) AuthService {
	return &authService{
		pgClient:     pgClient,
		idGenerator:  idGenerator,
		tknGenerator: tknGenerator,

		// for repositories
		userRepo: repositories.NewUserRepository(),
	}
}
func (s *authService) Login(ctx context.Context, req *entities.User) (*entities.User, *token_utils.Token, error) {
	user, err := s.userRepo.GetUserByUserName(ctx, s.pgClient, req.UserName)
	if err != nil {
		return nil, nil, err
	}

	if err := crypto_utils.CheckPassword(req.Password, user.Password); err != nil {
		return nil, nil, err
	}

	tkn, err := s.tknGenerator.Generate(&token_utils.Payload{
		UserID: user.UserName,
		Role:   string(user.Role),
	})
	if err != nil {
		return nil, nil, err
	}

	return user, tkn, nil
}