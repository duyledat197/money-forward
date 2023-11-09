package services

import (
	"context"
	"user-management/internal/entities"
	"user-management/internal/repositories"
	"user-management/pkg/database"
	"user-management/pkg/postgres_client"
)

type UserService interface {
	CreateUser(context.Context, *entities.User) error
}

// userService is a representation of service that implements business logic for user domain.
// For using skeleton: s *userService UserService
type userService struct {
	pgClient *postgres_client.PostgresClient

	userRepo interface {
		Create(ctx context.Context, db database.Executor, data *entities.User) error
	}
}

func NewUserService(pgClient *postgres_client.PostgresClient) UserService {
	return &userService{
		pgClient: pgClient,
		userRepo: repositories.NewUserRepository(),
	}
}

func (s *userService) CreateUser(ctx context.Context, data *entities.User) error {
	if err := s.userRepo.Create(ctx, s.pgClient, data); err != nil {
		return err
	}

	return nil
}
