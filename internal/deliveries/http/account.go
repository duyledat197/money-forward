package deliveries

import (
	"context"
	"user-management/internal/models"
)

// using skeleton with cmd (d *accountDelivery AccountDelivery)

type AccountDelivery interface {
	CreateAccount(context.Context, *models.CreateAccountRequest) (*models.CreateAccountResponse, error)
	GetAccountByID(context.Context, *models.GetAccountByIDRequest) (*models.GetAccountByIDResponse, error)
}

type accountDelivery struct {
}

func NewAccountDelivery() AccountDelivery {
	return &accountDelivery{}
}

func (d *accountDelivery) CreateAccount(_ context.Context, _ *models.CreateAccountRequest) (*models.CreateAccountResponse, error) {
	return &models.CreateAccountResponse{}, nil
}

func (d *accountDelivery) GetAccountByID(_ context.Context, req *models.GetAccountByIDRequest) (*models.GetAccountByIDResponse, error) {

	return &models.GetAccountByIDResponse{}, nil
}
