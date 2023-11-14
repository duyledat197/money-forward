package services

import (
	"context"

	"user-management/internal/entities"
	"user-management/internal/repositories"
	"user-management/pkg/cache"
	"user-management/pkg/database"
	"user-management/pkg/postgres_client"
)

// AccountService is a service exporter to account for other layers.
type AccountService interface {
	GetAccountByID(context.Context, int64) (*entities.Account, error)
}

type accountService struct {
	pgClient     *postgres_client.PostgresClient
	accountCache cache.Cache[int64, *entities.Account]

	accountRepo interface {
		GetAccountByID(ctx context.Context, db database.Executor, id int64) (*entities.Account, error)
	}
}

func NewAccountService(
	pgClient *postgres_client.PostgresClient,
	accountCache cache.Cache[int64, *entities.Account],
) AccountService {
	return &accountService{
		pgClient:     pgClient,
		accountCache: accountCache,
		accountRepo:  repositories.NewAccountRepository(),
	}
}

func (s *accountService) GetAccountByID(ctx context.Context, id int64) (*entities.Account, error) {
	if account, err := s.accountCache.Get(ctx, id); err == nil {
		return account, nil
	}

	account, err := s.accountRepo.GetAccountByID(ctx, s.pgClient, id)
	if err != nil {
		return nil, err
	}

	s.accountCache.Add(ctx, id, account)

	return account, nil
}
