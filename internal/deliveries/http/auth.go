package deliveries

import (
	"context"
	"fmt"
	"net/http"
	"user-management/internal/entities"
	"user-management/internal/models"
	"user-management/internal/services"
	"user-management/pkg/http_server"
)

type authDelivery struct {
	server      *http_server.HttpServer
	authService services.AuthService
}

// RegisterAuthDelivery is registration of user delivery APIs to http server.
func RegisterAuthDelivery(
	server *http_server.HttpServer,
	authService services.AuthService,
) {
	delivery := &authDelivery{
		server:      server,
		authService: authService,
	}

	http_server.Register(server, http.MethodPost, "/auth/login", delivery.Login)
}

func (d *authDelivery) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	if req.UserName == "" {
		return nil, fmt.Errorf("user must not be empty")
	}

	if req.Password == "" {
		return nil, fmt.Errorf("password must not be empty")
	}

	user, token, err := d.authService.Login(ctx, &entities.User{
		UserName: req.UserName,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Name:  user.Name.String,
		Role:  string(user.Role),
		ID:    user.ID,
		Token: token,
	}, nil
}
