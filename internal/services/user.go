package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"user-management/internal/entities"
	"user-management/internal/repositories"
	"user-management/pkg/cache"
	"user-management/pkg/crypto_utils"
	"user-management/pkg/database"
	"user-management/pkg/http_server/xcontext"
	"user-management/pkg/id_utils"
	"user-management/pkg/postgres_client"
)

// UserService is a service exporter to used for other layers.
type UserService interface {
	CreateUser(context.Context, *entities.User) (int64, error)
	GetUserByID(context.Context, int64) (*entities.UserWithAccounts, error)
	Update(ctx context.Context, data *entities.User) error

	// for account
	CreateAccount(ctx context.Context, data *entities.Account) (int64, error)
	ListAccountByID(ctx context.Context, id int64) ([]*entities.Account, error)
}

// userService is a representation of service that implements business logic for user domain.
type userService struct {
	pgClient    *postgres_client.PostgresClient
	idGenerator id_utils.IDGenerator

	// using memories cache for user entity
	userCache           cache.Cache[int64, *entities.UserWithAccounts]
	userByUserNameCache cache.Cache[string, *entities.User]

	userRepo interface {
		Create(ctx context.Context, db database.Executor, data *entities.User) error
		UpdateByID(ctx context.Context, db database.Executor, id int64, data *entities.User) error
		GetUserByID(ctx context.Context, db database.Executor, id int64) (*entities.UserWithAccounts, error)
		GetUserByUserName(ctx context.Context, db database.Executor, userName string) (*entities.User, error)
		DeleteByID(ctx context.Context, db database.Executor, id int64) error
	}
	accountRepo interface {
		Create(ctx context.Context, db database.Executor, data *entities.Account) error
		ListAccountByUserID(ctx context.Context, db database.Executor, userID int64) ([]*entities.Account, error)
	}
}

func NewUserService(
	pgClient *postgres_client.PostgresClient,
	idGenerator id_utils.IDGenerator,
	userCache cache.Cache[int64, *entities.UserWithAccounts],
	userByUserNameCache cache.Cache[string, *entities.User],

) UserService {
	return &userService{
		pgClient:            pgClient,
		idGenerator:         idGenerator,
		userCache:           userCache,
		userByUserNameCache: userByUserNameCache,

		// for repositories
		userRepo:    repositories.NewUserRepository(),
		accountRepo: repositories.NewAccountRepository(),
	}
}

// CreateUser is implementation to business logic for create user.
func (s *userService) CreateUser(ctx context.Context, data *entities.User) (int64, error) {
	userCtx, err := xcontext.ExtractUserInfoFromContext(ctx)
	if err != nil {
		return 0, err
	}

	data.CreatedBy = userCtx.UserID
	if existedUser, _ := s.userByUserNameCache.Get(ctx, data.UserName); existedUser != nil {
		return 0, fmt.Errorf("username already exists")
	}

	// If user exists we should return an existed user error
	existedUser, err := s.userRepo.GetUserByUserName(ctx, s.pgClient, data.UserName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	if existedUser != nil {
		return 0, fmt.Errorf("username already exists")
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

func (s *userService) GetUserByID(ctx context.Context, id int64) (*entities.UserWithAccounts, error) {
	// If user exists in cache, we no need call to database.
	if data, err := s.userCache.Get(ctx, id); err == nil {
		return data, nil
	}

	data, err := s.userRepo.GetUserByID(ctx, s.pgClient, id)
	if err != nil {
		return nil, err
	}

	// cache user entity by id
	s.userCache.Add(ctx, id, data)

	return data, nil
}

// Update is representation of business logic to update user by id
func (s *userService) Update(ctx context.Context, data *entities.User) error {
	oldUser, err := s.userRepo.GetUserByID(ctx, s.pgClient, data.ID)
	if err != nil {
		return err
	}
	// just update, no need check exists because we will check row affected.
	if err := s.userRepo.UpdateByID(ctx, s.pgClient, data.ID, data); err != nil {
		return err
	}

	// remove from cache because user info changed
	s.userCache.Remove(ctx, data.ID)
	s.userByUserNameCache.Remove(ctx, oldUser.UserName)

	return nil
}

// DeleteByID is representation of business logic to update user by id
func (s *userService) DeleteByID(ctx context.Context, id int64) error {
	user, err := s.userRepo.GetUserByID(ctx, s.pgClient, id)
	if err != nil {
		return err
	}

	if err := s.userRepo.DeleteByID(ctx, s.pgClient, id); err != nil {
		return err
	}

	// remove from cache because user info removed
	s.userCache.Remove(ctx, id)
	s.userByUserNameCache.Remove(ctx, user.UserName)

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

	s.userCache.Remove(ctx, data.ID)

	return data.ID, nil
}

func (s *userService) ListAccountByID(ctx context.Context, id int64) ([]*entities.Account, error) {
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

// For using skeleton: s *userService UserService
