package deliveries

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"user-management/internal/entities"
	"user-management/internal/models"
	"user-management/internal/services"
	"user-management/pkg/database"
	"user-management/pkg/http_server"
	"user-management/pkg/reflect_utils"
)

// using skeleton with cmd (d *userDelivery UserDelivery)
type userDelivery struct {
	server      *http_server.HttpServer
	userService services.UserService
}

// RegisterUserDelivery is registration of user delivery APIs to http server.
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
	http_server.Register(server, http.MethodPut, "/users/{id}", delivery.UpdateUser)

	http_server.Register(server, http.MethodPost, "/users/{id}", delivery.CreateAccountByUserID)
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

func (d *userDelivery) GetUserByID(ctx context.Context, req *models.GetUserByIDRequest) (*models.GetUserByIDResponse, error) {
	if req.ID == 0 {
		return nil, fmt.Errorf("id must not be empty")
	}

	data, err := d.userService.GetUserByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve user by id: %w", err)
	}

	return &models.GetUserByIDResponse{
		ID:   data.ID,
		Name: data.Name.String,
	}, nil
}

func (d *userDelivery) UpdateUser(ctx context.Context, req *models.UpdateUserRequest) (*models.UpdateUserResponse, error) {
	if req.ID == 0 {
		return nil, fmt.Errorf("id must not be empty")
	}

	if err := d.userService.Update(ctx, &entities.User{
		ID:   req.ID,
		Name: database.NullString(req.Name),
	}); err != nil {
		return nil, fmt.Errorf("unable to update user by id: %w", err)
	}

	return &models.UpdateUserResponse{}, nil
}

func (d *userDelivery) CreateAccountByUserID(ctx context.Context, req *models.CreateAccountByUserIDRequest) (*models.CreateAccountByUserIDResponse, error) {
	if req.UserID == 0 {
		return nil, fmt.Errorf("user id must not be empty")
	}
	id, err := d.userService.CreateAccount(ctx, &entities.Account{
		Name:   database.NullString(req.Name),
		UserID: req.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to update user by id: %w", err)
	}

	return &models.CreateAccountByUserIDResponse{
		ID: id,
	}, nil
}
