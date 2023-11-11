package deliveries

import (
	"context"
	"user-management/internal/models"
)

// using skeleton with cmd (d *accountDelivery AccountDelivery)

type AccountDelivery interface {
	GetAccountByID(context.Context, *models.GetAccountByIDRequest) (*models.GetAccountByIDResponse, error)
}

type accountDelivery struct {
}

func NewAccountDelivery() AccountDelivery {
	return &accountDelivery{}
}

func (d *accountDelivery) GetAccountByID(_ context.Context, req *models.GetAccountByIDRequest) (*models.GetAccountByIDResponse, error) {

	return &models.GetAccountByIDResponse{}, nil
}
