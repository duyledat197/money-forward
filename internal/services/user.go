package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"user-management/internal/entities"
	"user-management/internal/repositories"
	"user-management/pkg/crypto_utils"
	"user-management/pkg/database"
	"user-management/pkg/id_utils"
	"user-management/pkg/postgres_client"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

type UserService interface {
	CreateUser(context.Context, *entities.User) (int64, error)
	GetUserByID(context.Context, int64) (*entities.User, error)
}

// userService is a representation of service that implements business logic for user domain.
// For using skeleton: s *userService UserService
type userService struct {
	pgClient    *postgres_client.PostgresClient
	idGenerator id_utils.IDGenerator

	// using memories cache for user entity
	userCache *expirable.LRU[int64, *entities.User]

	userRepo interface {
		Create(ctx context.Context, db database.Executor, data *entities.User) error
		GetUserByID(ctx context.Context, db database.Executor, id int64) (*entities.User, error)
		GetUserByUserName(ctx context.Context, db database.Executor, userName string) (*entities.User, error)
	}
}

func NewUserService(pgClient *postgres_client.PostgresClient, idGenerator id_utils.IDGenerator) UserService {
	return &userService{
		pgClient:    pgClient,
		idGenerator: idGenerator,
		userCache:   expirable.NewLRU[int64, *entities.User](128, nil, 24*time.Hour),

		// for repositories
		userRepo: repositories.NewUserRepository(),
	}
}

func (s *userService) CreateUser(ctx context.Context, data *entities.User) (int64, error) {
	// If user exists we should return an existed user error
	existedUser, err := s.userRepo.GetUserByUserName(ctx, s.pgClient, data.UserName)
	if err != nil {
		return 0, err
	}
	if existedUser != nil {
		return 0, fmt.Errorf("username already exists")
	}
	// If error is not no rows error it means unpredicted error
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	// we should storing a hashed password to user table
	pwd, err := crypto_utils.HashPassword(data.Password)
	if err != nil {
		return 0, err
	}

	data.Password = pwd

	// generate a new id for new users
	data.ID = s.idGenerator.Int64()

	if err := s.userRepo.Create(ctx, s.pgClient, data); err != nil {
		return 0, err
	}

	return data.ID, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*entities.User, error) {
	// If user exists in cache, we no need call to database.
	if data, ok := s.userCache.Get(id); ok {
		return data, nil
	}

	data, err := s.userRepo.GetUserByID(ctx, s.pgClient, id)
	if err != nil {
		return nil, err
	}

	// cache user entity by id
	s.userCache.Add(id, data)

	return data, nil
}
