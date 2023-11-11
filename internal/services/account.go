package services

import (
	"context"
	"user-management/internal/entities"
	"user-management/internal/repositories"
	"user-management/pkg/database"
	"user-management/pkg/postgres_client"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

// AccountService is a service exporter to account for other layers.
type AccountService interface {
	GetAccountByID(context.Context, int64) (*entities.Account, error)
}

type accountService struct {
	pgClient     *postgres_client.PostgresClient
	accountCache *expirable.LRU[int64, *entities.Account]

	accountRepo interface {
		GetAccountByID(ctx context.Context, db database.Executor, id int64) (*entities.Account, error)
	}
}

func NewAccountService(
	pgClient *postgres_client.PostgresClient,
	accountCache *expirable.LRU[int64, *entities.Account],
) AccountService {
	return &accountService{
		pgClient:     pgClient,
		accountCache: accountCache,
		accountRepo:  repositories.NewAccountRepository(),
	}
}

func (s *accountService) GetAccountByID(_ context.Context, _ int64) (*entities.Account, error) {
	panic("not implemented") // TODO: Implement
}
