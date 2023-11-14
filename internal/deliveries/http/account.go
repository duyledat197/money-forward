package deliveries

import (
	"context"
	"net/http"
	"user-management/internal/models"
	"user-management/internal/services"
	"user-management/pkg/http_server"
)

// using skeleton with cmd (d *accountDelivery AccountDelivery)

type AccountDelivery interface {
	GetAccountByID(context.Context, *models.GetAccountByIDRequest) (*models.GetAccountByIDResponse, error)
}

type accountDelivery struct {
	server         *http_server.HttpServer
	accountService services.AccountService
}

// RegisterAccountDelivery is registration of account delivery APIs to http server.
func RegisterAccountDelivery(
	server *http_server.HttpServer,
	accountService services.AccountService,
) {
	delivery := &accountDelivery{
		server:         server,
		accountService: accountService,
	}

	http_server.Register(server, http.MethodGet, "/accounts/{id}", delivery.GetAccountByID)
}
func (d *accountDelivery) GetAccountByID(ctx context.Context, req *models.GetAccountByIDRequest) (*models.GetAccountByIDResponse, error) {
	resp, err := d.accountService.GetAccountByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &models.GetAccountByIDResponse{
		Account: &models.Account{
			ID:      resp.ID,
			Name:    resp.Name.String,
			Balance: resp.Balance.Int64,
		},
	}, nil
}
