package deliveries

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"user-management/internal/entities"
	"user-management/internal/models"
	"user-management/internal/services"
	"user-management/pkg/http_server"
	"user-management/pkg/reflect_utils"
)

// using skeleton with cmd (d *userDelivery UserDelivery)
type userDelivery struct {
	server      *http_server.HttpServer
	userService services.UserService
}

func RegisterUserDelivery(
	server *http_server.HttpServer,
	userService services.UserService,
) {
	delivery := &userDelivery{
		server:      server,
		userService: userService,
	}

	http_server.Register(server, http.MethodPost, "/users", delivery.CreateUser)
	http_server.Register(server, http.MethodGet, "/users/{id}", delivery.GetUserByID)
}

func (d *userDelivery) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.CreateUserResponse, error) {
	if req.UserName == "" {
		return nil, fmt.Errorf("user must not be empty")
	}

	if req.Password == "" {
		return nil, fmt.Errorf("password must not be empty")
	}

	if req.Name == "" {
		return nil, fmt.Errorf("name must not be empty")
	}

	if !slices.Contains(entities.UserRoleList, entities.User_Role(req.Role)) {
		return nil, fmt.Errorf("user role is not valid")
	}

	user := &entities.User{}
	if err := reflect_utils.CopyStruct(req, user); err != nil {
		return nil, fmt.Errorf("unable to copy user from req: %w", err)
	}
	id, err := d.userService.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("unable to create user: %w", err)
	}

	return &models.CreateUserResponse{
		ID: id,
	}, nil
}

func (d *userDelivery) GetUserByID(_ context.Context, req *models.GetUserByIDRequest) (*models.GetUserByIDResponse, error) {

	return &models.GetUserByIDResponse{}, nil
}
