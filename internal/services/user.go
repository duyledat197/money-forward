package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"user-management/internal/entities"
	"user-management/internal/models"
	"user-management/pkg/crypto_utils"
	"user-management/pkg/database"
	"user-management/pkg/id_utils"
	"user-management/pkg/postgres_client"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

// UserService is an exporter to used for other layers.
type UserService interface {
	CreateUser(context.Context, *entities.User) (int64, error)
	GetUserByID(context.Context, int64) (*entities.User, error)
	Update(ctx context.Context, data *entities.User) error

	// for account
	CreateAccount(ctx context.Context, data *entities.Account) (int64, error)
	ListAccountByID(ctx context.Context, id int64, paging *models.Paging) ([]*entities.Account, error)
}

// userService is a representation of service that implements business logic for user domain.
// For using skeleton: s *userService UserService
type userService struct {
	pgClient    *postgres_client.PostgresClient
	idGenerator id_utils.IDGenerator

	// using memories cache for user entity
	userCache *expirable.LRU[int64, *entities.User]

	userRepo    userRepo
	accountRepo accountRepo
}

func NewUserService(
	pgClient *postgres_client.PostgresClient,
	idGenerator id_utils.IDGenerator,
	userCache *expirable.LRU[int64, *entities.User],
	userRepo userRepo,
	accountRepo accountRepo,
) UserService {
	return &userService{
		pgClient:    pgClient,
		idGenerator: idGenerator,
		userCache:   userCache,

		// for repositories
		userRepo:    userRepo,
		accountRepo: accountRepo,
	}
}

type userRepo interface {
	Create(ctx context.Context, db database.Executor, data *entities.User) error
	UpdateByID(ctx context.Context, db database.Executor, id int64, data *entities.User) error
	GetUserByID(ctx context.Context, db database.Executor, id int64) (*entities.User, error)
	GetUserByUserName(ctx context.Context, db database.Executor, userName string) (*entities.User, error)
	DeleteByID(ctx context.Context, db database.Executor, id int64) error
}

type accountRepo interface {
	Create(ctx context.Context, db database.Executor, data *entities.Account) error
	ListAccountByUserID(ctx context.Context, db database.Executor, userID int64) ([]*entities.Account, error)
}

// CreateUser is implementation to business logic for create user.
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

	// We should storing a hashed password to user table
	pwd, err := crypto_utils.HashPassword(data.Password)
	if err != nil {
		return 0, err
	}

	data.Password = pwd

	// Generate a new id for new users
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

// Update is representation of business logic to update user by id
func (s *userService) Update(ctx context.Context, data *entities.User) error {
	// just update, no need check exists because we will check row affected.
	if err := s.userRepo.UpdateByID(ctx, s.pgClient, data.ID, data); err != nil {
		return err
	}

	// remove from cache because user info changed
	s.userCache.Remove(data.ID)

	return nil
}

// DeleteByID is representation of business logic to update user by id
func (s *userService) DeleteByID(ctx context.Context, id int64) error {
	if err := s.userRepo.DeleteByID(ctx, s.pgClient, id); err != nil {
		return err
	}

	// remove from cache because user info removed
	s.userCache.Remove(id)

	return nil
}

// CreateAccount is implementation to business logic for create account by user id.
func (s *userService) CreateAccount(ctx context.Context, data *entities.Account) (int64, error) {
	// checking use existed
	if _, err := s.userRepo.GetUserByID(ctx, s.pgClient, data.UserID); err != nil {
		// custom exists user error
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user does not exists")
		}

		// throw err for other error
		return 0, err
	}
	// Generate a new id for new accounts
	data.ID = s.idGenerator.Int64()

	err := s.accountRepo.Create(ctx, s.pgClient, data)
	if err != nil {
		return 0, err
	}

	s.userCache.Remove(data.ID)

	return data.ID, nil
}

func (s *userService) ListAccountByID(ctx context.Context, id int64, paging *models.Paging) ([]*entities.Account, error) {
	// checking use existed
	if _, err := s.userRepo.GetUserByID(ctx, s.pgClient, id); err != nil {
		// custom exists user error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user does not exists")
		}

		return nil, err
	}

	accounts, err := s.accountRepo.ListAccountByUserID(ctx, s.pgClient, id)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}
